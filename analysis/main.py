
import json
import logging
import yaml
import pandas as pd
import pytz

from datetime import date, datetime, timedelta
from database_io import OrdersDatabase
from analyze_and_plot import CreateFigure, PlotDaySales, PlotDailyForecast, PlotSalesHistory, SaveFigure2GoogleCloudStorage

from google.cloud import storage, bigquery

# Stuff needed for google-cloud-functions.
from markupsafe import escape
import functions_framework


# Settings files.
f_bigquery = "bigquery_ids.json"
f_cloud_storage = "cloud_storage_ids.json"
f_queries = "queries.yaml"


DAILY = "1D"
HOURLY = "1H"


def GetCurrentDate(timezone: str) -> datetime:
    """Return current date in timezone."""

    tz = pytz.timezone(timezone)
    return datetime.now().astimezone(tz)


# Register function for google-cloud-functions framework.
@functions_framework.cloud_event
def Run(event):
    # Export to global namespace for using interactively.
    global fig, ax, db
    global gcs_client, bq_client
    global t_bins_daily, t_bins_longterm, t_end, t_start_arima
    global previous_sales_forecast, sales_forecast

    logging.basicConfig(level=logging.DEBUG)

    # NOTE: This script plots data from last 24 hours,
    # Currently it is assumed that this will be run after midnight.
    t_end = GetCurrentDate("Europe/Helsinki")
    t_start_daily = t_end - timedelta(days=1)
    t_start_longterm = t_end - timedelta(days=7)
    t_start_arima = t_end - timedelta(days=14)
    t_title = t_start_daily

    t_bins_daily = pd.date_range(t_start_daily, t_end, freq=HOURLY)
    t_bins_longterm = pd.date_range(t_start_longterm, t_end + timedelta(days=3), freq=DAILY)

    # Load settings.
    bigquery_settings = json.load(open(f_bigquery, 'r'))
    cloud_storage_settings = json.load(open(f_cloud_storage, 'r'))

    # Load queries.
    queries = yaml.safe_load(open(f_queries, 'r'))

    # Create clients for google cloud.
    bq_client = bigquery.Client(project=bigquery_settings['project'])
    gcs_client = storage.Client(project=cloud_storage_settings['project'])

    # Retrieve data.
    db = OrdersDatabase(bq_client, bigquery_settings, queries)
    db.GetAll(t_start_daily, t_end)
    db.CalculateOrderPrices()

    # Save yesterdays forecast, create new model, and get new forecast.
    previous_sales_forecast = db.GetHourlySalesForecast(t_start_daily, t_end)
    db.UpdateOrderTotals()
    db.MakeARIMAHourly(t_start_arima, t_end)
    db.ForecastARIMAHourly(24)
    # sales_forecast = db.GetHourlySalesForecast()

    # Plot and save figure.
    fig, ax = CreateFigure()
    ax_daily, ax_longterm = ax
    title = "Hourly sales {}.".format(t_title.strftime("%d. %B %Y"))
    filename = "sales_{}.svg".format(t_title.strftime("%Y_%m_%d"))

    PlotDaySales(ax_daily, db.orders, t_bins_daily, title)
    PlotDailyForecast(ax_daily, previous_sales_forecast)

    db.orders = db.GetOrders(t_start_longterm, t_end)
    db.CalculateOrderPrices()
    PlotSalesHistory(ax_longterm, db.orders, t_bins_longterm, "Weekly sales")

    cloud_storage_settings['filename'] = filename
    cloud_storage_settings['filename'] = "testing"
    # SaveFigure2GoogleCloudStorage(fig, gcs_client, cloud_storage_settings)


# For executing Run-function in local machine.
class DummyEvent():
    def __init__(self):
        pass

# For running in local machine...
if __name__ == "__main__":
    Run(DummyEvent())
