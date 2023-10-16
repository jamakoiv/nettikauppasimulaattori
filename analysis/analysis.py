#!/usr/bin/python

import numpy as np
import pandas as pd
import seaborn as sns
import matplotlib as mpl
import matplotlib.pyplot as plt
from datetime import date, datetime, timedelta

import logging


bigquery_ids = {    "project": "nettikauppasimulaattori",
                    "dataset": "store_operational",
                    "orders_table": "orders",
                    "order_items_table": "order_items",
                    "products_table": "products" }

class OrdersDatabase():
    def __init__(self, bigquery_ids: dict):
        self.bq_ids = bigquery_ids

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
        if date_start != None and date_end != None:
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
            self.bq_ids["products_table"] )
        logging.debug(f"GetProducts-query: {query}")

        return pd.read_gbq(query, project_id=self.bq_ids["project"])
    

    def GetOrderitems(self):
        """Get the order-items -table from BigQuery."""

        query = """SELECT * FROM {}.{}""".format(
            self.bq_ids["dataset"],
            self.bq_ids["order_items_table"] )
        logging.debug(f"GetOrderItems-query: {query}")
        return pd.read_gbq(query, project_id=self.bq_ids["project"])
    

    def CalculateOrderPrices(self):
        """Calculate price, wholesale price, tax, and profit
        for each order and add these to the dataframe.
        """
        
        for i, ID in enumerate(self.orders['id']):
            tmp = db.order_items[db.order_items['order_id'] == ID]
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


def plot(orders: pd.DataFrame,
         bins: pd.DatetimeIndex,
         day: datetime):

    sns.set_theme()
    fig = plt.figure(figsize=(6, 7))
    ax = fig.subplots(2,1)
    plt.subplots_adjust(top=0.95, hspace=0.60)

    ax_daily, ax_weekly = ax

    __bins__ =  mpl.dates.date2num(bins)
    ax_daily.hist([orders['order_placed'], orders['order_placed']],
            weights=[orders['price'], orders['profit']],
            bins=__bins__)

    ax_daily.legend(['Sales', 'Profit'])
    ax_daily.set_xticks(__bins__[::2],
                        ["{:02d}".format(x) for x in np.arange(0,25,2)],
                        # rotation=45,
                        ha='right',
                        fontsize='small')
    # ax_daily.set_xlabel("Date")
    ax_daily.set_ylabel("Money €")

    ax_daily.set_title("Daily sales {}.".format(day.strftime("%d. %B %Y")))

    return fig, ax


if __name__ == "__main__":
    #logging.basicConfig(level=logging.DEBUG)

    t_start = datetime(2023, 10, 13)
    t_end = t_start + timedelta(days=1)
    t_bins = pd.date_range(t_start, t_end, freq="1H")

    db = OrdersDatabase(bigquery_ids)
    # db.orders = db.GetOrders()
    db.orders = db.GetOrders(t_start, t_end)
    db.order_items = db.GetOrderitems()
    db.products = db.GetProducts()
    db.CalculateOrderPrices()