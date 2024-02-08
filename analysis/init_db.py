import pandas as pd
import numpy as np
import random
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



def create_customers(N: int, 
                     income_per_area: pd.DataFrame,
                     destination_table: pd.DataFrame) -> pd.DataFrame: 

    weights = np.random.random(size=len(income_per_area))
    areas = random.choices(income_per_area['code'], weights, k=N)

    # 
    income_floor = 8000
    income_coeff = 0.02 

    max_purchase_chance = 0.10

    for i, area in enumerate(areas):
        active = np.random.normal(15.0, 6) % 24  # modulo forces value to 0-24 range.

        mask = income_per_area['code'] == area
        budget = (income_per_area[mask]['income'].values[0] - income_floor)
        budget *= np.random.random() * income_coeff

        purchase_chance = np.random.random() * max_purchase_chance

        destination_table.loc[i] = [i, '', area, 
                                    active, purchase_chance, 
                                    budget, '{1;2;3;4}']

    return destination_table