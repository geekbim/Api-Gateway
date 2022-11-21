#!/bin/sh
CURRENT_DIR=$(pwd)
   KRAKEND_PORT="4400"
   KRAKEND_BACKEND_HOST_API="http://localhost"
   KRAKEND_BACKEND_COOKIE="token"


echo "---- Generate Krakend Config ----"

FILE_ORIGINAL="$CURRENT_DIR/env/krakend.json"
FILE="$CURRENT_DIR/dist/krakend.json"
cp -f $FILE_ORIGINAL $FILE

sed -i "s~{KRAKEND_PORT}~${KRAKEND_PORT}~" $FILE
sed -i "s~{KRAKEND_BACKEND_HOST_API}~${KRAKEND_BACKEND_HOST_API}~" $FILE
sed -i "s~{KRAKEND_BACKEND_COOKIE}~${KRAKEND_BACKEND_COOKIE}~" $FILE

echo "---- Krakend Config Generated ----"
