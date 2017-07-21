package graph

import (
	"fmt"

	"errors"
	cmodel "github.com/Cepave/open-falcon-backend/common/model"
	cutils "github.com/Cepave/open-falcon-backend/common/utils"
	log "github.com/sirupsen/logrus"
	spool "github.com/toolkits/pool/simple_conn_pool"
	"math"
	"time"
)

func QueryBatch(para []cmodel.GraphQueryParam) (results []cmodel.GraphQueryResponse, errs []error) {
	addressMaper, errorCollector := shardRespPools2(para)

	for k, v := range addressMaper {
		log.Infof("QueryBatch debugging: %v has %d size", k, len(v))
		rpcResource, err := fetchRpcConn(k)
		if err != nil {
			errorCollector = append(errorCollector, fmt.Errorf("[fetchRpcConn] lost query num: %v response, because got error %v", len(v), err.Error()))
			continue
		}
		resp, err := lastBatchReq2(rpcResource, v)
		if err != nil {
			errorCollector = append(errorCollector, fmt.Errorf("[lastBatchReq] lost query num: %v response, because got error %v", len(v), err.Error()))
			continue
		}
		if resp != nil {
			for _, rt := range *resp {
				if rt != nil {
					results = append(results, fixedResp(*rt))
				}
			}
		}
	}
	return results, errorCollector
}

func fixedResp(resp cmodel.GraphQueryResponse) cmodel.GraphQueryResponse {
	if len(resp.Values) < 1 {
		resp.Values = []*cmodel.RRDData{}
		return resp
	}

	// TODO query不该做这些事情, 说明graph没做好
	fixed := []*cmodel.RRDData{}
	for _, v := range resp.Values {
		if v == nil {
			continue
		}
		//FIXME: 查询数据的时候，把所有的负值都过滤掉，因为transfer之前在设置最小值的时候为U
		if (resp.DsType == "DERIVE" || resp.DsType == "COUNTER") && v.Value < 0 {
			fixed = append(fixed, &cmodel.RRDData{Timestamp: v.Timestamp, Value: cmodel.JsonFloat(math.NaN())})
		} else {
			fixed = append(fixed, v)
		}
	}
	resp.Values = fixed
	return resp
}

func shardRespPools2(para []cmodel.GraphQueryParam) (map[string][]cmodel.GraphQueryParam, []error) {
	var errorCollector []error
	addressMaper := map[string][]cmodel.GraphQueryParam{}
	for _, p := range para {
		addr, err := selectAddress(p.Endpoint, p.Counter)
		mkey := fmt.Sprintf("%s, %s", p.Endpoint, p.Counter)
		if err != nil {
			errorCollector = append(errorCollector, fmt.Errorf("\"%s\" got error: %v", mkey, err.Error()))
			continue
		}
		if val, ok := addressMaper[addr]; !ok {
			addressMaper[addr] = []cmodel.GraphQueryParam{p}
		} else {
			addressMaper[addr] = append(val, p)
		}
	}
	return addressMaper, errorCollector
}

func lastBatchReq2(rpcResource rpcConnectionResurce, para []cmodel.GraphQueryParam) (*[]*cmodel.GraphQueryResponse, error) {
	type ChResult struct {
		Err  error
		Resp *cmodel.GraphQueryResponseList
	}
	callMethod := "Graph.BatchQuery"
	ch := make(chan *ChResult, 1)
	go func() {
		resp := &cmodel.GraphQueryResponseList{}
		err := rpcResource.RpcConn.Call(callMethod, para, resp)
		ch <- &ChResult{Err: err, Resp: resp}
	}()

	select {
	case <-time.After(time.Duration(callTimeout) * time.Millisecond):
		rpcResource.ForceClosePool()
		return nil, fmt.Errorf("%s, call timeout. proc: %s", rpcResource.Address, callMethod)
	case r := <-ch:
		if r.Err != nil {
			rpcResource.ForceClosePool()
			return nil, fmt.Errorf("%s, call failed, err %v. proc: %s", rpcResource.Address, r.Err, callMethod)
		} else {
			rpcResource.ForceClosePool()
			return r.Resp.List, nil
		}
	}
}

func LastBatch(para []cmodel.GraphLastParam) (results []cmodel.GraphLastResp, errs []error) {
	log.Info("LastBatch")
	addressMaper, errorCollector := shardRespPools(para)
	for k, v := range addressMaper {
		rpcResource, err := fetchRpcConn(k)
		if err != nil {
			errorCollector = append(errorCollector, fmt.Errorf("[fetchRpcConn] lost query num: %v response, because got error %v", len(v), err.Error()))
			continue
		}
		resp, err := lastBatchReq(rpcResource, v)
		if err != nil {
			errorCollector = append(errorCollector, fmt.Errorf("[lastBatchReq] lost query num: %v response, because got error %v", len(v), err.Error()))
			continue
		}
		if resp != nil {
			results = append(results, *resp...)
		}
	}
	return results, errorCollector
}

func shardRespPools(para []cmodel.GraphLastParam) (map[string][]cmodel.GraphLastParam, []error) {
	var errorCollector []error
	addressMaper := map[string][]cmodel.GraphLastParam{}
	for _, p := range para {
		addr, err := selectAddress(p.Endpoint, p.Counter)
		mkey := fmt.Sprintf("%s, %s", p.Endpoint, p.Counter)
		if err != nil {
			errorCollector = append(errorCollector, fmt.Errorf("\"%s\" got error: %v", mkey, err.Error()))
			continue
		}
		if val, ok := addressMaper[addr]; !ok {
			addressMaper[addr] = []cmodel.GraphLastParam{p}
		} else {
			addressMaper[addr] = append(val, p)
		}
	}
	return addressMaper, errorCollector
}

func selectAddress(endpoint, counter string) (string, error) {
	pkey := cutils.PK2(endpoint, counter)
	node, err := GraphNodeRing.GetNode(pkey)
	if err != nil {
		return "", err
	}
	addr, found := clusterMap[node]
	if !found {
		return "", errors.New("node not found")
	}
	return addr, nil
}

type rpcConnectionResurce struct {
	Address string
	Pool    *spool.ConnPool
	Conn    spool.NConn
	RpcConn *spool.RpcClient
}

func (mine rpcConnectionResurce) ForceClosePool() {
	mine.Pool.ForceClose(mine.Conn)
}

func fetchRpcConn(address string) (rpcConnectionResurce, error) {
	rcpRes := rpcConnectionResurce{Address: address}
	pool, found := GraphConnPools.Get(address)
	if !found {
		log.Errorf("pool :%v not found, address: %v", pool, address)
		return rcpRes, fmt.Errorf("select pool: %v got error with address: %v", pool, address)
	}
	rcpRes.Pool = pool
	conn, err := rcpRes.Pool.Fetch()
	if err != nil {
		log.Errorf("fetch %v connection got error: %v", address, err)
		return rcpRes, fmt.Errorf("fetch %v connection got error: %v", address, err)
	}
	rcpRes.Conn = conn
	rpcConn := rcpRes.Conn.(spool.RpcClient)
	if rpcConn.Closed() {
		pool.ForceClose(conn)
		return rcpRes, errors.New("conn closed")
	}
	rcpRes.RpcConn = &rpcConn
	return rcpRes, nil
}

func lastBatchReq(rpcResource rpcConnectionResurce, para []cmodel.GraphLastParam) (*[]cmodel.GraphLastResp, error) {
	type ChResult struct {
		Err  error
		Resp *cmodel.GraphLastRespList
	}
	callMethod := "Graph.LastBatch"
	ch := make(chan *ChResult, 1)
	go func() {
		resp := &cmodel.GraphLastRespList{}
		err := rpcResource.RpcConn.Call(callMethod, para, resp)
		ch <- &ChResult{Err: err, Resp: resp}
	}()

	select {
	case <-time.After(time.Duration(callTimeout) * time.Millisecond):
		rpcResource.ForceClosePool()
		return nil, fmt.Errorf("%s, call timeout. proc: %s", rpcResource.Address, callMethod)
	case r := <-ch:
		if r.Err != nil {
			rpcResource.ForceClosePool()
			return nil, fmt.Errorf("%s, call failed, err %v. proc: %s", rpcResource.Address, r.Err, callMethod)
		} else {
			rpcResource.ForceClosePool()
			return r.Resp.List, nil
		}
	}
}
