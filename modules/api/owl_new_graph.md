# OWL new graph 文件

### `GET` http://localhost/api/v1/owlgraph/keyword_search

| prams | description | example | 备注 |  
| ----- | ----------- | ------- | --- |
| q | 关键字查询 | 联通 | case sensitive |
| filter_type | 查询项目类别  | hostgroup | 支援项目 "platform", "idc", "isp", "province", "hostname", "hostgroup", 如不指定或是 "all" 表示全域关键字查询 |

### Dashboard `Screen` & `Graph`
文件请参考 doc/dashboard.json or [links](https://htmlpreview.github.io/?https://github.com/Cepave/open-falcon-backend/modules/api/blob/owl_new_graph/doc/dahsboard.html)


### 透過hostname列表取得counters list
| prams | description | example | 备注 |  
| ----- | ----------- | ------- | --- |
| q | 使用 regex 查询字符  | cpu | - |
| endpoints | hostname list  | agent01,agent02 | - |
請參考 [links](https://masato25.github.io/owl_backend/#/endpointstr_counter)
