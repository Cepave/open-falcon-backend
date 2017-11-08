#!/bin/bash

SUFFIX=
PREFIX=
databases=(imdb falcon_portal uic falcon_links grafana graph boss dashboard)
mysql_args=()
action=recreate
yes=

current_script=$(basename ${BASH_SOURCE[0]})

help="
${current_script} [--action=recreate] [--mysql=<args>] [--prefix=<db prefix>] [--suffix=<db suffix>]
\n\nThis script would drop and create databases: \"${databases[@]}\"
\n\n\t--action=<action> - The action to be perfomed
\n\n\t\tValue domain: recreate(default), drop, create(if not existing)
\n\n\t--prefix=<prefix> - The prefix to be added to name of database
\n\n\t--suffix=<suffix> - The suffix to be appended to name of database
\n\n\t--mysql=<args> - The arguments to be fed to \"mysql\" command
\n\n\t--yes - Applys \"yes\" for any questions
\n\n\t--help - Show this message
"

function parseParam()
{
	for param in "$@"; do
		case $param in
			--suffix=*)
			SUFFIX=${param#--suffix=}
			;;
			--prefix=*)
			PREFIX=${param#--prefix=}
			;;
			--mysql=*)
			param=${param#--mysql=}
			mysql_args=(${param[@]})
			;;
			--action=*)
			action=${param#--action=}
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
}

function ask_execute()
{
	test $yes -eq 1 && return 0

	echo Databases: "${databases[@]}"
	echo -n "Are you sure to **[$action]** databases(PREFIX=$PREFIX, SUFFIX=$SUFFIX)? [Y/N]: "

	read answer

	answer=$(tr '[:upper:]' '[:lower:]' <<<$answer)

	case $answer in
		y|yes)
		return 0
		;;
	esac

	return 1
}

parseParam "$@"

ask_execute || exit 0

for dbname in "${databases[@]}"; do
	finalDbName="$PREFIX$dbname$SUFFIX"

	case $action in
		recreate)
		finalSql="DROP DATABASE IF EXISTS \`$finalDbName\`; CREATE DATABASE \`$finalDbName\` DEFAULT CHARSET utf8;"
		;;
		drop)
		finalSql="DROP DATABASE IF EXISTS \`$finalDbName\`;"
		;;
		create)
		finalSql="CREATE DATABASE IF NOT EXISTS \`$finalDbName\` DEFAULT CHARSET utf8;"
		;;
	esac

	echo Run mysql: $finalSql

	mysql "${mysql_args[@]}" <<<$finalSql
done
