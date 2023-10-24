from datetime import datetime, timedelta
import logging
import pandas as pd


class OrdersDatabase():
    """

    """

    def __init__(self, bigquery_settings: dict):
        self.bq_ids = bigquery_settings

        self.orders = None
        self.products = None
        self.order_items = None

    @classmethod
    def datetime2GoogleSQL(self, d: datetime):
        """Convert python datetime-object to google-SQL
        CAST(... AS DATETIME)."""

        return """CAST("{:04d}-{:02d}-{:02d}T{:02d}:{:02d}:{:02d}" AS DATETIME)""".format(
            d.year, d.month, d.day,
            d.hour, d.minute, d.second)

    def GetOrders(self,
                  date_start: datetime = None,
                  date_end: datetime = None):
        """Get orders in timewindow from database in BigQuery."""

        # Create the SQL query
        query = """SELECT * FROM {}.{}""".format(
            self.bq_ids["dataset"],
            self.bq_ids["orders_table"]
            )
        if date_start is not None and date_end is not None:
            query += """ WHERE order_placed BETWEEN {} AND {}""".format(
                self.datetime2GoogleSQL(date_start),
                self.datetime2GoogleSQL(date_end)
            )
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

        query = """SELECT * FROM {}.{}""".format(
            self.bq_ids["dataset"],
            self.bq_ids["products_table"])
        logging.debug(f"GetProducts-query: {query}")

        return pd.read_gbq(query, project_id=self.bq_ids["project"])

    def GetOrderitems(self):
        """Get the order-items -table from BigQuery."""

        query = """SELECT * FROM {}.{}""".format(
            self.bq_ids["dataset"],
            self.bq_ids["order_items_table"])
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
