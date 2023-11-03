#!/usr/bin/python3

import json
import yaml

import numpy as np
import pandas as pd
import seaborn as sns

from customers_raw_data import customers_ref
from analysis.database_io import OrdersDatabase

from google.cloud import bigquery
from sklearn.cluster import KMeans


f_bigquery_settings = "analysis/bigquery_ids.json"
f_queries = "analysis/queries.yaml"

k_means_N = 8


def load_settings(f_bq, f_q):
    # Kind of stupid having different file types for different settings.
    with open(f_bq, 'r') as f:
        bq_ids = json.load(f)

    with open(f_q, 'r') as f:
        queries = yaml.safe_load(f)

    return bq_ids, queries


def sanitize_dataframe(df: pd.DataFrame) -> pd.DataFrame:
    # Convert pandas.Int64Dtype to numpy.int64.

    # NOTE: np.where returns two-element tuple, hence the [0] in the end.
    for i in np.where(df.dtypes.values == pd.Int64Dtype())[0]:
        df[df.columns[i]] = pd.Series([np.int64(x) for x in df.iloc[:, i]])

    return df


def main():
    global customers_ref, customers

    bigquery_settings, queries = load_settings(f_bigquery_settings, f_queries)

    bq_client = bigquery.Client(project=bigquery_settings['project'])
    db = OrdersDatabase(bq_client, bigquery_settings, queries)

    db.UpdateCustomerStats()
    customers = db.GetCustomerStats()

    customers = sanitize_dataframe(customers)
    customers_ref = sanitize_dataframe(customers_ref)

    customers_kmeans = KMeans(n_clusters=k_means_N, n_init="auto")
    customers_ref_kmeans = KMeans(n_clusters=k_means_N, n_init="auto")

    customers_kmeans.fit(customers[['peak_activity_hour',
                                    'product_category',
                                    'average_order_price',
                                    'number_of_orders']])
    customers_ref_kmeans.fit(customers_ref[['peak_activity_hour',
                                            'max_budget',
                                            'category',
                                            'expected_value']])

    customers['cluster'] = customers_kmeans.labels_
    customers_ref['cluster'] = customers_ref_kmeans.labels_

    sns.pairplot(data=customers[['peak_activity_hour',
                                 'product_category',
                                 'average_order_price',
                                 'number_of_orders',
                                 'cluster']],
                 hue='cluster')
    sns.pairplot(data=customers_ref[['peak_activity_hour',
                                     'max_budget',
                                     'category',
                                     'expected_value',
                                     'cluster']],
                 hue='cluster')


if __name__ == "__main__":
    main()
