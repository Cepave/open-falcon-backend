package db

import (
	"github.com/Cepave/common/model"
	"log"
)

func QueryGroupTemplates() (map[int][]int, error) {
	m := make(map[int][]int)

	sql := "select grp_id, tpl_id from grp_tpl"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var gid, tid int
		err = rows.Scan(&gid, &tid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if _, exists := m[gid]; exists {
			m[gid] = append(m[gid], tid)
		} else {
			m[gid] = []int{tid}
		}
	}

	return m, nil
}

// 获取所有的策略模板列表
func QueryTemplates() (map[int]*model.Template, error) {

	templates := make(map[int]*model.Template)

	sql := "select id, tpl_name, parent_id, action_id, create_user from tpl"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return templates, err
	}

	defer rows.Close()
	for rows.Next() {
		t := model.Template{}
		err = rows.Scan(&t.Id, &t.Name, &t.ParentId, &t.ActionId, &t.Creator)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		templates[t.Id] = &t
	}

	return templates, nil
}

// 一个机器ID对应了多个模板ID
func QueryHostTemplateIds() (map[int][]int, error) {
	ret := make(map[int][]int)
	rows, err := DB.Query("select a.tpl_id, b.host_id from grp_tpl as a inner join grp_host as b on a.grp_id=b.grp_id")
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		var tid, hid int

		err = rows.Scan(&tid, &hid)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		if _, ok := ret[hid]; ok {
			ret[hid] = append(ret[hid], tid)
		} else {
			ret[hid] = []int{tid}
		}
	}

	return ret, nil
}
