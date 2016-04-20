package sender

import (
	"fmt"
)

type nqmEndpoint struct {
	Id         int32 `json:"id"`
	IspId      int16 `json:"isp_id"`
	ProvinceId int16 `json:"province_id"`
	CityId     int16 `json:"city_id"`
	NameTagId  int32 `json:"name_tag_id"`
}

func (end nqmEndpoint) String() string {
	return fmt.Sprintf(
		"Id:[%d] IspId:(%d) ProvinceId:(%d), CityId:[%d], NameTagId:[%d]",
		end.Id,
		end.IspId,
		end.ProvinceId,
		end.CityId,
		end.NameTagId,
	)
}

type nqmMetrics struct {
	Rttmin      int32   `json:"min"`
	Rttavg      float32 `json:"avg"`
	Rttmax      int32   `json:"max"`
	Rttmdev     float32 `json:"mdev"`
	Rttmedian   float32 `json:"med"`
	Pkttransmit int32   `json:"sentPackets"`
	Pktreceive  int32   `json:"receivedPackets"`
}

func (metric nqmMetrics) String() string {
	return fmt.Sprintf(
		"Rttmin:%v, Rttavg:%v, Rttmax:%v, Rttmdev:%v, Rttmedian:%v, Pkttransmit:%v, Pktreceive:%v",
		metric.Rttmin,
		metric.Rttavg,
		metric.Rttmax,
		metric.Rttmdev,
		metric.Rttmedian,
		metric.Pkttransmit,
		metric.Pktreceive,
	)
}

type nqmRpcItem struct {
	Timestamp int64       `json:"time"`
	Agent     nqmEndpoint `json:"agent"`
	Target    nqmEndpoint `json:"target"`
	Metrics   nqmMetrics  `json:"metrics"`
}

func (this nqmRpcItem) String() string {
	return fmt.Sprintf(
		"<TS:%d, Src:<%v>, Dst:<%v>, Metrics:<%v>>",
		this.Timestamp,
		this.Agent,
		this.Target,
		this.Metrics,
	)
}
