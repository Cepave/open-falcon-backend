package owl

type NameTag struct {
	Id int16 `json:"id" db:"nt_id"`
	Value string `json:"value" db:"nt_value"`
}
