#!/bin/bash

DATABASE_NAME=
MYSQLDUMP_OPTIONS=

function load_params()
{
	if [[ $# -eq 0 ]]; then
		echo "Usage export-changelog.sh -db-name=<database_name> {{ mysqldump options... }}"
		echo Options:
		echo
		echo "-db-name=<database_name>"
		exit 1
	fi

	while [[ $# -gt 0 ]]; do
		key="$1"

		case $key in
			-db-name=*)
			DATABASE_NAME=${key#-db-name=}
			;;
			*)
			MYSQLDUMP_OPTIONS="$MYSQLDUMP_OPTIONS $key"
			;;
		esac
		shift
	done
}

load_params ${BASH_ARGV[@]}

if [ -z $DATABASE_NAME ]; then
	echo "Needs -db-name=<database name>"
	exit 1
fi

mysqldump $MYSQLDUMP_OPTIONS $DATABASE_NAME sysdb_change_log
