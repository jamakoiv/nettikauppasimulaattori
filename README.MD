
Program for simulating an internet store. Includes a simple program for simulating customers 
who place orders to the store, and workers who ship the orders.
Ultimate goal is to forecast daily and weekly sales by hour from the sales-data 
(ARIMA or linear regression), and classify different types of customers 
by their behaviour (time of day, class of items, money used etc.) 
(k-clustering or RandomForestClassifier). Maybe 

Includes automated scripts for plotting daily/weekly/montly sales and profits
and projections for weekly and monthly sales. Save graphs to google-cloud-storage.


Automated back-up of operational database to google-cloud-storage.
Weekly.

The customer/worker program is run on google cloud functions 
and orders are sent to BigQuery tables. Analysis and backup are also
run on GC functions. Timing is done via GC schedule (cron) and Pub/Sub.


The project in BigQuery contains saved SQL-queries for creating the necessary tables.
Datasets must be added manually.



TODO: Create a settings-file for BigQuery credentials, dataset- & table-ids etc.

TODO: Create the analysis-scripts.
TODO: Create the backup-scripts.

TODO: Create tests...
TODO: Create option for selecting if we want to upload orders from local test-runs.



Table schemas in BigQuery (Data type in []):
-------------------------------------------------

In dataset `store_operational`:

table:	customers
columns:	id			[int]
			first_name	[string]
			last_name	[string]

table:	products
columns:	id				[int]
			name			[string]
			wholesale_price [int]
			price			[int]
			vat				[float]

table:	orders
columns:	id				[int]
			customer_id		[int]
			delivery_type	[int]
			status			[int]
			order_placed	[datetime]
			shipping_date	[date]
			tracking_number	[int]

# Contains items for each order-id.
table:	order_items	
columns:	order_id	[int]
			product_id	[int]


In dataset `store_analysis`:

table:	order_totals
columns:	order_id	[int]
			order_placed[datetime]
			price		[int]
			wholesale_price [int]
			tax			[int]
			profit		[int]


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
