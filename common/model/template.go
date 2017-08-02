package model

import (
	"fmt"
)

type Template struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentId"`
	ActionId int    `json:"actionId"`
	Creator  string `json:"creator"`
}

func (this *Template) String() string {
	return fmt.Sprintf(
		"<Id:%d, Name:%s, ParentId:%d, ActionId:%d, Creator:%s>",
		this.Id,
		this.Name,
		this.ParentId,
		this.ActionId,
		this.Creator,
	)
}

type NewTemplate struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
	ActionID int    `json:"action_id"`
	Creator  string `json:"creator"`
}

func (this *NewTemplate) String() string {
	return fmt.Sprintf(
		"<ID:%d, Name:%s, ParentID:%d, ActionID:%d, Creator:%s>",
		this.ID,
		this.Name,
		this.ParentID,
		this.ActionID,
		this.Creator,
	)
}
