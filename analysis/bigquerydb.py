import pandas as pd
from google.cloud import bigquery


class BigQueryDB:
    def __init__(self, project: str, dataset: str) -> None:
        self.project = project
        self.dataset = dataset

        self.DB = bigquery.Client(project=project)


    def getTable(self, table: str) -> pd.DataFrame:
        query = f"SELECT * FROM `{self.project}.{self.dataset}.{table}`"

        return pd.read_gbq(query, project_id=self.project)


    def writeTable(self, 
                   data: pd.DataFrame,
                   table: str,
                   *args, **kwargs) -> None:

        dest = f"{self.dataset}.{table}"
        data.to_gbq(destination_table=dest,
                    project_id=self.project,
                    *args,
                    **kwargs) 


    def deleteTable(self, table: str) -> None:
        dest = f"{self.dataset}.{table}"
        self.DB.delete_table(dest)


    def close(self) -> None:
        self.DB.close()