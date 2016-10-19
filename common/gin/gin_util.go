package gin

import (
	"gopkg.in/gin-gonic/gin.v1"
	model "github.com/Cepave/open-falcon-backend/common/model"
	"strconv"
)

const (
	headerPageSize = "page-size"
	headerPagePos = "page-pos"
	headerTotalCount = "total-count"
	headerPageMore = "page_more"
)

// PagingByHeader would initialize paging object by header
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

	return &finalPaging
}

// HeaderWithPaging would set headers with information of paging
func HeaderWithPaging(context *gin.Context, defaultPaging *model.Paging) {
	context.Header(headerPageSize, string(defaultPaging.Size))
	context.Header(headerPagePos, string(defaultPaging.Position))
	context.Header(headerTotalCount, string(defaultPaging.TotalCount))

	pageMore := "false"
	if defaultPaging.PageMore {
		pageMore = "true"
	}
	context.Header(headerPageMore, pageMore)
}
