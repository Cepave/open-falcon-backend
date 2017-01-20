#!/bin/sh

host="localhost"
id=root
pw=password
AUTH=$(echo -u $id:$pw)
DB=extend
MEASUREMENT="service.http.lvs.443.port"

createdb(){
    echo "create database $DB"
    curl -G "http://$host:8086/query" $AUTH --data-urlencode "q=CREATE DATABASE $DB"
}

database(){
    curl -G "http://$host:8086/query?pretty=true" $AUTH --data-urlencode "q=SHOW DATABASES"
}

measurement(){
    curl -G "http://$host:8086/query" $AUTH --data-urlencode "db=$DB" --data-urlencode "q=SHOW MEASUREMENTS"
}

query(){
    curl -G "http://$host:8086/query?pretty=true" $AUTH --data-urlencode "db=$DB" --data-urlencode "q=SELECT * FROM \"$MEASUREMENT\""
}

scope_query(){
    curl -G "http://$host:8086/query?pretty=true" $AUTH --data-urlencode "db=$DB" --data-urlencode "q=SELECT MEAN(system) FROM \"$MEASUREMENT\" WHERE time > now() - 3m GROUP BY time(10s)" 
}

stats() {
    QUERY="SHOW STATS"
    curl -G "http://$host:8086/query?pretty=true" $AUTH --data-urlencode "db=$DB" --data-urlencode "q=$QUERY"
}

diagnostics() {
    QUERY="SHOW DIAGNOSTICS"
    curl -G "http://$host:8086/query?pretty=true" $AUTH --data-urlencode "db=$DB" --data-urlencode "q=$QUERY"
}

usage() {
        echo "$0 createdb"
        echo "$0 database"
        echo "$0 measurement"
        echo "$0 query"
        echo "$0 scope_query"
        echo "$0 stats"
        echo "$0 diagnostics"
}

action=$1
case $action in
    "createdb")
        createdb
        ;;
    "database")
        database
        ;;
    "measurement")
        measurement
        ;;
    "query")
        query
        ;;
    "scope_query")
        scope_query
        ;;
    "stats")
        stats
        ;;
    "diagnostics")
        diagnostics
        ;;
    *)
        usage
        ;;
esac
