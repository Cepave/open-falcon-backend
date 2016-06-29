#!/bin/bash
export PATH=$PATH:/usr/local/bin/
API_ADDR=$1
message=`echo -e "$2\n--------\n$3"|od -t x1 -A n -v -w100000 | tr " " %`
get_gid()
{
    echo curl -s $API_ADDR/openqq/get_group_info
    GID=`curl -s $API_ADDR/openqq/get_group_info | jq '.[0].ginfo.gid'| tr -d '"'`
}
send_messege()
{
    get_gid
    echo curl $API_ADDR/openqq/send_group_message?gid=$GID&content=$message
    api_url="$API_ADDR/openqq/send_group_message?gid=$GID&content=$message"
    curl $api_url
}
send_messege
