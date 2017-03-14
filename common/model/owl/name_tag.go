package owl

type NameTag struct {
	Id    int16  `json:"id" db:"nt_id"`
	Value string `json:"value" db:"nt_value"`
}

func (NameTag) TableName() string {
	return "owl_name_tag"
}

type NameTagOfPingtaskView struct {
	Id    int    `json:"id" db:"nt_id"`
	Value string `json:"value" db:"nt_value"`
}

func (NameTagOfPingtaskView) TableName() string {
	return "owl_name_tag"
}
