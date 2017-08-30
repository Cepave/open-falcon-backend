if [ "$1" == "init" ]; then
  echo "template_create_test.go"
  go test -v template_create_test.go -test.run TestTplCreate
fi
echo "template_get_test.go"
go test -v template_get_test.go -test.run TestTplGet
echo "template_update_test.go"
go test -v template_update_test.go -test.run TestTplUpdate
echo "template_delete_test.go"
go test -v template_delete_test.go -test.run TestTplDelete
echo "template_clone_test.go"
go test -v template_clone_test.go -test.run TestTplClone
