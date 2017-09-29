echo "object_tag_create_test.go"
go test -v object_tag_create_test.go --test.run TestObjectTagCreate

echo "object_tag_update_test.go"
go test -v object_tag_update_test.go --test.run TestObjectTagUpdate

echo "object_tag_delete_test.go"
go test -v object_tag_delete_test.go --test.run TestObjectTagDelete
