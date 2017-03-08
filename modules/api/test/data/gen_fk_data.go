package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/fe/dmodel"
	"github.com/Pallinder/go-randomdata"
)

var (
	platform  = []string{"a-01", "a-02", "b-01", "c-04", "d-01", "c-11"}
	idc       = []string{"北京一区讯通", "上海一区联通", "武汉一区移动", "上海二区联通"}
	isp       = []string{"bgp", "cmb", "ctt", "hkk"}
	province  = []string{"北京", "上海", "武汉"}
	hostgroup = []string{"北京hg1", "外线hg2", "伺服器hg3", "testhg"}
)

func getPlatform() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(platform) - 1)
	return platform[indx]
}

func getIdc() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(idc) - 1)
	return idc[indx]
}

func getIsp() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(isp) - 1)
	return isp[indx]
}

func getProvince() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(province) - 1)
	return province[indx]
}
func getHostGroup() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(hostgroup) - 1)
	return hostgroup[indx]
}

func main() {
	f, err := os.Create("../fakeData.json")
	defer f.Close()
	nsize := 20
	res := make([]dmodel.BossObj, nsize)
	for i := 0; i < nsize; i++ {
		b := dmodel.BossObj{
			Platform: getPlatform(),
			Province: getProvince(),
			Isp:      getIsp(),
			Idc:      getIdc(),
			Ip:       randomdata.IpV4Address(),
			Hostname: randomdata.SillyName() + "_" + randomdata.StringNumberExt(4, "-", 2),
		}
		res[i] = b
	}
	slcB, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
	}
	f.Write(slcB)
}
