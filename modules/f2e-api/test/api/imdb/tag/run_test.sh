echo "tag_test.go"
go test -v tag_test.go --test.run TestTagCreate

go test -v tag_test.go --test.run TestTagDelete
