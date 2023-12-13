import io
import numpy as np
import pandas as pd
import seaborn as sns
import matplotlib as mpl
import matplotlib.pyplot as plt

from google.cloud import storage


def CreateFigure():
    sns.set_theme()
    fig = plt.figure(figsize=(6, 7))
    ax = fig.subplots(2, 1)
    plt.subplots_adjust(top=0.95, hspace=0.60)

    return fig, ax


def PlotDaySales(ax: plt.axis,
                 orders: pd.DataFrame,
                 bins: pd.DatetimeIndex,
                 title: str):
    """Plot sales and profits as histogram."""

    # BUG: X-axis labels give 0-24 hours regardless of
    # actual bins given in the input.

    # Plot histogram of sales and profit per hour.
    __bins__ = mpl.dates.date2num(bins)
    ax.hist([orders['order_placed'], orders['order_placed']],
            weights=[orders['price'], orders['profit']],
            bins=__bins__)

    ax.legend(['Sales', 'Profit'])
    ax.set_xticks(
            __bins__[::2],
            ["{:02d}".format(x) for x in np.arange(0, 25, 2)],
            # rotation=45,
            ha='right',
            fontsize='small')
    ax.set_title(title)
    ax.set_xlabel("Hour", fontsize='small')
    ax.set_ylabel("Money €")

    # return ax


def PlotDailyForecast(ax: plt.axis,
                      forecast: pd.DataFrame) -> None:

    # x_ = mpl.dates.date2num(forecast['forecast_timestamp'])
    forecast = forecast.sort_values(by='forecast_timestamp')
    x_ = forecast['forecast_timestamp']

    ax.plot(x_, forecast['forecast_value'], 'k:')
    # ax.plot(x_, forecast['prediction_interval_lower_bound'], 'r:')
    # ax.plot(x_, forecast['prediction_interval_upper_bound'], 'r:')


def PlotDayOrders(ax: plt.axis,
                  orders: pd.DataFrame,
                  title: str):
    """Plot order amounts."""

    # TODO...
    dates = pd.to_datetime(orders['order_placed'])


def PlotSalesHistory(ax: plt.axis,
                     orders: pd.DataFrame,
                     bins: pd.DatetimeIndex,
                     title: str):
    """Plot sales history frame and projection of future sales."""
    __bins__ = mpl.dates.date2num(bins)

    ax.hist(orders['order_placed'], weights=orders['price'], bins=__bins__)
    ax.set_xticks(ax.get_xticks())
    ax.set_xticklabels(ax.get_xticklabels(), rotation=45, ha='right', fontsize='small')
    ax.set_title(title)
    ax.set_ylabel("Sales €")


def SaveFigure2GoogleCloudStorage(fig: mpl.figure.Figure,
                                  storage_client: storage.Client,
                                  settings: dict):
    """Save figure 'fig' to google cloud storage bucket."""

    buf = io.BytesIO()
    fig.savefig(buf, format='svg')

    bucket = storage_client.bucket(settings['bucket'])
    blob = bucket.blob(settings['filename'])
    blob.upload_from_file(buf, content_type='image/svg', rewind=True)
