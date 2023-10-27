
import json
import pandas as pd
import pytz

from datetime import date, datetime, timedelta
from database_io import OrdersDatabase
from analyze_and_plot import CreateFigure, PlotDaySales, PlotSalesHistory, SaveFigure2GoogleCloudStorage

from google.cloud import storage

# Stuff needed for google-cloud-functions.
from markupsafe import escape
import functions_framework


# Settings files.
f_bigquery = "bigquery_ids.json"
f_cloud_storage = "cloud_storage_ids.json"


DAILY = "1D"
HOURLY = "1H"


def GetCurrentDate(timezone: str) -> datetime:
    """Return current date in timezone."""

    tz = pytz.timezone(timezone)
    return datetime.now().astimezone(tz)


# Register function for google-cloud-functions framework.
@functions_framework.cloud_event
def Run(event):
    global fig, ax, db, gcs_clienst
    global t_bins_daily, t_bins_longterm

    # logging.basicConfig(level=logging.DEBUG)

    # NOTE: This script plots data from last 24 hours,
    # Currently it is assumed that this will be run after midnight.
    t_end = GetCurrentDate("Europe/Helsinki")
    t_start_daily = t_end - timedelta(days=1)
    t_start_longterm = t_end - timedelta(days=14)
    t_title = t_start_daily

    t_bins_daily = pd.date_range(t_start_daily, t_end, freq=HOURLY)
    t_bins_longterm = pd.date_range(t_start_longterm, t_end, freq=DAILY)

    # Load settings.
    bigquery_settings = json.load(open(f_bigquery, 'r'))
    cloud_storage_settings = json.load(open(f_cloud_storage, 'r'))

    # Retrieve data.
    db = OrdersDatabase(bigquery_settings)
    db.GetAll(t_start_daily, t_end)
    db.CalculateOrderPrices()

    # Plot and save figure.
    fig, ax = CreateFigure()
    ax_daily, ax_longterm = ax
    title = "Hourly sales {}.".format(t_title.strftime("%d. %B %Y"))
    filename = "sales_{}.svg".format(t_title.strftime("%Y_%m_%d"))

    gcs_client = storage.Client()
    PlotDaySales(ax_daily, db.orders, t_bins_daily, title)

    db.orders = db.GetOrders(t_start_longterm, t_end)
    db.CalculateOrderPrices()
    PlotSalesHistory(ax_longterm, db.orders, t_bins_longterm, "Weekly sales")

    cloud_storage_settings['filename'] = filename
    SaveFigure2GoogleCloudStorage(fig, gcs_client, cloud_storage_settings)


# For running in local machine...
if __name__ == "__main__":
    Run()
