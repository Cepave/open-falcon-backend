echo "value_model_test.go"
go test -v value_model_test.go --test.run TestValueModelCreate
go test -v value_model_test.go --test.run TestValueModelUpdate
go test -v value_model_test.go --test.run TestValueModelDelete
go test -v value_model_test.go --test.run TestValueModelGet
