
import json
import pandas as pd

from datetime import datetime, timedelta
from database_io import OrdersDatabase
from analyze_and_plot import CreateFigure, PlotDaySales, PlotHistoryAndProjetion, SaveFigure2GoogleCloudStorage

from google.cloud import storage

f_bigquery = "bigquery_ids.json"
f_cloud_storage = "cloud_storage_ids.json"

DAILY = "1D"
HOURLY = "1H"


def main():
    global fig, ax, db, gcs_client
    global t_bins_daily, t_bins_longterm

    # logging.basicConfig(level=logging.DEBUG)

    t_end = datetime(2023, 10, 18)
    t_start_daily = t_end - timedelta(days=1)
    t_start_longterm = t_end - timedelta(days=7)

    t_bins_daily = pd.date_range(t_start_daily, t_end, freq=HOURLY)
    t_bins_longterm = pd.date_range(t_start_longterm, t_end, freq="3H")

    bigquery_settings = json.load(open(f_bigquery, 'r'))
    cloud_storage_settings = json.load(open(f_cloud_storage, 'r'))

    db = OrdersDatabase(bigquery_settings)
    db.GetAll(t_start_daily, t_end)
    db.CalculateOrderPrices()

    fig, ax = CreateFigure()
    ax_daily, ax_longterm = ax
    title = "Daily sales {}.".format(t_end.strftime("%d. %B %Y"))
    filename = "sales_{}.svg".format(t_end.strftime("%Y_%m_%d"))

    gcs_client = storage.Client()
    PlotDaySales(ax_daily, db.orders, t_bins_daily, title)

    db.orders = db.GetOrders(t_start_longterm, t_end)
    db.CalculateOrderPrices()
    PlotHistoryAndProjetion(ax_longterm, db.orders, t_bins_longterm, "Weekly sales")

    cloud_storage_settings['filename'] = filename
    SaveFigure2GoogleCloudStorage(fig, gcs_client, cloud_storage_settings)


if __name__ == "__main__":
    main()
