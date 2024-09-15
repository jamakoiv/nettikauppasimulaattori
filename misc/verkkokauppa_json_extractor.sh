#!/usr/bin/bash

# TODO: Make these command line arguments as well.
USER=jamakoiv
PASS=$(cat passwd)
CONN="mongodb+srv://$USER:$PASS@cosmos-mongo-testi.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000wtimeoutMS=0"
DATABASE=reviews
SIGNATURE=reviewText

# Parse input arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --save-json)
      SAVE_JSON=true
      shift # past argument
      ;;
    --upload)
      UPLOAD=true
      shift
      ;;
    -*|--*)
      echo "Unknown option $1"
      exit 1
      ;;
    *)
      URL=$1 # save positional arg
      break
      ;;
  esac
done

if [ "$SAVE_JSON" = true ]; then
  JSON_FILE=$(echo $URL | tr -d [/]).json
else 
  JSON_FILE=/dev/null
fi

# Extract the product JSON-file from the website.
JSON=$(curl -s $URL | grep $SIGNATURE | tee $JSON_FILE)
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
if [ "$UPLOAD" = true ]; then
  echo "db = db.getSiblingDB('$DATABASE');" >upload.js
  echo "db.$DATABASE.insertOne($JSON);" >>upload.js
  echo "printjson(db.$DATABASE.find());" >>upload.js

  mongosh --file=upload.js $CONN
  rm upload.js
fi

