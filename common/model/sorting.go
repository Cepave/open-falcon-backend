package model

import (
	"fmt"
	"strings"
)

// Deprecated: use utils instead
const (
	DefaultDirection byte = 0
	// Sorting by ascending
	Ascending byte = 1
	// Sorting by descending
	Descending byte = 2
)

// Dialect used to mapping property to column.
// And direction of sorting
type OrderByDialect struct {
	Separator          string
	PropertyMapping    map[string]string
	DirectionMapping   map[byte]string
	FuncEntityToSyntax func(*OrderByEntity) (string, error)
}

func GetOrderByAndLimit(paging *Paging, dialect *OrderByDialect) string {
	orderBy, err := dialect.ToQuerySyntax(paging.OrderBy)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("\nORDER BY %s\n%s", orderBy, GetSqlLimit(paging))
}

func GetSqlLimit(paging *Paging) string {
	return fmt.Sprintf("LIMIT %d, %d", paging.GetOffset(), paging.Size)
}

// Converts the entities of order to Query syntax
//
// If some of mapping could be found, the returned error would be non-nil
func (dialect *OrderByDialect) ToQuerySyntax(entities []*OrderByEntity) (string, error) {
	var querySyntaxForOrderBy []string

	for _, v := range entities {
		syntax, err := dialect.FuncEntityToSyntax(v)
		if err != nil {
			return "", err
		}

		querySyntaxForOrderBy = append(querySyntaxForOrderBy, syntax)
	}

	return strings.Join(querySyntaxForOrderBy, dialect.Separator), nil
}

var sqlDirectionMapping = map[byte]string{
	DefaultDirection: "",
	Ascending:        "ASC",
	Descending:       "DESC",
}

// Builds a dialect for default SQL language
//
// DefaultDirection - nothing(by default value of database)
// Ascending - Omit "ASC"
// Descending - Omit "DESC"
func NewSqlOrderByDialect(propertyMapping map[string]string) *OrderByDialect {
	var newMapOfDirection map[byte]string = make(map[byte]string)
	for k, v := range sqlDirectionMapping {
		newMapOfDirection[k] = v
	}

	var newMapOfProperties map[string]string = make(map[string]string)
	for k, v := range propertyMapping {
		newMapOfProperties[k] = v
	}

	dialect := &OrderByDialect{
		Separator:        ", ",
		PropertyMapping:  newMapOfProperties,
		DirectionMapping: newMapOfDirection,
	}
	dialect.FuncEntityToSyntax = func(entity *OrderByEntity) (string, error) {
		return entityToSqlSyntax(dialect, entity)
	}

	return dialect
}

func entityToSqlSyntax(dialect *OrderByDialect, entity *OrderByEntity) (string, error) {
	var propValue, dirValue string
	var ok bool

	if propValue, ok = dialect.PropertyMapping[entity.Expr]; !ok {
		return "", fmt.Errorf("Cannot find mapping for property: [%s]", entity.Expr)
	}
	if dirValue, ok = dialect.DirectionMapping[entity.Direction]; !ok {
		return "", fmt.Errorf("Cannot find mapping for direction: [%d]", entity.Direction)
	}

	if dirValue == "" {
		return propValue, nil
	}

	return fmt.Sprintf("%s %s", propValue, dirValue), nil
}

// The order by
type OrderByEntity struct {
	// Could name of column, property or any user-defined text
	Expr string
	// See Asc/Desc constant
	Direction byte
}

func (entity *OrderByEntity) String() string {
	var direction = "<DEFAULT>"
	switch entity.Direction {
	case Ascending:
		direction = "Asc"
	case Descending:
		direction = "Desc"
	}

	return fmt.Sprintf("OrderBy: [%s # %s]", entity.Expr, direction)
}
