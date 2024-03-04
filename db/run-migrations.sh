#!/bin/bash

cd /
echo "Migration status"
neo4j-migrations -v -p neo4j -a neo4j://neo4j:7687 --location file://migrations info
echo "Running migrations"
neo4j-migrations -v -p neo4j -a neo4j://neo4j:7687 --location file://migrations migrate
echo "Migrations ran with return code $?"
exit $?