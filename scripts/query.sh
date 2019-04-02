#!/bin/bash

SERVER="localhost:8080"

echo "Get first 10 buildings with pagination"
rc=$(curl -s ${SERVER}/v1/buildings?from=0&limit=10)
echo -e "$rc\n\n"

echo "Get 10 buildings from id 11 with pagination"
rc=$(curl -s ${SERVER}/v1/buildings?from=11&limit=10)
echo -e "$rc\n\n"

echo "Get the buildings of id 1"
rc=$(curl -s ${SERVER}/v1/buildings/1)
echo -e "$rc\n\n"

echo "Get first 10 buildings constucted in year 2008 order by id with pagination"
rc=$(curl -s ${SERVER}/v1/buildings/years/2008?from=0&limit=10)
echo -e "$rc\n\n"

echo "Get the height's and area's statistics(min max avg count) of buildings constructed group by years"
rc=$(curl -s ${SERVER}/v1/buildings/stats/years)
echo -e "$rc\n\n"

echo "Get the height's and area's statistics(min max avg count) of buildings constructed in year 2008"
rc=$(curl -s ${SERVER}/v1/buildings/stats/years/2008)
echo -e "$rc\n\n"
