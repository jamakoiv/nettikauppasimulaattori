#!/usr/bin/python3

import seaborn as sns
import pandas as pd

import customer_clustering as cc
from bigquerydb import BigQueryDB

import matplotlib.pyplot as plt
from matplotlib.figure import Figure


def MakeFigure() -> Figure:
    """ """
    figsize = (9.0, 3.5)

    fig = plt.figure(figsize=figsize)
    fig.subplots(1, 3)

    return fig


def PlotClustersHistogram(
    df: pd.DataFrame, columns: list[str], fig: Figure, **kwargs
) -> Figure:
    if len(columns) > len(fig.axes):
        print(
            f"""Not enough axes to plot all required columns: 
            Figure has {len(fig.axes)} axes, received {len(columns)} 
            data columns to plot."""
        )

    for column, ax in zip(columns, fig.axes):
        sns.histplot(
            data=df,
            x=column,
            hue="group",
            multiple="dodge",
            ax=ax,
            **kwargs,
        )

    for ax in fig.axes[1:]:
        ax.get_legend().remove()
        ax.set_ylabel("")

    return fig


if __name__ == "__main__":
    sns.set_theme()

    db = BigQueryDB("nettikauppasimulaattori", "store_analysis_prod")
    data = cc.GetCustomerData(db)
    data, sil_avg, sil_samples = cc.KMeansCustomers(data, n_clusters=4, normalize=True)

    fig = MakeFigure()
    fig = PlotClustersHistogram(
        data,
        ["average_order_price", "average_profit", "favourite_product_category"],
        fig,
        bins=10,
        palette="tab10",
    )

    labels = ["Avg. order total [€]", "Avg. profit [€]", "Product category"]
    for label, ax in zip(labels, fig.axes):
        ax.set_xlabel(label)

    fig.subplots_adjust(
        top=0.948, bottom=0.193, left=0.08, right=0.98, hspace=0.2, wspace=0.29
    )
