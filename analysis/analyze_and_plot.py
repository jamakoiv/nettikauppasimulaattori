import io
import numpy as np
import pandas as pd
import seaborn as sns
import matplotlib as mpl
import matplotlib.pyplot as plt

from google.cloud import storage

storage_ids = { "project": "nettikauppasimulaattori",
                "bucket": "nettikauppasimulaattori_analysis"}

def CreateFigure():
    sns.set_theme()
    fig = plt.figure(figsize=(6, 7))
    ax = fig.subplots(2,1)
    plt.subplots_adjust(top=0.95, hspace=0.60)

    return fig, ax

def PlotDaySales(ax: plt.axis,
              orders: pd.DataFrame,
              bins: pd.DatetimeIndex,
              title: str):
    """Plot sales and profits as histogram."""

    __bins__ =  mpl.dates.date2num(bins)
    ax.hist([orders['order_placed'], orders['order_placed']],
            weights=[orders['price'], orders['profit']],
            bins=__bins__)

    ax.legend(['Sales', 'Profit'])
    ax.set_xticks(__bins__[::2],
                        ["{:02d}".format(x) for x in np.arange(0,25,2)],
                        # rotation=45,
                        ha='right',
                        fontsize='small')
    ax.set_title(title)
    ax.set_xlabel("Hour", fontsize='small')
    ax.set_ylabel("Money €")

    # return ax

def PlotHistoryAndProjetion():
    ...


def SaveFigure2GoogleCloudStorage(fig: mpl.figure.Figure,
                                  storage_client: storage.Client,
                                  filename: str):
    """Save figure 'fig' to google cloud storage bucket."""
    
    buf = io.BytesIO()
    fig.savefig(buf, format='svg')

    bucket = storage_client.bucket(storage_ids['bucket'])
    blob = bucket.blob(filename)
    blob.upload_from_file(buf, content_type='image/svg', rewind=True)