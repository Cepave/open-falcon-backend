package utils

import (
	"fmt"
	"github.com/toolkits/ldap"
)

func LdapBind(addr, user, password string) (bool, error) {
	conn, err := ldap.Dial("tcp", addr)
	if err != nil {
		return false, fmt.Errorf("dial ldap fail %s", err.Error())
	}

	defer conn.Close()

	err = conn.Bind(user, password)
	return err == nil, nil
}
