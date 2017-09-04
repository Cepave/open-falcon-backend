echo "lambda_query_mock_test.go"
go test -v lambda_query_mock_test.go --test.run TestLambdaQuery
echo "func_top_test.go"
go test -v func_top_test.go --test.run TestFuncTop
echo "func_top_diff_test.go"
go test -v func_top_diff_test.go --test.run TestFuncTopDiff
echo "func_avg_compare_test.go"
go test -v func_avg_compare_test.go --test.run TestFuncAvgCompare
