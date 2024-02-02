#!/usr/bin/python3

import pandas as pd
import numpy as np
from bigquerydb import BigQueryDB
from sklearn import cluster, preprocessing, metrics

from typing import Tuple


def GetCustomerData(db: BigQueryDB) -> pd.DataFrame:
    """ """

    table = "CA_aggregate_customer_analysis_data"
    query = f"SELECT * FROM `{db.project}.{db.dataset}.{table}`"

    res = pd.read_gbq(query, project_id=db.project)
    return res


def NormalizeCustomerData(df: pd.DataFrame) -> pd.DataFrame:
    """ """

    res = preprocessing.normalize(df)
    res = pd.DataFrame(res, columns=df.keys())

    return res


def KMeansCustomers(
    df: pd.DataFrame, *args, **kwargs
) -> Tuple[pd.DataFrame, float, np.ndarray]:
    """ """

    fit_data = df.drop("customer_id", axis=1)

    if "normalize" in kwargs.keys():
        fit_data = NormalizeCustomerData(fit_data)
        kwargs.pop("normalize")

    model = cluster.KMeans(*args, **kwargs)
    model.fit(fit_data)

    df = df.assign(group=model.labels_)
    sil_score = metrics.cluster.silhouette_score(fit_data, model.labels_)
    sil_sample = metrics.cluster.silhouette_samples(fit_data, model.labels_)

    return df, sil_score, sil_sample


def AnalyzeSilhouetteScore(
    df: pd.DataFrame, sil_score: np.ndarray, sil_sample: np.ndarray
) -> None:
    """ """

    pass


if __name__ == "__main__":
    db = BigQueryDB("nettikauppasimulaattori", "store_analysis_prod")
    data = GetCustomerData(db)
    data_norm = NormalizeCustomerData(data)

    data, sil_avg, sil_samples = KMeansCustomers(data, n_clusters=4, normalize=True)
