#!/bin/bash

if [[ ${BASH_VERSION:0:1} -lt 4 ]]; then
	echo "Need version of BASH to be at least \"4.x\"" >&2
	exit 1
fi

databases=(uic falcon_portal falcon_links grafana graph boss dashboard imdb)

declare -A dbAndChangelog
dbAndChangelog[uic]="uic"
dbAndChangelog[falcon_portal]="portal"
dbAndChangelog[imdb]="imdb"
dbAndChangelog[falcon_links]="links"
dbAndChangelog[grafana]="grafana"
dbAndChangelog[graph]="graph"
dbAndChangelog[boss]="boss"
dbAndChangelog[dashboard]="dashboard"

URL=
PREFIX=
SUFFIX=

liquibase_options=()
liquibase_command=(update)
mysql_conn=192.168.20.50:3306
change_log_base="./liquibase-changelog/"
yes=

java_property=()
current_script=$(basename ${BASH_SOURCE[0]})

help="
${current_script} [--mysql_conn=192.168.20.50:3306] [--command=update] [--change-log-base=./liquibase-changelog/] [--options=<args>] [--database=<database>] [--prefix=<db prefix>] [--suffix=<db suffix>]
\n\nThis script would update databases(By Liquibase): \"${databases[@]}\"
\n\n\t--mysql_conn=<host:port> - host and port for database connection of MySql.
\n\n\t\tFor example: 192.168.20.50:3306(default value)
\n\n\t--command=<command> - The command of Liquibase
\n\n\t\tFor example: \"update\", \"updateCount 13\"
\n\n\t--database=<database> - Name of OWL database.
\n\n\t\tValue domain: \"${databases[@]}\"
\n\n\t--options=<args> - Arguments of liquibase.
\n\n\t\tFor example: \"--username=abc --password=cepave\"
\n\n\t\tDefault value: \"update\"
\n\n\t--change-log-base=<directory> - Base directory of files for changelog
\n\n\t\tDefault value: \"./liquibase-changelog/\"
\n\n\t--prefix=<prefix> - The prefix to be added to name of database
\n\n\t--suffix=<suffix> - The suffix to be appended to name of database
\n\n\t--yes - Applys \"yes\" to any question
\n\n\t--help - Show this message
"

function parseParam()
{
	for param in "$@"; do
		case $param in
			--change-log-base=*)
			change_log_base=${param#--change-log-base=}
			;;
			--mysql_conn=*)
			mysql_conn=${param#--mysql_conn=}
			;;
			--options=*)
			param=${param#--options=}
			liquibase_options=(${param[@]})
			;;
			--suffix=*)
			SUFFIX=${param#--suffix=}
			;;
			--prefix=*)
			PREFIX=${param#--prefix=}
			;;
			--database=*)
			databases=(${param#--database=})
			;;
			--command=*)
			param=${param#--command=}
			liquibase_command=(${param[@]})
			;;
			--yes)
			yes=1
			;;
			--help)
			echo -e $help
			exit 0
			;;
			*)
			echo "Unknown parameter: $param" >&2
			exit 1
			;;
		esac
	done

	java_property=("-Ddbname.portal=$PREFIX"falcon_portal"$SUFFIX" "-Ddbname.uic=$PREFIX"uic"$SUFFIX")
}

function ask_execute()
{
	[[ $yes -eq 1 ]] && return 0

	echo Databases: ${databases[@]}.
	echo Command: ${liquibase_command[@]}
	echo MySql Host: $mysql_conn
	echo -n "Are you sure to perform command on databases by Liquibase(PREFIX=$PREFIX, SUFFIX=$SUFFIX)? [Y/N]: "

	read answer

	case $(tr '[:upper:]' '[:lower:]' <<<$answer) in
		y|yes)
		return 0
		;;
	esac

	return 1
}

parseParam "$@"

ask_execute || exit 0

for dbname in "${databases[@]}"; do
	test -z "${dbAndChangelog[$dbname]}" && { echo "No such database: $dbname" >&2; exit 1; }

	finalDbName="$PREFIX$dbname$SUFFIX"
	finalUrl="jdbc:mysql://$mysql_conn/$finalDbName?useSSL=false&loc=Local&parseTime=true&characterEncoding=utf8"
	changeLogFile="${dbAndChangelog[$dbname]}.yaml"

	echo "Target Database: \"$finalDbName\". Change log file: \"$changeLogFile\"."

	./liquibase/liquibase "--url=$finalUrl" "--changeLogFile=$change_log_base$changeLogFile" \
		--databaseChangeLogTableName=lq_change_log --databaseChangeLogLockTableName=lg_lock \
			"${liquibase_options[@]}" "${liquibase_command[@]}" "${java_property[@]}"

	result=$?
	if [[ $result -ne 0 ]]; then
		echo "Update database: $finalDbName has error. Code: $?" >&2
		exit $result
	fi
done
