#!/bin/bash

function parseOpt
{
	local USAGE='
	go-test-all.sh -t <test folders> -e <exclude folders>

	-t - Folders separated by space
	-e - Exclude folders separated by space
	-v - Verbose(-test.v -gocheck.vv -ginkgo.v)
	-f - Put "-owl.test.propfile=<file>"
	-p - Put "-owl.test=<properties>"
	-s - Put "-owl.test.sep=<properties>"
	-a - Put "<flags>"(seperated by space) to "go test <flags>"

	For example:

	# Tests folder(recursively) "modules/fe" and "modules/hbs"
	# but excludes "modules/fe/ex1" and "modules/fe/ex9"
	go-test-all.sh -t "modules/fe modules/hbs" -e "modules/fe/ex1 modules/fe/ex9" -p "mysql=root:cepave@tcp(192.16.20.50:3306)/falcon_portal"
	'

	local OPTS="t:e:p:s:f:a:v"
	while getopts $OPTS opt; do
		case $opt in
			a) TEST_FLAGS=(${OPTARG[@]})
				;;
			t) TEST_FOLDER=($OPTARG)
				;;
			e) TEST_FOLDER_EXCLUDE=($OPTARG)
				;;
			v) VERBOSE=true
				;;
			f) FLAG_OWL_TEST_PROPS_FILE=($OPTARG)
				;;
			p) FLAG_OWL_TEST_PROPS="$OPTARG"
				;;
			s) FLAG_OWL_TEST_PROPS_SEP="$OPTARG"
				;;
			*) echo -e "Usage: \n$USAGE" >&2; exit 1
				;;
		esac
	done
}

# Converts the list of folders to exclusions of path used by "find" command:
#
# ! ( -path <path-1> -prune -o -path <path-2> -prune  ... )
function setupFindExcludeSyntax
{
	local var_name=$1
	shift
	local excludes=($@)

	test ${#excludes[@]} -eq 0 && return 0

	local final_syntax=()

	for exclude_path in ${excludes[@]}; do
		if [[ ${#final_syntax[@]} -gt 0 ]]; then
			final_syntax+=("-o")
		fi

		final_syntax+=("-path" "$exclude_path" "-prune")
	done

	final_syntax=("!" "(" "${final_syntax[@]}" ")")

	eval "$var_name=(\"\${final_syntax[@]}\")"
}

function hasFramework
{
	local folder=$1
	local checkPackage=$2

	grep -q -e "$checkPackage" $folder/*_test.go &>/dev/null
}
function setupVerbose
{
	local folder=$1
	local var_name=$2

	eval "$var_name=()"

	test $VERBOSE == "false" && return 0

	local final_args=("-test.v")

	if hasFramework "$folder" "gopkg.in/check.v1"; then
		final_args+=("-gocheck.vv")
	fi
	if hasFramework "$folder" "github.com/onsi/ginkgo"; then
		final_args+=("-ginkgo.v")
	fi

	eval "$var_name=(\"\${final_args[@]}\")"
}
function setupOwlProps
{
	local folder=$1
	local var_name=$2

	cmd_flags=()
	if grep -q -Ee "common/testing/(db|http|jsonrpc)" $folder/*_test.go &>/dev/null; then
		if test -n "$FLAG_OWL_TEST_PROPS_FILE"; then
			cmd_flags+=("-owl.test.propfile=$FLAG_OWL_TEST_PROPS_FILE")
		fi
		if test -n "$FLAG_OWL_TEST_PROPS"; then
			cmd_flags+=("-owl.test=$FLAG_OWL_TEST_PROPS")
		fi
		if test -n "$FLAG_OWL_TEST_PROPS_SEP"; then
			cmd_flags+=("-owl.test.sep=$FLAG_OWL_TEST_PROPS_SEP")
		fi
	fi

	eval "$var_name=(\"\${cmd_flags[@]}\")"
}

TEST_FOLDER=
TEST_FOLDER_EXCLUDE=
TEST_FLAGS=

FLAG_OWL_TEST_PROPS=
FLAG_OWL_TEST_PROPS_FILE=
FLAG_OWL_TEST_PROPS_SEP=

VERBOSE=false

parseOpt "$@"

echo Test folders: ${TEST_FOLDER[@]}
echo -e Exclude folders: ${TEST_FOLDER_EXCLUDE[@]} "\n"

setupFindExcludeSyntax find_exclude_syntax "${TEST_FOLDER_EXCLUDE[@]}"

error_count=0
success_count=0
for go_test_folder in ${TEST_FOLDER[@]}; do
	for folder in `find $go_test_folder -type d "${find_exclude_syntax[@]}"`; do
		if [[ `find $folder -maxdepth 1 -type f -name "*_test.go" | wc -l` -gt 0 ]]; then
			setupVerbose "$folder" verbose_args
			setupOwlProps "$folder" owl_flag_args

			echo go test ./$folder "${verbose_args[@]}" "${owl_flag_args[@]}" "${TEST_FLAGS[@]}"
			go test ./$folder "${verbose_args[@]}" "${owl_flag_args[@]}" "${TEST_FLAGS[@]}"
			TEST_RESULT=$?

			if [[ $TEST_RESULT -eq 0 ]]; then
				(( success_count++ ))
			else
				echo -e "\nTest failed. Code: $TEST_RESULT\n" >&2
				(( error_count++ ))
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

if [[ $error_count -gt 0 ]]; then
	exit 1
fi

exit 0
