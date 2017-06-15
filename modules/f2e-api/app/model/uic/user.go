package uic

import (
	con "github.com/Cepave/open-falcon-backend/modules/f2e-api/config"
	"github.com/spf13/viper"
)

type User struct {
	ID      int64  `json:"id" `
	Name    string `json:"name"`
	Cnname  string `json:"cnname"`
	Passwd  string `json:"-"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	IM      string `json:"im" gorm:"column:im"`
	QQ      string `json:"qq" gorm:"column:qq"`
	Role    int    `json:"role"`
	Creator int    `json:"creator"`
}

func skipAccessControll() bool {
	return !viper.GetBool("access_control")
}

func (this User) IsAdmin() bool {
	if this.Role == 2 || this.Role == 1 {
		return true
	}
	return false
}

func (this User) IsSuperAdmin() bool {
	if this.Role == 2 {
		return true
	}
	return false
}

func (this User) IsThirdPartyUser() bool {
	if this.Creator == 1 {
		return true
	}
	return false
}

func (this User) FindUser() (user User, err error) {
	db := con.Con()
	user = this
	dt := db.Uic.Find(&user)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	return
}

func (this User) UserNameExist() bool {
	db := con.Con()
	if this.Name == "" {
		return false
	}
	counter := 0
	db.Uic.Model(&this).Where("name = ?", this.Name).Count(&counter)
	if counter != 0 {
		return true
	}
	return false
}

func (this User) TableName() string {
	return "user"
}
