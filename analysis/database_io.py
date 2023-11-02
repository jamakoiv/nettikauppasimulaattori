from datetime import datetime, timedelta

import logging
import pandas as pd

from google.cloud import bigquery


class OrdersDatabase():
    """

    """

    def __init__(self,
                 client: bigquery.Client,
                 bigquery_settings: dict,
                 queries: dict):

        self.bq_ids = bigquery_settings
        self.queries = queries
        self.client = client

        self.orders = None
        self.products = None
        self.order_items = None

    @classmethod
    def datetime2GoogleSQL(self, d: datetime):
        """Convert python datetime-object to SQL-style date-string."""

        return """CAST("{:04d}-{:02d}-{:02d}T{:02d}:{:02d}:{:02d}" AS DATETIME)""".format(
            d.year, d.month, d.day,
            d.hour, d.minute, d.second)

    def GetOrders(self,
                  date_start: datetime = None,
                  date_end: datetime = None):
        """Get orders in timewindow from database in BigQuery."""

        query = self.queries['get_table_between_dates']['sql'].format(
                    dataset=self.bq_ids["dataset_operational"],
                    table=self.bq_ids["orders_table"],
                    date_column="order_placed",
                    start_date=self.datetime2GoogleSQL(date_start),
                    end_date=self.datetime2GoogleSQL(date_end))
        logging.debug(f"GetOrders-query: {query}")

        df = pd.read_gbq(query, project_id=self.bq_ids["project"])

        # TODO: Do we even need these if using
        # mpl.dates.date2num and pd.date_range?
        tmp = df['order_placed']
        df['order_placed_date'] = [t.to_pydatetime().date() for t in tmp]
        df['order_placed_hour'] = [t.to_pydatetime().hour for t in tmp]

        return df

    def GetProducts(self):
        """Get the products-table from BigQuery."""

        query = self.queries['get_table']['sql'].format(
                    dataset=self.bq_ids["dataset_operational"],
                    table=self.bq_ids["products_table"])
        logging.debug(f"GetProducts-query: {query}")

        return pd.read_gbq(query, project_id=self.bq_ids["project"])

    def GetOrderitems(self):
        """Get the order-items -table from BigQuery."""

        query = self.queries['get_table']['sql'].format(
                    dataset=self.bq_ids["dataset_operational"],
                    table=self.bq_ids["order_items_table"])
        logging.debug(f"GetOrderItems-query: {query}")

        return pd.read_gbq(query, project_id=self.bq_ids["project"])

    def CalculateOrderPrices(self):
        """Calculate price, wholesale price, tax, and profit
        for each order and add these to the dataframe.
        """
        # TODO: Maybe should do this server-side.

        for i, ID in enumerate(self.orders['id']):
            tmp = self.order_items[self.order_items['order_id'] == ID]
            tmp = tmp.merge(right=self.products,
                            left_on='product_id',
                            right_on='id')

            wholesale_price = tmp['wholesale_price'].sum()
            price = tmp['price'].sum()
            tax = (tmp['price'] * tmp['vat']).sum()
            profit = price - wholesale_price - tax

            self.orders.loc[i, 'wholesale_price'] = wholesale_price
            self.orders.loc[i, 'price'] = price
            self.orders.loc[i, 'tax'] = tax
            self.orders.loc[i, 'profit'] = profit

    def GetAll(self,
               date_start: datetime = None,
               date_end: datetime = None):
        """Helper function to run GetOrders, GetOrderItems,
        and GetProducts in one call."""

        self.orders = self.GetOrders(date_start, date_end)
        self.order_items = self.GetOrderitems()
        self.products = self.GetProducts()

    def UpdateOrderTotals(self) -> None:
        """Update order-totals table."""

        query = self.queries['update_order_totals']['sql'].format(
                    source_dataset=self.bq_ids['dataset_operational'],
                    source_table=self.bq_ids['orders_table'],
                    insert_dataset=self.bq_ids['dataset_analysis'],
                    insert_table=self.bq_ids['order_totals_table'])
        logging.debug("UpdateOrderTotals query: {}".format(query))

        conf = bigquery.QueryJobConfig(use_legacy_sql=False)
        res = self.client.query(query, job_config=conf)
        while not res.done():
            pass

        if res.errors:
            logging.error("Error updating table {}: {}.".format(
                    self.bq_ids['order_totals_table'], res.errors))
        else:
            logging.info("Updated table {}.".format(self.bq_ids['order_totals_table']))


    def MakeARIMAHourly(self,
                        t_start: datetime,
                        t_end: datetime) -> None:
        """Create ARIMA-model for hourly sales.
        Datasets and tables are defined in self.bq_ids."""

        query = self.queries['create_arima_model']['sql'].format(
                    model_dataset=self.bq_ids['dataset_analysis'],
                    model_name=self.bq_ids['sales_model'],
                    dataset=self.bq_ids['dataset_analysis'],
                    table=self.bq_ids['order_totals_table'],
                    time_column="order_placed",
                    data_column="price",
                    start_date=self.datetime2GoogleSQL(t_start),
                    end_date=self.datetime2GoogleSQL(t_end))
        logging.debug("MakeARIMAHourly query: {}".format(query))

        # Can't use default dataset here or the we get error complaining
        # of default project missing.
        # conf = bigquery.QueryJobConfig(default_dataset=self.bq_ids['dataset_analysis'],
        #                                use_legacy_sql=False)
        # res = self.client.query(query, job_config=conf)

        conf = bigquery.QueryJobConfig(use_legacy_sql=False)
        res = self.client.query(query, job_config=conf)
        while not res.done():
            pass

        if res.errors:
            logging.error("Error creating ARIMA-model: {}".format(res.errors))
        else:
            logging.info("Succesfully created ARIMA-model.")

    def ForecastARIMAHourly(self,
                            hours_to_forecast: int):
        """Use the ML-model for hourly sales to forecast future sales."""

        query = self.queries['forecast_model']['sql'].format(
                    model_dataset=self.bq_ids['dataset_analysis'],
                    model_name=self.bq_ids['sales_model'],
                    forecast_N=hours_to_forecast)
        logging.debug("ForecastARIMAHourly query: {}".format(query))

        dest = "{}.{}.{}".format(self.bq_ids['project'],
                                 self.bq_ids['dataset_analysis'],
                                 self.bq_ids['sales_model_forecast'])
        conf = bigquery.QueryJobConfig(destination=dest,
                                       create_disposition="CREATE_IF_NEEDED",
                                       write_disposition="WRITE_APPEND",
                                       use_legacy_sql=False)

        res = self.client.query(query, job_config=conf)
        while not res.done():
            pass

        if res.errors:
            logging.error("Error forecasting with ARIMA-model: {}".format(res.errors))
        else:
            logging.info("""Succesfully created forecast from ARIMA-model {}
                         into table {}.""".format(self.bq_ids['sales_model'],
                                                  dest))

    def GetHourlySalesForecast(self,
                               start_date: datetime,
                               end_date: datetime) -> pd.DataFrame:
        """Retrieve forecasted hourly sales data."""

        query = self.queries['get_table_between_dates']['sql'].format(
                    dataset=self.bq_ids['dataset_analysis'],
                    table=self.bq_ids['sales_model_forecast'],
                    start_date=self.datetime2GoogleSQL(start_date),
                    end_date=self.datetime2GoogleSQL(end_date))

        logging.debug("GetHourlySalesForecast query: {}".format(query))

        return pd.read_gbq(query, project_id=self.bq_ids['project'])

    def UpdateCustomerStats(self) -> None:
        """Update customer_stats table."""

        query = self.queries['get_customer_stats']['sql'].format(
                    dataset_operational=self.bq_ids['dataset_operational'],
                    orders_table=self.bq_ids["orders_table"],
                    customers_table=self.bq_ids["customers_table"],
                    items_table=self.bq_ids["order_items_table"],
                    products_table=self.bq_ids["products_table"],
                    dataset_analysis=self.bq_ids["dataset_analysis"],
                    order_totals_table=self.bq_ids["order_totals_table"])

        dest = "{}.{}.{}".format(self.bq_ids["project"],
                                 self.bq_ids["dataset_analysis"],
                                 self.bq_ids["customer_stats_table"])
        conf = bigquery.QueryJobConfig(destination=dest,
                                       create_disposition="CREATE_IF_NEEDED",
                                       write_disposition="WRITE_TRUNCATE",
                                       use_legacy_sql=False)

        res = self.client.query(query, job_config=conf)
        while not res.done():
            pass

        if res.errors:
            logging.error("Error updating table {}: {}".format(dest, res.errors))
        else:
            logging.info("Succesfully updated table {}.".format(dest))
