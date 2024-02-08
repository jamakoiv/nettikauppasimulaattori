
#### Table schemas in BigQuery (Data type in []):

Dataset `store_operational` contains database of customer information,
individual orders placed by customers, and products available in the store.

**customers:**
 - id [*int*]
 - first_name [*string*]
 - last_name [*string*]

**products:** 
 - id [*int*], 
 - name [*string*]
 - wholesale_price [*int*]
 - price [*int*]
 - vat [*float*]

**orders:**
 - id				[*int*]
 - customer_id		[*int*]
 - delivery_type	[*int*]
 - status			[*int*]
 - order_placed	[*datetime*]
 - shipping_date	[*date*]
 - tracking_number	[*int*]

**order_items:**
 - order_id		[*int*]
 - product_id	[*int*]


Dataset `store_analysis` contains aggregate data created from the data in


Dataset `store_constants` contains helper data for 
converting constant values stored as integers 
in `store_operational` to human-readable form.

**delivery_types**
 - delivery_id  [*int*],
 - type_name    [*string*]

**orderstatus_types**
 - status_id	[*int*]
 - type_name	[*string*]

**product_categories**
 - category_id  [*int*]
 - type_name 	[*string*]