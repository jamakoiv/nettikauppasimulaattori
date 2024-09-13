#!/usr/bin/bash

USER="jamakoiv"
CONN="mongodb+srv://$USER@cosmos-mongo-testi.mongocluster.cosmos.azure.com/?tls=true&authMechanism=SCRAM-SHA-256&retrywrites=false&maxIdleTimeMS=120000"

JSON=$(curl -s https://www.verkkokauppa.com/fi/product/893623/FWD-HP-EliteBook-840-G5-14-kaytetty-kannettava-tietokone-B-l | grep reviewText)

echo $CONN
echo $JSON

# db = connect($CONN)
