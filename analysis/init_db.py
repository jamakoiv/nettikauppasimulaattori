import pandas as pd
import numpy as np

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

