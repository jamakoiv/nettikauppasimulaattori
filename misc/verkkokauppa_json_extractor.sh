#!/usr/bin/bash

USER=jamakoiv
PASS=$(cat passwd)
CONN="mongodb+srv://$USER:$PASS@cosmos-mongo-testi.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000wtimeoutMS=0"
DATABASE=reviews

URL=$1
SIGNATURE=reviewText # Some grep-line from the JSON which we want to capture.

# Extract the product JSON-file from the website.
JSON_FILE=$(echo $URL | tr -d [/])
JSON=$(curl -s $URL | grep $SIGNATURE | tee $JSON_FILE.json)
RES=$?
if [ $RES -eq 1 ]; then
  echo Error number $RES: probably grep did not find the SIGNATURE string it was looking for, or the automatically generated JSON_FILE filename contains illegal characters.
  exit 10
elif [ $RES -ne 0 ]; then
  echo Error number $RES: probably retrieving the website '$URL' failed.
  exit 11
fi

# NOTE: Trying to send the data using --eval=<command> usually fails due to
# the JSON string exceeding command argument limit. That's why we make
# a separate script file instead.

# Create a mongosh-script for sending the JSON to the database.
# echo "db = db.getSiblingDB('$DATABASE');" >testi.js
# echo "db.$DATABASE.insertOne($JSON);" >>testi.js
# echo "printjson(db.$DATABASE.find());" >>testi.js
# 
# mongosh --file=testi.js $CONN

# rm testi.js
