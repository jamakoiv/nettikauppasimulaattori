
Generate and analyze data for a simulated online-store. 


### Program outline

<img src="readme_assets/project_outline.svg"/>

Main part of the program consist of simulated customers sending orders to the operational database.
The simulation-code can be run in either 
[Google-cloud functions](https://cloud.google.com/functions/), Docker-container 
on a cloud platform, or locally. First two are preferred due to easy automation of
running the code periodically.
Operational and analysis databases use [Google BigQuery](https://cloud.google.com/bigquery/)
as back-end for storing the data. [DBT](https://www.getdbt.com/) is used for 
importing and transforming the raw data in the operational database to a form 
suitable for analyzing business-related values, such as weekly income etc.
DBT related code can be found from the repository 
[nettikauppasimulaattori-dbt](https://github.com/jamakoiv/nettikauppasimulaattori_dbt).

<img src="readme_assets/histograms.png" width="90%">

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
