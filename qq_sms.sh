#!/bin/bash
export PATH=$PATH:/usr/local/bin/
# 你login.pl中定义的host和port
API_ADDR=$1
# 处理下编码，用于合并告警内容的标题和内容，即$2和$3
#message=`echo -e "$1\n$2"|od -t x1 -A n -v -w1000000000 | tr " " %`
message=`echo -e "$2\n--------\n$3"|od -t x1 -A n -v -w100000 | tr " " %`
#——————- main body dont’ modify blow ——————–#
get_gid()
{
    # 获取gid，需要安装jd处理json
    GID=`curl -s http://$API_ADDR/openqq/get_group_info | jq '.[0].ginfo.gid'| tr -d '"'`
}
send_messege()
{
    get_gid
    # 这没什么好说的
    api_url="http://$API_ADDR/openqq/send_group_message?gid=$GID&content=$message"
    curl $api_url
}
# 发送消息，执行函数，很多朋友复制的时候漏了
send_messege
