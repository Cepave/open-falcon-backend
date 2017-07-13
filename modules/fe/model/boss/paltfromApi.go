package boss

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	"github.com/astaxie/beego/orm"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/sets/hashset"
	log "github.com/sirupsen/logrus"
)

func getFctoken() string {
	hasher := md5.New()
	io.WriteString(hasher, g.Config().Api.Token)
	s := hex.EncodeToString(hasher.Sum(nil))

	t := time.Now()
	now := t.Format("20060102")
	s = now + s

	hasher = md5.New()
	io.WriteString(hasher, s)
	fctoken := hex.EncodeToString(hasher.Sum(nil))
	return fctoken
}

func GetPlatformASJSON() (repons PlatformList, err error) {
	config := g.Config()
	fctoken := getFctoken()
	url := config.Api.Map + "/fcname/" + config.Api.Name + "/fctoken/" + fctoken
	url += "/show_active/yes/hostname/yes/pop_id/yes/ip/yes.json"
	log.Debugf("platform get url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &repons)
	if repons.Status != 1 {
		err = errors.New("Platform list quer failed")
	}
	return
}

//generate platform mapping of triggered alarm cases
func GenPlatMap(repons PlatformList, filterList *hashset.Set) (ipMapper *hashmap.Map, platList *hashset.Set, popIds []int) {
	log.Debugf("filterLIST: %v", filterList.Values())
	res := repons.Result
	ipMapper = hashmap.New()
	platList = hashset.New()
	popIDSet := hashset.New()
	for _, platform := range res {
		for _, ipInfo := range platform.IPList {
			//only get platform info that matched triggered alarm case
			if filterList.Contains(ipInfo.HostName) {
				platList.Add(platform.Platform)
				ipInfo.Platform = platform.Platform
				popID, _ := strconv.Atoi(ipInfo.POPID)
				popIDSet.Add(popID)
				ip := ipInfo.IP
				if len(ip) > 0 && ip == getIPFromHostname(ipInfo.HostName) {
					ipMapper.Put(ipInfo.HostName, ipInfo)
				}
			}
		}
	}
	log.Debugf("popID set: %v", popIDSet.Values())
	popIds = []int{}
	for _, popID := range popIDSet.Values() {
		popIds = append(popIds, popID.(int))
	}
	return
}

func getIPFromHostname(hostname string) string {
	ip := ""
	fragments := strings.Split(hostname, "-")
	slice := []string{}
	if len(fragments) == 6 {
		fragments := fragments[2:]
		for _, fragment := range fragments {
			num, err := strconv.Atoi(fragment)
			if err == nil {
				slice = append(slice, strconv.Itoa(num))
			}
		}
		if len(slice) == 4 {
			ip = strings.Join(slice, ".")
		}
	}
	return ip
}

func IdcMapping(popIDs []int) (idcmapping map[int]string, err error) {
	if len(popIDs) == 0 {
		err = errors.New("popIdDs is empty")
		return
	}
	q := orm.NewOrm()
	q.Using("grafana")
	popFilter := ""
	for idx, pop := range popIDs {
		if idx == 0 {
			popFilter = fmt.Sprintf("%d", pop)
		}
		popFilter = fmt.Sprintf("%s, %d", popFilter, pop)
	}
	sqlcommd := fmt.Sprintf("select name,pop_id from grafana.idc where pop_id IN (%s)", popFilter)
	idcs := []IDC{}
	_, err = q.Raw(sqlcommd).QueryRows(&idcs)
	idcmapping = map[int]string{}
	for _, idc := range idcs {
		idcmapping[idc.PopId] = idc.Name
	}
	return
}
