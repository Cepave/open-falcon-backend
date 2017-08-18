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

	uuidValue, err := uuid.FromString(sourceAsString)
	if err != nil {
		panic(errors.Details(
			errors.Annotatef(err, "Value of \"%s\" cannot be parsed to uuid.UUID", sourceAsString),
		))
	}

	return uuidValue
}
