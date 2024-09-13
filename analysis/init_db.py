import pandas as pd
import numpy as np
import random
import copy

from pathlib import Path
from dataclasses import dataclass


delivery_types = pd.DataFrame(
    [[0, "SHIP TO CUSTOMER"], [1, "COLLECT FROM STORE"]],
    columns=["id", "label"]
)

orderstatus_types = pd.DataFrame(
    [[0, "SHIPPED"], [1, "PENDING"], [2, "EMPTY"]],
    columns=['id', 'label']
)

orders = pd.DataFrame({
    'id': pd.Series(dtype='int'),
    'customer_id': pd.Series(dtype='int'),
    'status': pd.Series(dtype='int'),
    'order_placed': pd.Series(dtype='datetime64[ns]'),
    'shipping_date': pd.Series(dtype='datetime64[ns]'),
    'tracking_number': pd.Series(dtype='int'),
    'last_modified': pd.Series(dtype='datetime64[ns]')
})



order_items = pd.DataFrame({
    'order_id': pd.Series(dtype='int'),
    'product_id': pd.Series(dtype='int')
})

customers = pd.DataFrame({
    'id': pd.Series(dtype='int'),
    'name': pd.Series(dtype='str'),
    'area': pd.Series(dtype='int'),
    'most_active': pd.Series(dtype='int'),
    'purchase_chance': pd.Series(dtype='float'),
    'max_budget': pd.Series(dtype='float'),
    'product_categories': pd.Series(dtype='str')
})


@dataclass
class Customer:
    id: int
    name: str
    area: int
    most_active: int
    max_budget: float
    purchase_chance: float
    product_categories: list[int]

    ...


def import_income_per_area(f: str | Path) -> pd.DataFrame | None:
    """Import CSV-file containing median income per area.
    ------------
    f: path to CSV-file. See 'seed_median_income.csv' for file format.

    -> on success DataFrame, on failure None.
    """
    try: 
        res = pd.read_csv(f, 
                        header=1, 
                        names=['region', 'area', 'income'], 
                        thousands=" ")
    except FileNotFoundError:
        return

    res = res.drop_duplicates(ignore_index=True)
    return res


def create_customers(N: int, 
                     income_per_area: pd.DataFrame) -> pd.DataFrame: 
    """Create customers with income based on median income per area.
    -------------
    N: number of customers to create.
    income_per_area: DataFrame containing median income for each area.
    
    -> DataFrame containing the new customers.
    """

    # Index change later creates bug in random.choices.
    # Make local copy so original table is not mangled. 
    income_per_area = copy.copy(income_per_area)

    weights = np.random.random(size=len(income_per_area))
    areas = random.choices(income_per_area['area'], weights, k=N)

    #  
    income_floor = 8000
    income_coeff = 0.02 

    income_per_area.index = income_per_area['area']
    income = [ (income_per_area.loc[area]['income'] - income_floor) * 
              income_coeff * np.random.random()
              for area in areas ]

    max_purchase_chance = 0.10
    purchase_chance = np.random.random(size=N) * max_purchase_chance

    # modulo forces value to 0-24 range.
    active = np.random.normal(15.0, 6, size=N) % 24  

    id = np.arange(N)

    res = pd.DataFrame({'id': pd.Series(id, dtype='int'), 
        'name': pd.Series(dtype='str'),
        'area': pd.Series(areas, dtype='int'),
        'most_active': pd.Series(active, dtype='float'),
        'purchase_chance': pd.Series(purchase_chance, dtype='float'),
        'max_budget': pd.Series(income, dtype='float'),
        'product_categories': pd.Series(dtype='str')
    })

    return res


if __name__ == "__main__":
    income_per_area = import_income_per_area("seed_median_income.csv")

    df = create_customers(5000, income_per_area)
