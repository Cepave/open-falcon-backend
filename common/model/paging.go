package model

import (
	"fmt"
	"reflect"
)

// The paging object used to hold information
type Paging struct {
	Size       int32
	Position   int32
	TotalCount int32
	PageMore   bool
	OrderBy    []*OrderByEntity
}

// Initialize the paging with default values
//
// 	Paging.Size: -1
// 	Paging.Position: -1
// 	Paging.TotalCount: -1
// 	Paging.PageMore: false
func NewUndefinedPaging() *Paging {
	return &Paging{
		Size:       -1,
		Position:   -1,
		TotalCount: -1,
		PageMore:   false,
	}
}

// Gets the offset to be skipped
func (self *Paging) GetOffset() int32 {
	if self.Position <= 1 {
		return 0
	}

	return (self.Position - 1) * self.Size
}

// Shows the information of paging
func (self *Paging) String() string {
	return fmt.Sprintf(
		"Page Size:[%d]. Page Position:[%d]. Total Count:[%d]. Has More:[%v]. Order By: [%v]",
		self.Size, self.Position, self.TotalCount, self.PageMore, self.OrderBy,
	)
}

// Sets the total count of data and set-up PageMore flag
func (self *Paging) SetTotalCount(totalCount int32) {
	self.TotalCount = totalCount
	self.PageMore = self.GetOffset()+self.Size < totalCount
}

// Extracts page from array or slice
//
// This function would call Paging.SetTotalCount function on input paging
func ExtractPage(arrayObject interface{}, paging *Paging) interface{} {
	arrayValue := reflect.ValueOf(arrayObject)
	if arrayValue.Kind() != reflect.Array &&
		arrayValue.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Input value should be array or slice, got %T", arrayObject))
	}

	resultSlice := reflect.MakeSlice(arrayValue.Type(), 0, 0)

	/**
	 * Calculates the boundary of indexes on input array
	 */
	startIndex, endIndex := int(paging.GetOffset()), int(paging.GetOffset()+paging.Size)
	lenOfArray := arrayValue.Len()
	if endIndex > lenOfArray {
		endIndex = lenOfArray
	}
	// :~)

	for i := startIndex; i < endIndex; i++ {
		resultSlice = reflect.Append(resultSlice, arrayValue.Index(i))
	}

	return resultSlice.Interface()
}
