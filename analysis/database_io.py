from datetime import datetime, timedelta

import logging
import pandas as pd

from google.cloud import bigquery


class OrdersDatabase():
    """

    """

    def __init__(self, bigquery_settings: dict, queries: dict):
        self.bq_ids = bigquery_settings

        self.orders = None
        self.products = None
        self.order_items = None
        self.queries = queries

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


    def MakeARIMAHourly(self,
                        client: bigquery.Client,
                        t_start: datetime,
                        t_end: datetime) -> None:

        query = self.queries['create_arima_model']['sql'].format(
                    model_dataset=self.bq_ids['dataset_analysis'],
                    model_name=self.bq_ids['sales_model'],
                    dataset=self.bq_ids['dataset_analysis'],
                    table=self.bq_ids['order_totals_table'],
                    time_column="order_placed",
                    data_column="price",
                    start_date=self.datetime2GoogleSQL(t_start),
                    end_date=self.datetime2GoogleSQL(t_end))

        conf = bigquery.QueryJobConfig(default_dataset=self.bq_ids['dataset_analysis'],
                                       use_legacy_sql=False)
        res = client.query(query, job_config=conf)
        while not res.done():
            pass

        if res.errors:
            logging.error("Error creating ARIMA-model: {}".format(res.errors))
        else:
            logging.info("Succesfully created ARIMA-model.")

    def QueryARIMAHourly(self, hours_to_forecast: int) -> pd.DataFrame:

        query = self.queries['forecast_model']['sql'].format(
                    model_dataset=self.bq_ids['dataset_analysis'],
                    model_name=self.bq_ids['sales_model'],
                    forecast_N=hours_to_forecast)

        df = pd.read_gbq(query, project_id=self.bq_ids['project'])

        return df
