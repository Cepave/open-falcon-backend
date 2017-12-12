package boss

import (
	model "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	underscore "github.com/ahl5esoft/golang-underscore"
)

const (
	DEFAULT_GRP = "Owl_Default_Group"
)

// nonNull grp-host tuple
type tHost struct {
	Grpname  string
	Hostname string
}

func toHostnameSlice(tuples []tHost) []string {
	rs := []string{}
	for _, h := range tuples {
		rs = append(rs, h.Hostname)
	}
	return rs
}

func Boss2cmdb(hosts []*model.BossHost) *model.SyncForAdding {
	sync := model.SyncForAdding{}
	sync.Relations = make(map[string][]string)
	tuples := []tHost{}
	owl_default_grp_hosts := []string{}
	for _, h := range hosts {
		// transform hosts
		// Hostname and IP not Null property guaranteed by schema
		// Activate and Platform might be Null value
		if h.Activate.Valid {
			sync.Hosts = append(sync.Hosts, &model.SyncHost{
				Activate: int(h.Activate.Int64),
				Name:     h.Hostname,
				IP:       h.Ip,
			})
		} else {
			// treat activate NULL value as not activated
			sync.Hosts = append(sync.Hosts, &model.SyncHost{
				Activate: 0,
				Name:     h.Hostname,
				IP:       h.Ip,
			})
		}
		// initialize grp-host tuples
		if h.Platform.Valid {
			tuples = append(tuples, tHost{
				Grpname:  h.Platform.String,
				Hostname: h.Hostname,
			})
		} else {
			// skip Platform NULL value
		}
		// initialize hostname for default group
		owl_default_grp_hosts = append(owl_default_grp_hosts, h.Hostname)
	}

	// begin grps and relations transform
	// group-by host tuples
	t := underscore.GroupBy(tuples, "grpname")
	if grps_dict, ok := t.(map[string][]tHost); ok {
		for grp_name, tuple_list := range grps_dict {
			// transform grps
			sync.Hostgroups = append(sync.Hostgroups, &model.SyncHostGroup{
				Creator: "root",
				Name:    grp_name,
			})
			// transform relations
			sync.Relations[grp_name] = toHostnameSlice(tuple_list)
		}
		sync.Hostgroups = append(sync.Hostgroups, &model.SyncHostGroup{
			Creator: "root",
			Name:    DEFAULT_GRP,
		})
	}
	sync.Relations[DEFAULT_GRP] = owl_default_grp_hosts
	return &sync
}
