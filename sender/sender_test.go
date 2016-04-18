package sender

import (
	"encoding/json"
	cmodel "github.com/open-falcon/common/model"
	"testing"
)

func createMetaData() *cmodel.MetaData {
	in := cmodel.MetaData{
		Metric:      "test.metric.niean.1",
		Timestamp:   1460366463,
		Step:        60,
		Value:       0.000000,
		CounterType: "",
		Tags: map[string]string{
			"rttmin":             "18611",
			"rttavg":             "21626.5",
			"rttmax":             "26951",
			"rttmdev":            "234.2",
			"rttmedian":          "21.522",
			"pkttransmit":        "13",
			"pktreceive":         "12",
			"dstpoint":           "test.endpoint.niean.2",
			"agent-id":           "1334",
			"agent-isp-id":       "12",
			"agent-province-id":  "13",
			"agent-city-id":      "14",
			"agent-name-tag-id":  "123",
			"target-id":          "2334",
			"target-isp-id":      "22",
			"target-province-id": "23",
			"target-city-id":     "24",
			"target-name-tag-id": "223",
		},
	}

	return &in
}

func TestConvert2NqmMetrics(t *testing.T) {
	in := createMetaData()
	out_ptr, _ := convert2NqmMetrics(in)
	out := nqmMetrics{
		Rttmin:      18611,
		Rttavg:      21626.5,
		Rttmax:      26951,
		Rttmdev:     234.2,
		Rttmedian:   21.522,
		Pkttransmit: 13,
		Pktreceive:  12,
	}

	if out != *out_ptr {
		t.Error("Expected output: ", out)
		t.Error("Real output:     ", *out_ptr)
	}
}

func TestConvert2NqmEndpoint(t *testing.T) {
	in := createMetaData()
	out_ptr, _ := convert2NqmEndpoint(in, "agent")
	out := nqmEndpoint{
		Id:         1334,
		IspId:      12,
		ProvinceId: 13,
		CityId:     14,
		NameTagId:  123,
	}

	if out != *out_ptr {
		t.Error("Expected output: ", out)
		t.Error("Real output:     ", *out_ptr)

	}

	out_ptr, _ = convert2NqmEndpoint(in, "target")
	out = nqmEndpoint{
		Id:         2334,
		IspId:      22,
		ProvinceId: 23,
		CityId:     24,
		NameTagId:  223,
	}

	if out != *out_ptr {
		t.Error("Expected output: ", out)
		t.Error("Real output:     ", *out_ptr)
	}
}

func TestConvert2NqmRpcItem(t *testing.T) {
	in := createMetaData()
	out, _ := convert2NqmRpcItem(in)
	t.Log("convert2NqmRpcItem:", out)
}

func TestJsonMarshal(t *testing.T) {
	in := createMetaData()
	out, _ := convert2NqmEndpoint(in, "agent")
	check, _ := json.Marshal(out)
	t.Log("JsonMarshal of agent: ", string(check))
	var dat map[string]int
	json.Unmarshal(check, &dat)

	expected := map[string]int{
		"name_tag_id": 123,
		"id":          1334,
		"isp_id":      12,
		"province_id": 13,
		"city_id":     14,
	}

	for k, v := range expected {
		if v != dat[k] {
			t.Error("Expected output: ", expected)
			t.Error("Real output:     ", dat)
		}
	}

	out, _ = convert2NqmEndpoint(in, "target")
	check, _ = json.Marshal(out)
	t.Log("JsonMarshal of target: ", string(check))
	json.Unmarshal(check, &dat)

	expected = map[string]int{
		"name_tag_id": 223,
		"id":          2334,
		"isp_id":      22,
		"province_id": 23,
		"city_id":     24,
	}

	for k, v := range expected {
		if v != dat[k] {
			t.Error("Expected output: ", expected)
			t.Error("Real output:     ", dat)
		}
	}

	m_out, _ := convert2NqmMetrics(in)
	check, _ = json.Marshal(m_out)
	t.Log("JsonMarshal of metrics: ", string(check))
	var int_dat map[string]int32
	json.Unmarshal(check, &int_dat)

	var min int32 = 18611
	var max int32 = 26951

	if v, p := int_dat["min"]; p {
		if v != min {
			t.Error("Expected output: ", min)
			t.Error("Real output:     ", v)
		}
	}
	if v, p := int_dat["max"]; p {
		if v != max {
			t.Error("Expected output: ", max)
			t.Error("Real output:     ", v)
		}
	}
}
