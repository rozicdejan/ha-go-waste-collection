#!/bin/bash

# Start waste collection service
echo "Starting Waste Collection Add-on..."

# Run the Go script in a loop to update data every 15 minutes
while true; do
    ./waste-collection
    echo "Data updated. Sleeping for 15 minutes."
    sleep 900 # 900 seconds = 15 minutes
done
