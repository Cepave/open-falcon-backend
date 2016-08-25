#!/bin/bash

declare -A confs
confs=(
    [%%PORTAL%%]=10.20.30.40
    [%%DASHBOARD%%]=10.20.30.40
    [%%GRAFANA%%]=10.20.30.40
    [%%ALARM%%]=10.20.30.40
    [%%FE_HTTP%%]=0.0.0.0:1234
    [%%QUERY_HTTP%%]=0.0.0.0:9966
    [%%GIN_HTTP%%]=0.0.0.0:9967
    [%%HBS_HTTP%%]=0.0.0.0:6031
    [%%TRANSFER_HTTP%%]=0.0.0.0:6060
    [%%GRAPH_HTTP%%]=0.0.0.0:6071
    [%%HBS_RPC%%]=0.0.0.0:6030
    [%%REDIS%%]=172.17.0.2:6379
    [%%GRAPH_RPC%%]=0.0.0.0:6070
    [%%TRANSFER_RPC%%]=0.0.0.0:8433
    [%%FE%%]=127.0.0.1:1235
    [%%MYSQL%%]="root:@tcp(172.17.0.3:3306)"
)

configurer() {
    for i in "${!confs[@]}"
    do
        search=$i
        replace=${confs[$i]}
        # Note the "" after -i, needed in OS X
        sed -i "s/${search}/${replace}/g" ./config/*.json
    done
}
configurer
