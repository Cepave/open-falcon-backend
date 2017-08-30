if [ "$1" == "init" ]; then
  echo "hostgroup_test.go"
   go test -v hostgroup_test.go --test.run TestHostGroupCreate
fi

echo "hostgroup_test.go"
# go test -v hostgroup_test.go --test.run TestBindHostGroupCreate1
# go test -v hostgroup_test.go --test.run TestBindHostGroupCreate2
# go test -v hostgroup_test.go --test.run TestGetHostGroup
# go test -v hostgroup_test.go --test.run TestUnBindHostGroupCreate1
# go test -v hostgroup_test.go --test.run TestDeleteHostGroup
