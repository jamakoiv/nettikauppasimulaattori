
# Nettikauppasimulaattori

Program for simulating an internet store. Includes a simple program for simulating customers 
who place orders to the store, and workers who ship the orders.
Ultimate goal is to forecast daily and weekly sales by hour from the sales-data 
(ARIMA or linear regression), and classify different types of customers 
by their behaviour (time of day, class of items, money used etc.) 
(k-clustering or RandomForestClassifier).
Other goal is to forecast possible new product sales using linear regression on the sales-data,
and using logistical regression to classify which customers might buy the
new product.

<img src="readme_assets/project_outline.svg"/>

Includes automated scripts for plotting daily/weekly/montly sales and profits
and projections for weekly and monthly sales. Save graphs to google-cloud-storage.

The customer/worker program is run on google cloud functions 
and orders are sent to BigQuery tables. Analysis and backup are also
run on GC functions. Timing is done via GC schedule (cron) and Pub/Sub.

The project in BigQuery contains saved SQL-queries for creating the necessary tables.
Datasets must be added manually.

Moved most of data-transforms to DBT. 


TODO: Fix bugs/shoddy logic in analysis-scripts.
TODO: Change analysis-scripts to use DBT-created tables.

TODO: Add DBT-tests.
TODO: Create separate devel- and production environments for DBT.
TODO: Add snapshots in DBT.
TODO: Add data-freshness warnings in DBT.

TODO: Create rest of unit-tests.
TODO: Create option for selecting if we want to upload orders from local test-runs.


### Table schemas in BigQuery (Data type in []):

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
`store_operational'. `order_totals` contains data used for forecasting future sales,
and `customer_stats` contains data used for analyzing customers.
In dataset `store_analysis`:

table:	order_totals
columns:	order_id	[int]
			order_placed[datetime]
			price		[int]
			wholesale_price [int]
			tax			[int]
			profit		[int]

table:	customer_stats
columns:	customer_id			[int]
			peak_activity_hour	[int]
			product_category	[int]
			products_in_category	[int]
			average_order_price		[int]
			number_of_orders		[int]


Dataset `store_constants` contains helper data for converting some values
stored as integers in `store_operational` to human-readable form.
In dataset `store_constants`:

table:	delivery_types
columns:	delivery_id [int],
			type_name	[string]

table:	orderstatus_types
columns:	status_id	[int]
			type_name	[string]

table:	product_categories
columns:	category_id [int]
			type_name 	[string]
