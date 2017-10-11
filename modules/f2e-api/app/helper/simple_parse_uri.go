package helper

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Params struct {
}

func (Params) GetInt(params gin.Params, name string) (*int, error) {
	val, exist := params.Get(name)
	if !exist {
		err := fmt.Errorf("missing params: %v", name)
		return nil, err
	}
	intval, err := strconv.Atoi(val)
	return &intval, err
}

func (Params) GetString(params gin.Params, name string) (*string, error) {
	val, exist := params.Get(name)
	if !exist {
		err := fmt.Errorf("missing params: %v", name)
		return nil, err
	}
	return &val, nil
}
