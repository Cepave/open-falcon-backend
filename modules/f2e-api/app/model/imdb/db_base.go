package imdb

import (
	"fmt"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

var db = config.Con().IMDB

func Init() {
	var tag_type TagType
	var tag Tag
	db.Model(&tag_type).Related(&tag, "TagTypeId")
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
