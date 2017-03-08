package filter

import (
	"strings"

	"github.com/Cepave/open-falcon-backend/modules/api/app/model/boss"
)

func PlatformFilter(dat []boss.BossHost, filterTxt string) []boss.BossHost {
	res := []boss.BossHost{}
	for _, n := range dat {
		if strings.Contains(n.Platform, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func IspFilter(dat []boss.BossHost, filterTxt string) []boss.BossHost {
	res := []boss.BossHost{}
	for _, n := range dat {
		if strings.Contains(n.Isp, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func IdcFilter(dat []boss.BossHost, filterTxt string) []boss.BossHost {
	res := []boss.BossHost{}
	for _, n := range dat {
		if strings.Contains(n.Idc, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func IpFilter(dat []boss.BossHost, filterTxt string) []boss.BossHost {
	res := []boss.BossHost{}
	for _, n := range dat {
		if strings.Contains(n.Ip, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func ProvinceFilter(dat []boss.BossHost, filterTxt string) []boss.BossHost {
	res := []boss.BossHost{}
	for _, n := range dat {
		if strings.Contains(n.Province, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}

func HostNameFilter(dat []boss.BossHost, filterTxt string) []boss.BossHost {
	res := []boss.BossHost{}
	for _, n := range dat {
		if strings.Contains(n.Hostname, filterTxt) {
			res = append(res, n)
		}
	}
	return res
}
