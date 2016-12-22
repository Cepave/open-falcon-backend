package gin

import (
	"gopkg.in/gin-gonic/gin.v1"
	"strconv"
)

// QueryWrapper for gin.Context
type QueryWrapper gin.Context

// NewQueryWrapper converts *gin.Context to *QueryWrapper
func NewQueryWrapper(context *gin.Context) *QueryWrapper {
	return (*QueryWrapper)(context)
}

// GetInt64 gets query parameter as int64
func (wrapper *QueryWrapper) GetInt64(key string) *ViableParamValue {
	return wrapper.getValueWithParseFunc(
		key, valueToInt,
	)
}
// GetInt64Default gets query parameter as int64 with default value
func (wrapper *QueryWrapper) GetInt64Default(key string, defaultValue int64) *ViableParamValue {
	return wrapper.getValueDefaultWithParseFunc(
		key, defaultValue, valueToInt,
	)
}

// GetUint64 gets query parameter as uint64
func (wrapper *QueryWrapper) GetUint64(key string) *ViableParamValue {
	return wrapper.getValueWithParseFunc(
		key, valueToUInt,
	)
}
// GetUint64Default gets query parameter as uint64 with default value
func (wrapper *QueryWrapper) GetUint64Default(key string, defaultValue uint64) *ViableParamValue {
	return wrapper.getValueDefaultWithParseFunc(
		key, defaultValue, valueToUInt,
	)
}

// GetFloat64 gets query parameter as uint64
func (wrapper *QueryWrapper) GetFloat64(key string) *ViableParamValue {
	return wrapper.getValueWithParseFunc(
		key, valueToFloat64,
	)
}
// GetFloat64Default gets query parameter as uint64 with default value
func (wrapper *QueryWrapper) GetFloat64Default(key string, defaultValue float64) *ViableParamValue {
	return wrapper.getValueDefaultWithParseFunc(
		key, defaultValue, valueToFloat64,
	)
}

// GetBool gets query parameter as bool
//
// It accepts:
// 1, t, T, TRUE, true, True,
// 0, f, F, FALSE, false, False.
func (wrapper *QueryWrapper) GetBool(key string) *ViableParamValue {
	return wrapper.getValueWithParseFunc(
		key, valueToBool,
	)
}
// GetBoolDefault gets query parameter as bool with default value
func (wrapper *QueryWrapper) GetBoolDefault(key string, defaultValue bool) *ViableParamValue {
	return wrapper.getValueDefaultWithParseFunc(
		key, defaultValue, valueToBool,
	)
}

/**
 * Private function/objects
 */
type parseFunc func(value string) (interface{}, error)

func valueToInt(value string) (interface{}, error) {
	return strconv.ParseInt(value, 10, 64)
}
func valueToUInt(value string) (interface{}, error) {
	return strconv.ParseUint(value, 10, 64)
}
func valueToFloat64(value string) (interface{}, error) {
	return strconv.ParseFloat(value, 64)
}
func valueToBool(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}

func (wrapper *QueryWrapper) getValueWithParseFunc(key string, parseImpl parseFunc) *ViableParamValue {
	viableParam := &ViableParamValue {
		Viable: true,
		Value: ((*gin.Context)(wrapper)).Query(key),
		Error: nil,
	}

	if viableParam.Value.(string) == "" {
		viableParam.Viable = false
		viableParam.Value = nil
	} else {
		viableParam.Value, viableParam.Error = parseImpl(viableParam.Value.(string))
		if viableParam.Error != nil {
			viableParam.Viable = false
			viableParam.Value = nil
		}
	}

	return viableParam
}

func (wrapper *QueryWrapper) getValueDefaultWithParseFunc(key string, defaultValue interface{}, parseImpl parseFunc) *ViableParamValue {
	viableParam := wrapper.getValueWithParseFunc(key, parseImpl)
	if viableParam.Error != nil {
		return viableParam
	}
	if !viableParam.Viable {
		viableParam.Value = defaultValue
	}

	return viableParam
}
