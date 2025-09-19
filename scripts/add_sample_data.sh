#!/bin/bash

# Sample data creation script for Home Automation Dashboard
# This script adds some sample tasks and devices to demonstrate the functionality

BASE_URL="http://localhost:8080"

echo "Adding sample tasks..."

curl -X POST "$BASE_URL/tasks" -d "title=Check HVAC filters&description=Replace or clean HVAC filters to maintain efficiency"

curl -X POST "$BASE_URL/tasks" -d "title=Test smoke detectors&description=Test all smoke detectors and replace batteries if needed"

curl -X POST "$BASE_URL/tasks" -d "title=Review energy usage&description=Check monthly energy consumption and identify optimization opportunities"

curl -X POST "$BASE_URL/tasks" -d "title=Update security cameras&description=Check security camera firmware and update if necessary"

echo ""
echo "Adding sample devices..."

curl -X POST "$BASE_URL/devices" -d "name=Front Door Lock&type=lock&location=Front Door"

curl -X POST "$BASE_URL/devices" -d "name=Kitchen Light&type=light&location=Kitchen"

curl -X POST "$BASE_URL/devices" -d "name=Bedroom Temperature Sensor&type=sensor&location=Bedroom"

curl -X POST "$BASE_URL/devices" -d "name=Garage Door Opener&type=switch&location=Garage"

curl -X POST "$BASE_URL/devices" -d "name=Security Camera&type=camera&location=Backyard"

echo ""
echo "Sample data added successfully!"
echo "Visit http://localhost:8080 to see the tasks"
echo "Visit http://localhost:8080/devices to see the devices"