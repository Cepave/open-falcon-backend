package types

import (
	"reflect"

	"github.com/juju/errors"
	"github.com/satori/go.uuid"

	t "github.com/Cepave/open-falcon-backend/common/reflect/types"
)

// Adds additional converters:
//
// String to satori/go.uuid.UUID
func AddDefaultConverters(registry ConverterRegistry) {
	registry.AddConverter(
		t.TypeOfString, reflect.TypeOf(uuid.Nil),
		stringToUuid,
	)
}

func stringToUuid(source interface{}) interface{} {
	sourceAsString, ok := source.(string)

	if !ok {
		panic(errors.Details(
			errors.Errorf("Only support string type to uuid.UUID"),
		))
	}

	var targetUuid Uuid
	targetUuid.MustFromString(sourceAsString)

	return uuid.UUID(targetUuid)
}
