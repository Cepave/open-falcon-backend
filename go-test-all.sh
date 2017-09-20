#!/bin/bash

function parseOpt
{
	local USAGE='
	go-test-all.sh -t <test folders> -e <exclude folders>

	-t - Folders separated by space
	-e - Exclude folders separated by space
	-v - Verbose(-test.v -gocheck.vv -ginkgo.v)
	-o - Put "-owl.test=<properties>"
	-s - Put "-owl.test.sep=<properties>"

	For example:

	# Tests folder(recursively) "modules/fe" and "modules/hbs"
	# but excludes "modules/fe/ex1" and "modules/fe/ex9"
	go-test-all.sh -t "modules/fe modules/hbs" -e "modules/fe/ex1 modules/fe/ex9" -o "mysql=root:cepave@tcp(192.16.20.50:3306)/falcon_portal"
	'

	local OPTS="t:e:o:s:v"
	while getopts $OPTS opt; do
		case $opt in
			t) TEST_FOLDER=($OPTARG)
				;;
			e) TEST_FOLDER_EXCLUDE=($OPTARG)
				;;
			v) VERBOSE="-test.v"
				;;
			o) OWL_TEST_PROPS="$OPTARG"
				;;
			s) OWL_TEST_SEP="$OPTARG"
				;;
			*) echo -e "Usage: \n$USAGE" >&2; exit 1
				;;
		esac
	done
}

# Converts the list of folders to exclusions of path used by "find" command:
#
# ! ( -path <path-1> -prune -o -path <path-2> -prune  ... )
function buildFindExcludeFolder
{
	local EXCLUDES=($@)

	if [[ ${#EXCLUDES[*]} -eq 0 ]]; then
		return 0
	fi

	local COMBINED_EXCLUDE=

	for exclude_path in ${EXCLUDES[*]}; do
		local path_syntax="-path $exclude_path -prune"

		local OR_CONN=" -o "
		if [[ ${#COMBINED_EXCLUDE} -eq 0 ]]; then
			OR_CONN=""
		fi
		COMBINED_EXCLUDE+="$OR_CONN$path_syntax"
	done

	echo ! \( $COMBINED_EXCLUDE \)
}

function frameworkVerbose
{
	local existingVerbose=$1
	local folder=$2
	local checkPackage=$3
	local frameworkFlag=$4

	if grep -q -e "$checkPackage" $folder/*_test.go; then
		echo "$existingVerbose" "$frameworkFlag"
	else
		echo "$existingVerbose"
	fi
}
function getVerbose
{
	test -z $VERBOSE && return 0

	local folder=$1
	local currentVerbose="$VERBOSE"

	currentVerbose=$(frameworkVerbose "$currentVerbose" "$folder" "gopkg.in/check.v1" "-gocheck.vv")
	currentVerbose=$(frameworkVerbose "$currentVerbose" "$folder" "github.com/onsi/ginkgo" "-ginkgo.v")

	echo $currentVerbose
}

TEST_FOLDER=
TEST_FOLDER_EXCLUDE=
VERBOSE=
OWL_TEST_PROPS=
OWL_TEST_SEP=

parseOpt "$@"

echo Test folders: ${TEST_FOLDER[*]}
echo -e Exclude folders: ${TEST_FOLDER_EXCLUDE[*]} "\n"

exclude_syntax=$(buildFindExcludeFolder ${TEST_FOLDER_EXCLUDE[*]})

final_result=0
error_count=0
success_count=0
for go_test_folder in ${TEST_FOLDER[*]}; do
	for folder in `find $go_test_folder -type d $exclude_syntax`; do
		if [[ `find $folder -maxdepth 1 -type f -name "*_test.go" | wc -l` -gt 0 ]]; then
			current_verbose=$(getVerbose $folder)
			echo go test ./$folder $current_verbose
			go test ./$folder $current_verbose
			TEST_RESULT=$?

			if [[ $TEST_RESULT -eq 0 ]]; then
				(( success_count++ ))
			else
				echo -e "\nTest failed. Code: $TEST_RESULT\n" >&2
				(( error_count++ ))
				final_result=1
			fi
		fi
	done
done

if [[ $error_count -gt 0 ]]; then
	echo -e "\n========================================"
	echo "Number of folders has failed for testing: $error_count" >&2
	echo -e "========================================\n"
fi
echo -e "\nNumber of folders has succeeded for testing: $success_count\n"

exit $final_result
