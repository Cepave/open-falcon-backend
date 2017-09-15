#!/bin/bash

function parseOpt
{
	local USAGE='
	go-test-all.sh -t <test folders> -e <exclude folders>

	-t - Folders separated by space
	-e - Exclude folders separated by space

	For example:

	# Tests folder(recursively) "modules/fe" and "modules/hbs"
	# but excludes "modules/fe/ex1" and "modules/fe/ex9"
	go-test-all.sh -t "modules/fe modules/hbs" -e "modules/fe/ex1 modules/fe/ex9"
	'

	local OPTS="t:e:"
	while getopts $OPTS opt; do
		case $opt in
			t) TEST_FOLDER=($OPTARG)
				;;
			e) TEST_FOLDER_EXCLUDE=($OPTARG)
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

TEST_FOLDER=
TEST_FOLDER_EXCLUDE=

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
			echo go test for: ./$folder;

			go test ./$folder
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
	echo "Number of folder has failed for testing: $error_count" >&2
	echo -e "========================================\n"
fi
echo -e "\nNumber of folders has succeeded for testing: $success_count\n"

exit $final_result
