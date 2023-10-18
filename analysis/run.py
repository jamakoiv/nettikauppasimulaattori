from database_io import *
from analyze_and_plot import *

from google.cloud import storage


def main():
    global fig, ax, db, gcs_client

     #logging.basicConfig(level=logging.DEBUG)

    t_start = datetime(2023, 10, 13)
    t_end = t_start + timedelta(days=1)
    t_bins = pd.date_range(t_start, t_end, freq="1H")

    db = OrdersDatabase(bigquery_ids)
    db.orders = db.GetOrders(t_start, t_end)
    db.order_items = db.GetOrderitems()
    db.products = db.GetProducts()
    db.CalculateOrderPrices()


    fig, ax = CreateFigure()
    ax_daily, ax_longterm = ax
    title = "Daily sales {}.".format(t_start.strftime("%d. %B %Y"))
    filename = "sales_{}.svg".format(t_start.strftime("%Y_%m_%d"))
    
    gcs_client = storage.Client()
    PlotDaySales(ax_daily, db.orders, t_bins, title)
    SaveFigure2GoogleCloudStorage(fig, gcs_client, filename)


if __name__ == "__main__":
    main()