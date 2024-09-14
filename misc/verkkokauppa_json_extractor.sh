#!/usr/bin/bash

USER=jamakoiv
PASS=$(cat passwd)
CONN="mongodb+srv://$USER:$PASS@cosmos-mongo-testi.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000wtimeoutMS=0"
DATABASE=reviews

# Extract the product JSON-file from the website.
# JSON=$(curl -s https://www.verkkokauppa.com/fi/product/893623/FWD-HP-EliteBook-840-G5-14-kaytetty-kannettava-tietokone-B-l | grep reviewText)
JSON=UGUU

# NOTE: Trying to send the data using --eval=<command> usually fails due to
# the JSON string exceeding command argument limit. That's why we make
# a separate script file instead.

# Create a mongosh-script for sending the JSON to the database.
echo "db = connect('$CONN');" >testi.js
echo "db.$DATABASE.find();" >>testi.js

mongosh --nodb --file=testi.js $CONN

# echo $CONN
# echo $JSON

# db = connect($CONN)
