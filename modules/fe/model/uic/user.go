package uic

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/toolkits/cache"
	"github.com/toolkits/logger"
	"github.com/toolkits/slice"
	"time"
)

func SelectUserById(id int64) *User {
	if id <= 0 {
		return nil
	}

	obj := User{Id: id}
	err := orm.NewOrm().Read(&obj, "Id")
	if err != nil {
		if err != orm.ErrNoRows {
			logger.Errorln(err)
		}
		return nil
	}
	return &obj
}

func ReadUserById(id int64) *User {
	if id <= 0 {
		return nil
	}

	key := fmt.Sprintf("user:obj:%d", id)
	var obj User
	if err := cache.Get(key, &obj); err != nil {
		objPtr := SelectUserById(id)
		if objPtr != nil {
			go cache.Set(key, objPtr, time.Hour)
		}
		return objPtr
	}

	return &obj
}

func SelectUserIdByName(name string) int64 {
	if name == "" {
		return 0
	}

	type IdStruct struct {
		Id int64
	}

	var idObj IdStruct
	err := orm.NewOrm().Raw("select id from user where name = ?", name).QueryRow(&idObj)
	if err != nil {
		return 0
	}

	return idObj.Id
}

func ReadUserIdByName(name string) int64 {
	if name == "" {
		return 0
	}

	key := fmt.Sprintf("user:id:%s", name)
	var id int64
	if err := cache.Get(key, &id); err != nil {
		id = SelectUserIdByName(name)
		if id > 0 {
			go cache.Set(key, id, time.Hour)
		}
	}

	return id
}

func ReadUserByName(name string) *User {
	return ReadUserById(ReadUserIdByName(name))
}

func (this *User) Save() (int64, error) {
	id, err := orm.NewOrm().Insert(this)
	if err != nil {
		this.Id = id
	}
	return id, err
}

func InsertRegisterUser(name, password, email string) (int64, error) {
	userPtr := &User{
		Name:   name,
		Passwd: password,
		Email:  email,
	}
	return userPtr.Save()
}

func (this *User) Update() (int64, error) {
	num, err := orm.NewOrm().Update(this)

	if err == nil && num > 0 {
		cache.Delete(fmt.Sprintf("user:obj:%d", this.Id))
	}

	return num, err
}

func (this *User) CanWrite(t *Team) bool {
	if this.Role > 0 {
		return true
	}

	uids, err := Uids(t.Id)
	if err != nil {
		return false
	}

	return slice.ContainsInt64(uids, this.Id)
}

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(User))
}

func QueryUsers(query string) orm.QuerySeter {
	qs := orm.NewOrm().QueryTable(new(User))
	if query != "" {
		cond := orm.NewCondition()
		cond = cond.Or("Name__icontains", query).Or("Email__icontains", query)
		qs = qs.SetCond(cond)
	}
	return qs
}

func (this *User) Remove() (int64, error) {
	num, err := DeleteUserById(this.Id)
	if err != nil {
		return num, err
	}

	cache.Delete(fmt.Sprintf("user:obj:%d", this.Id))
	cache.Delete(fmt.Sprintf("user:id:%s", this.Name))

	tids, err := Tids(this.Id)
	if err == nil && len(tids) > 0 {
		for i := 0; i < len(tids); i++ {
			cache.Delete(fmt.Sprintf("t:uids:%d", tids[i]))
		}
	}

	UnlinkByUserId(this.Id)

	return num, err
}

func DeleteUserById(id int64) (int64, error) {
	r, err := orm.NewOrm().Raw("DELETE FROM `user` WHERE `id` = ?", id).Exec()
	if err != nil {
		return 0, err
	}

	return r.RowsAffected()
}
