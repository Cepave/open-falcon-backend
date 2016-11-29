package model

import (
	"fmt"
)

// The paging object used to hold information
type Paging struct {
	Size int32
	Position int32
	TotalCount int32
	PageMore bool
	OrderBy []*OrderByEntity
}

// Initialize the paging with default values
//
// 	Paging.Size: -1
// 	Paging.Position: -1
// 	Paging.TotalCount: -1
// 	Paging.PageMore: false
func NewUndefinedPaging() *Paging {
	return &Paging{
		Size: -1,
		Position: -1,
		TotalCount: -1,
		PageMore: false,
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
	return fmt.Sprintf("Page Size:[%d]. Page Position:[%d]. Total Count:[%d]. Has More:[%v]. Order By: [%v]", self.Size, self.Position, self.TotalCount, self.PageMore, self.OrderBy)
}
