
Add simulated orders from simulated customers into google BigQuery-table.




Table schemas in BigQuery (Data type in []):
-------------------------------------------------

table:	customers
columns:	id [int], 
			first_name [string], 
			last_name [string], 
			most_active [int], 
			price_point [int],

table:	products
columns:	id [int],
			name [string],
			wholesale_price [int],
			price [int],
			vat [float]

table:	orders
columns:	id [int],
			customer_id [int],
			delivery_type [int],
			status [int],

table:	order_items
columns:	order_id [int],
			product_id [int],

table:	delivery_types
columns:	delivery_id [int],
			type_name [string],

table:	orderstatus_types
columns:	status_id [int],
			type_name [string],

