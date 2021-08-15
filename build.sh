#!/bin/bash

docker build --build-arg C2B_ONLINE_PARTY_B=${C2B_ONLINE_PARTY_B} \
             --build-arg C2B_ONLINE_CHECKOUT_CALLBACK_URL=${C2B_ONLINE_CHECKOUT_CALLBACK_URL} \
             --build-arg MPESA_SHORTCODE=${MPESA_SHORTCODE} \
             --build-arg MPESA_URL=${MPESA_URL} \
             --build-arg LNM_PASSKEY=${LNM_PASSKEY} \
             --build-arg CONSUMER_KEY=${CONSUMER_KEY} \
             --build-arg CONSUMER_SECRET=${CONSUMER_SECRET} \
             --build-arg MPESA_AUTH_URL=${MPESA_AUTH_URL} \
             --build-arg DB_USER=${DB_USER} \
             --build-arg DB_PASS=${DB_PASS} \
             --build-arg DB_NAME=${DB_NAME} \
             --build-arg DB_HOST=${DB_HOST} \
             -t samplegoapp:latest . --no-cache