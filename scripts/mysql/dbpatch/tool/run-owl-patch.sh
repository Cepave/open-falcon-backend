#!/bin/bash

PATCH_BIN=
PATCH_LOG_BASE=../change-log/
PATCH_DATABASE=

DATABASE_CONNECTION="root:cepave@tcp(192.168.20.50:3306)"
DATABASE_NAME=
DATABASE_TYPE=mysql

PATCH_YAML_FILE=
PATCH_CHANGLOG_BASE=

bin_paran=

bin_os=
bin_arch=32

HELP="
Options:\n\n
-bin=<bin file> - default value is \"dbpatch\"\n
-bin-os=<os for bin file> - If this value is provided, ignores \"-bin\" option\n
\t linux, windows, or osx\n
-bin-arch=<32> - If this value is provided, ignores \"-bin\" option\n
\t 64 or 32\n
-database=<boss|portal|uic|links|graph|grafana|dashboard|imdb>\n
-log-base=<base directory>. Defualt value is \"$PATCH_LOG_BASE\"\n
-db-connection=<connection string> - default value is \"$DATABASE_CONNECTION\"\n
-db-type=<type> - default value is \"mysql\"\n
"

function load_params()
{
	if [[ $# -eq 0 ]]; then
		echo -e $HELP
		exit 1
	fi

	while [[ $# -gt 0 ]]; do
		key="$1"

		case $key in
			-bin=*)
			bin_param=${key#-bin=}
			;;
			-bin-os=*)
			bin_os=dbpatch-${key#-bin-os=}
			;;
			-bin-arch=*)
			bin_arch=${key#-bin-arch=}
			;;
			-log-base=*)
			PATCH_LOG_BASE=${key#-log-base=}
			;;
			-database=*)
			PATCH_DATABASE=${key#-database=}
			;;
			-db-connection=*)
			DATABASE_CONNECTION=${key#-db-connection=}
			;;
			-db-type=*)
			DATABASE_TYPE=${key#-db-type=}
			;;
			*)
			>&2 echo "Unknown option: \"$key\""
			exit 1
			;;
		esac

		shift
	done

	find_bin
	check_logbase
	check_database
}
function check_logbase()
{
	if [ -z $PATCH_LOG_BASE ]; then
		>&2 echo "Needs \"-log-base=<directory>\""
		exit 1
	fi

	if [ ! -d $PATCH_LOG_BASE ]; then
		>&2 echo "Log base \"$PATCH_LOG_BASE\" is not a directory"
		exit 1
	fi
}
function check_database()
{
	case $PATCH_DATABASE in
		boss)
		DATABASE_NAME=boss
		;;
		portal)
		DATABASE_NAME=falcon_portal
		;;
		uic)
		DATABASE_NAME=uic
		;;
		links)
		DATABASE_NAME=falcon_links
		;;
		graph)
		DATABASE_NAME=graph
		;;
		grafana)
		DATABASE_NAME=grafana
		;;
		dashboard)
		DATABASE_NAME=dashboard
		;;
		imdb)
		DATABASE_NAME=imdb
		;;
		*)
		>&2 echo "Needs \"-database=<boss|portal|uic|links|graph|grafana|dashboard|imdb>\""
		exit 1
		;;
	esac

	DATABASE_CONNECTION="$DATABASE_CONNECTION/$DATABASE_NAME"

	PATCH_YAML_FILE="$PATCH_LOG_BASE/change-log-$PATCH_DATABASE.yaml"
	PATCH_CHANGLOG_BASE="$PATCH_LOG_BASE/schema-$PATCH_DATABASE"

	if [ ! -f $PATCH_YAML_FILE ]; then
		echo "The yaml file \"$PATCH_YAML_FILE\" is not viable"
		exit 1
	fi

	if [ ! -d $PATCH_CHANGLOG_BASE ]; then
		echo "The directory \"$PATCH_CHANGLOG_BASE\" is not viable"
		exit 1
	fi
}
function find_bin()
{
	test -n "$bin_os" && PATCH_BIN="$bin_os-$bin_arch"
	test -n $bin_param && test -z $PATCH_BIN && PATCH_BIN=$bin_param

	if [ -n $PATCH_BIN ] && ! [ -e $PATCH_BIN ] ; then
		>&2 echo "File: \"$PATCH_BIN\" is not executable"
		exit 1
	fi

	if [ -e "dbpatch" ]; then
		PATCH_BIN="./dbpatch"
	fi

	if [ -e "dbpatch.exe" ]; then
		PATCH_BIN="./dbpatch.exe"
	fi

	if [ -z $PATCH_BIN ]; then
		>&2 echo "Need set -bin=<file>"
		exit 1
	fi
}

function ask_execute()
{
	echo
	echo -n "Are you sure to patch? [Y/N]: "

	read answer

	answer=$( tr '[:upper:]' '[:lower:]' <<<$answer )

	case $answer in
		y|yes)
		return 1
		;;
	esac

	return 0
}

load_params ${BASH_ARGV[@]}

echo "Patching database: \"$DATABASE_CONNECTION\""
echo "Using bin: ./\"$PATCH_BIN\""
echo "Change log: \"$PATCH_LOG_BASE\""

ask_execute

if [[ $? == 0 ]]; then
	echo "Goodbye, your database is crashed... HA~HA~HA~"
	exit 0
fi

echo
./$PATCH_BIN "-driverName=$DATABASE_TYPE" "-dataSourceName=$DATABASE_CONNECTION" "-changeLog=$PATCH_YAML_FILE" "-patchFileBase=$PATCH_CHANGLOG_BASE"
