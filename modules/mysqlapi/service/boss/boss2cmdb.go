package boss

import (
	model "github.com/Cepave/open-falcon-backend/modules/mysqlapi/model"
	underscore "github.com/ahl5esoft/golang-underscore"
	"strings"
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
	for h := range tuples {
		append(rs, h.Hostname)
	}
	return rs
}

func boss2cmdb(hosts []*model.BossHost) *model.SyncForAdding {
	sync := model.SyncForAdding{}
	tuples := []tHost{}
	owl_default_grp_hosts := []string{}
	for h := range hosts {
		// transform hosts
		// Hostname and IP not Null property guaranted by schema
		// Activate and Platform might be Null value
		if h.Activate.Valid {
			append(sync.Hosts, &model.SyncHost{
				Activate: int(h.Activate.Int64),
				Name:     h.Hostname,
				IP:       h.Ip,
			})
		} else {
			// treat activate NULL value as not activated
			append(sync.Hosts, &model.SyncHost{
				Activate: 0,
				Name:     h.Hostname,
				IP:       h.Ip,
			})
		}
		// initialize grp-host tuples
		if h.Platform.Valid {
			append(tuples, tHost{
				Grpname:  h.Platform.String,
				Hostname: h.Hostname,
			})
		} else {
			// skip Platform NULL value
		}
		// initialize hostname for default group
		append(owl_default_grp_hosts, h.Hostname)
	}

	// begin grps and relations transform
	// group-by host tuples
	t := underscore.GroupBy(tuples, "grpname")
	if grps_dict, ok := t.(map[string][]tHost); ok {
		for grp_name, tuple_list := range grps_dict {
			// transform grps
			append(sync.Hostgroups, &model.SyncHostGroup{
				Creator: "root",
				Name:    grp_name,
			})
			// transform relations
			sync.Relations[grp_name] = toHostnameSlice(tuple_list)
		}
		append(sync.Hostgroups, &model.SyncHostGroup{
			Creator: "root",
			Name:    DEFAULT_GRP,
		})
	}
	sync.Relations[DEFAULT_GRP] = owl_default_grp_hosts
	return &sync
}
