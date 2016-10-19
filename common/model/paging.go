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

func (self *Paging) String() string {
	return fmt.Sprintf("Page Size:[%d]. Page Position:[%d]. Total Count:[%d]. Has More:[%v].", self.Size, self.Position, self.TotalCount, self.PageMore)
}
