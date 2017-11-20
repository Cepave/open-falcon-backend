#!/bin/bash

if [[ ${BASH_VERSION:0:1} -lt 4 ]]; then
	echo "Need version of BASH to be at least \"4.x\"" >&2
	exit 1
fi

wait_time=${1:=3s}
max_times=${2:=6}

counter=0

echo -n "Waiting(\"$wait_time\" with maximum $max_times times) for MySql($port) to be ready ... "
while ( mysql --connect_timeout=2 -h $mysql_host -P $mysql_port -u$mysql_user -p$mysql_password <<<"SELECT NOW();" &>mysql.ping.log; test $? -ne 0 ); do
	let counter++

	test $counter -ge $max_times && { cat mysql.ping.log; docker logs ci-mysql; exit 1; }

	sleep $wait_time
done

echo "Ready."
