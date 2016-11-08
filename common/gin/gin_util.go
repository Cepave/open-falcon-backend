package gin

import (
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	model "github.com/Cepave/open-falcon-backend/common/model"
	"regexp"
	"strconv"
	"strings"
)

const (
	headerPageSize = "page-size"
	headerPagePos = "page-pos"
	headerTotalCount = "total-count"
	headerPageMore = "page-more"
	headerOrderBy = "order-by"
)

type BindJsonError struct {
	sourceError error
}

func (err BindJsonError) Error() string {
	return err.sourceError.Error()
}

// Binds JSON and panic with JsonBindError if there is error
func BindJson(context *gin.Context, object interface{}) {
	if err := context.BindJSON(object); err != nil {
		panic(BindJsonError{ err })
	}
}

// PagingByHeader would initialize paging object by header
//
// This funcion would load header value:
// "page-size" - The size of page
// "page-pos" - The position of page, starting with "1"
// "order-by" - The order for paging
//
// 		<prop_1>#<dir>:<prop_2>#<dir>:...
//
// context - The gin context
// defaultPaging - The default value of paging
func PagingByHeader(context *gin.Context, defaultPaging *model.Paging) *model.Paging {
	finalPaging := *defaultPaging

	if pageSize := context.Request.Header.Get(headerPageSize)
		pageSize != "" {

		parsedValue, err := strconv.ParseInt(pageSize, 10, 32)
		// If the parsing has error, use the value of default paging
		if err == nil {
			finalPaging.Size = int32(parsedValue)
		}
	}
	if pagePos := context.Request.Header.Get(headerPagePos)
		pagePos != "" {

		parsedValue, err := strconv.ParseInt(pagePos, 10, 32)
		if err == nil {
			finalPaging.Position = int32(parsedValue)
		}
	}
	if orderBy := context.Request.Header.Get(headerOrderBy)
		orderBy != "" {

		parsedValue, err := parseOrderBy(orderBy)
		if err == nil {
			finalPaging.OrderBy = parsedValue
		}
	}

	return &finalPaging
}

// HeaderWithPaging would set headers with information of paging
func HeaderWithPaging(context *gin.Context, paging *model.Paging) {
	context.Header(headerPageSize, int32ToString(paging.Size))
	context.Header(headerPagePos, int32ToString(paging.Position))
	context.Header(headerTotalCount, int32ToString(paging.TotalCount))

	pageMore := "false"
	if paging.PageMore {
		pageMore = "true"
	}
	context.Header(headerPageMore, pageMore)
}

func int32ToString(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

var regexpOrderValue = regexp.MustCompile(`^\w+(?:#(?:a|d|asc|desc|ascending|descending))?(?::\w+(?:#(?:a|d|asc|desc|ascending|descending))?)*$`)
func parseOrderBy(headerValueOfOrderBy string) ([]*model.OrderByEntity, error) {
	var result []*model.OrderByEntity

	headerValueOfOrderBy = strings.ToLower(headerValueOfOrderBy)
	if !regexpOrderValue.MatchString(headerValueOfOrderBy) {
		return result, fmt.Errorf("Header value cannot be recognized. [%v]", headerValueOfOrderBy)
	}

	for _, entity := range strings.Split(headerValueOfOrderBy, ":") {
		if entity == "" {
			continue
		}

		var newOrderEntity = model.OrderByEntity{
			Direction: model.DefaultDirection,
		}
		if strings.Contains(entity, "#") {
			/**
			 * Process <prop_1>#<dir> syntax
			 */
			orderPair := strings.Split(entity, "#")

			newOrderEntity.Expr = orderPair[0]
			newOrderEntity.Direction = getOrderDirection(orderPair[1])
			// :~)
		} else {
			// Process only <prop_1> syntax
			newOrderEntity.Expr = entity
		}

		result = append(result, &newOrderEntity)
	}

	return result, nil
}

func getOrderDirection(value string) byte {
	switch value {
	case "asc", "ascending", "a":
		return model.Ascending
	case "desc", "descending", "d":
		return model.Descending
	}

	return model.DefaultDirection
}
