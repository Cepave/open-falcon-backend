package falcon_portal

import (
	"fmt"
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/model/uic"
	"github.com/Cepave/open-falcon-backend/modules/f2e-api/app/utils"
	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
)

////////////////////////////////////////////////////////////////////////////////////
// |id                    | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
// | uic                  | varchar(255)     | NO   |     |         |                |
// | url                  | varchar(255)     | NO   |     |         |                |
// | callback             | tinyint(4)       | NO   |     | 0       |                |
// | before_callback_sms  | tinyint(4)       | NO   |     | 0       |                |
// | before_callback_mail | tinyint(4)       | NO   |     | 0       |                |
// | after_callback_sms   | tinyint(4)       | NO   |     | 0       |                |
// | after_callback_mail  | tinyint(4)       | NO   |     | 0  		  |								 |
////////////////////////////////////////////////////////////////////////////////////
type Action struct {
	ID                 int64  `json:"id" gorm:"column:id"`
	UIC                string `json:"uic" gorm:"column:uic"`
	URL                string `json:"url" gorm:"column:url"`
	Callback           int    `json:"callback" orm:"column:callback"`
	BeforeCallbackSMS  int    `json:"before_callback_sms" orm:"column:before_callback_sms"`
	BeforeCallbackMail int    `json:"before_callback_mail" orm:"column:before_callback_mail"`
	AfterCallbackSMS   int    `json:"after_callback_sms" orm:"column:after_callback_sms"`
	AfterCallbackMail  int    `json:"after_callback_mail" orm:"column:after_callback_mail"`
}

func (this Action) TableName() string {
	return "action"
}

func (self Action) FindUics() (teams []uic.Team) {
	teams = []uic.Team{}
	db := con.Con()
	if self.UIC != "" {
		teamst := strings.Split(self.UIC, ",")
		teamstCt, _ := utils.ArrStringsToString(teamst)
		db.Uic.Model(teams).Where(fmt.Sprintf("name IN (%v)", teamstCt)).Scan(&teams)
	}
	return teams
}
