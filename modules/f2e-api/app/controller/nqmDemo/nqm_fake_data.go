package nqmDemo

var agentData = []byte(`
[
  {
    "city": {
      "id": 8,
      "name": "上海"
    },
    "comment": null,
    "connection_id": "aaa-mt-127-0-0-99@127.0.0.99",
    "group_tags": [],
    "hostname": "aaa-mt-127-0-0-99",
    "id": 3876,
    "ip_address": "1127.0.0.99",
    "isp": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "last_heartbeat_time": 1493966667,
    "name": null,
    "name_tag": {
      "id": -1,
      "value": "<UNDEFINED>"
    },
    "num_of_enabled_pingtasks": 0,
    "province": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "status": true
  },
  {
    "city": {
      "id": 1,
      "name": "北京市"
    },
    "comment": null,
    "connection_id": "aaa-mt-10-0-0-1@10.0.0.1",
    "group_tags": [],
    "hostname": "aaa-mt-10-0-0-1",
    "id": 4146,
    "ip_address": "0.0.0.0",
    "isp": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "last_heartbeat_time": null,
    "name": "aaa-mt-10-0-0-1",
    "name_tag": {
      "id": -1,
      "value": "<UNDEFINED>"
    },
    "num_of_enabled_pingtasks": 0,
    "province": {
      "id": 4,
      "name": "北京"
    },
    "status": true
  },
  {
    "city": {
      "id": 1,
      "name": "北京市"
    },
    "comment": null,
    "connection_id": "aaa-mt-10-0-0-3@10.0.0.2",
    "group_tags": [],
    "hostname": "aaa-mt-10-0-0-2",
    "id": 4155,
    "ip_address": "10.0.0.2",
    "isp": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "last_heartbeat_time": null,
    "name": "aaa-mt-10-0-0-2",
    "name_tag": {
      "id": -1,
      "value": "<UNDEFINED>"
    },
    "num_of_enabled_pingtasks": 0,
    "province": {
      "id": 4,
      "name": "北京"
    },
    "status": true
  },
  {
    "city": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "comment": null,
    "connection_id": "aaa-mt-10-0-0-3@10.0.0.3",
    "group_tags": [],
    "hostname": "aaa-mt-10.0.0.3",
    "id": 4500,
    "ip_address": "10.0.0.3",
    "isp": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "last_heartbeat_time": 1502169389,
    "name": null,
    "name_tag": {
      "id": -1,
      "value": "<UNDEFINED>"
    },
    "num_of_enabled_pingtasks": 0,
    "province": {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    "status": true
  }
]
`)

var IspData = []byte(`
[
  {
    "id": -1,
    "name": "<UNDEFINED>",
    "acronym": "<UNDEFINED>"
  },
  {
    "id": 1,
    "name": "北京三信时代",
    "acronym": "北京三信时代"
  },
  {
    "id": 1,
    "name": "教育网",
    "acronym": "教育网"
  },
  {
    "id": 3,
    "name": "移动",
    "acronym": "移动"
  },
  {
    "id": 4,
    "name": "电信",
    "acronym": "电信"
  },
  {
    "id": 5,
    "name": "联通",
    "acronym": "联通"
  }
]
`)

var provincesData = []byte(`
  [
    {
      "id": -1,
      "name": "<UNDEFINED>"
    },
    {
      "id": 8,
      "name": "上海"
    },
    {
      "id": 35,
      "name": "中国其它"
    },
    {
      "id": 25,
      "name": "云南"
    },
    {
      "id": 1,
      "name": "内蒙古"
    },
    {
      "id": 4,
      "name": "北京"
    },
    {
      "id": 6,
      "name": "吉林"
    },
    {
      "id": 24,
      "name": "四川"
    },
    {
      "id": 36,
      "name": "国外"
    },
    {
      "id": 12,
      "name": "天津"
    },
    {
      "id": 28,
      "name": "宁夏"
    },
    {
      "id": 10,
      "name": "安徽"
    },
    {
      "id": 11,
      "name": "山东"
    },
    {
      "id": 2,
      "name": "山西"
    },
    {
      "id": 20,
      "name": "广东"
    },
    {
      "id": 21,
      "name": "广西"
    },
    {
      "id": 30,
      "name": "新疆"
    },
    {
      "id": 9,
      "name": "江苏"
    },
    {
      "id": 14,
      "name": "江西"
    },
    {
      "id": 23,
      "name": "海南"
    },
    {
      "id": 18,
      "name": "湖北"
    },
    {
      "id": 17,
      "name": "湖南"
    },
    {
      "id": 33,
      "name": "澳门"
    },
    {
      "id": 27,
      "name": "甘肃"
    },
    {
      "id": 16,
      "name": "重庆"
    },
    {
      "id": 26,
      "name": "陕西"
    },
    {
      "id": 29,
      "name": "青海"
    },
    {
      "id": 32,
      "name": "香港"
    }
  ]
`)

var targetsData = []byte(`
  [
    {
      "available": true,
      "city": {
        "id": 137,
        "name": "宁波市"
      },
      "comment": null,
      "creation_time": 1490839434,
      "group_tags": [],
      "host": "180.0.0.1",
      "id": 44461,
      "isp": {
        "id": 6,
        "name": "联通"
      },
      "name": "180.0.0.1-浙江多线",
      "name_tag": {
        "id": 59,
        "value": "联通-浙江宁波多线"
      },
      "probed_by_all": false,
      "province": {
        "id": 13,
        "name": "浙江"
      },
      "status": true
    },
    {
      "available": true,
      "city": {
        "id": 168,
        "name": "洛阳市"
      },
      "comment": null,
      "creation_time": 1488339890,
      "group_tags": [
        {
          "id": 1,
          "name": "内网"
        }
      ],
      "host": "10.0.0.1",
      "id": 44459,
      "isp": {
        "id": 3,
        "name": "移动"
      },
      "name": "10.0.0.1-河南洛阳移动",
      "name_tag": {
        "id": -1,
        "value": "<UNDEFINED>"
      },
      "probed_by_all": false,
      "province": {
        "id": 19,
        "name": "河南"
      },
      "status": true
    },
    {
      "available": true,
      "city": {
        "id": 127,
        "name": "镇江市"
      },
      "comment": null,
      "creation_time": 1485156076,
      "group_tags": [
        {
          "id": 1,
          "name": "内网"
        }
      ],
      "host": "10.0.0.2",
      "id": 44444,
      "isp": {
        "id": 6,
        "name": "联通"
      },
      "name": "10.0.0.2-江苏镇江",
      "name_tag": {
        "id": 57,
        "value": "联通-江苏镇江"
      },
      "probed_by_all": false,
      "province": {
        "id": 9,
        "name": "江苏"
      },
      "status": true
    },
    {
      "available": true,
      "city": {
        "id": 238,
        "name": "青岛市"
      },
      "comment": null,
      "creation_time": 1485156076,
      "group_tags": [
        {
          "id": 1,
          "name": "内网"
        }
      ],
      "host": "180.0.0.2",
      "id": 44445,
      "isp": {
        "id": 3,
        "name": "移动"
      },
      "name": "180.0.0.2-山东青岛",
      "name_tag": {
        "id": 37,
        "value": "移动-山东青岛"
      },
      "probed_by_all": false,
      "province": {
        "id": 11,
        "name": "山东"
      },
      "status": true
    }
  ]`)

var pingTaskData = []byte(`
  [
    {
      "id": 2,
      "period": 1,
      "name": null,
      "enable": true,
      "comment": null,
      "num_of_enabled_agents": 16,
      "filter": {
        "isps": [
          {
            "id": 6,
            "name": "联通"
          }
        ],
        "provinces": [],
        "cities": [],
        "name_tags": [],
        "group_tags": []
      }
    },
    {
      "id": 15,
      "period": 1,
      "name": null,
      "enable": true,
      "comment": null,
      "num_of_enabled_agents": 11,
      "filter": {
        "isps": [
          {
            "id": 3,
            "name": "移动"
          }
        ],
        "provinces": [],
        "cities": [],
        "name_tags": [],
        "group_tags": []
      }
    },
    {
      "id": 1,
      "period": 1,
      "name": null,
      "enable": true,
      "comment": null,
      "num_of_enabled_agents": 8,
      "filter": {
        "isps": [
          {
            "id": 5,
            "name": "电信"
          }
        ],
        "provinces": [],
        "cities": [],
        "name_tags": [],
        "group_tags": []
      }
    },
    {
      "id": 73,
      "period": 1,
      "name": "内网",
      "enable": true,
      "comment": null,
      "num_of_enabled_agents": 95,
      "filter": {
        "isps": [],
        "provinces": [],
        "cities": [],
        "name_tags": [],
        "group_tags": [
          {
            "id": 1,
            "name": "内网"
          }
        ]
      }
    }
  ]
`)

var CityData = []byte(`
  [
    {
      "id": -1,
      "name": "<UNDEFINED>",
      "post_code": "<UNDEFINED>",
      "province": {
        "id": -1,
        "name": "<UNDEFINED>"
      }
    },
    {
      "id": 1,
      "name": "北京市",
      "post_code": "100000",
      "province": {
        "id": 4,
        "name": "北京"
      }
    },
    {
      "id": 2,
      "name": "天津市",
      "post_code": "300000",
      "province": {
        "id": 12,
        "name": "天津"
      }
    },
    {
      "id": 3,
      "name": "上海市",
      "post_code": "200000",
      "province": {
        "id": 8,
        "name": "上海"
      }
    },
    {
      "id": 4,
      "name": "重庆市",
      "post_code": "400000",
      "province": {
        "id": 16,
        "name": "重庆"
      }
    },
    {
      "id": 5,
      "name": "广州市",
      "post_code": "510000",
      "province": {
        "id": 20,
        "name": "广东"
      }
    },
    {
      "id": 6,
      "name": "深圳市",
      "post_code": "518000",
      "province": {
        "id": 20,
        "name": "广东"
      }
    },
    {
      "id": 7,
      "name": "东莞市",
      "post_code": "511700",
      "province": {
        "id": 20,
        "name": "广东"
      }
    },
    {
      "id": 294,
      "name": "国外",
      "post_code": "",
      "province": {
        "id": -1,
        "name": "<UNDEFINED>"
      }
    }
  ]
`)

var nameTagsData = []byte(`
  [
    {
      "id": 99,
      "value": "内网"
    },
    {
      "id": 1,
      "value": "北京移动（铁通）"
    },
    {
      "id": 2,
      "value": "教育网北京"
    },
    {
      "id": 3,
      "value": "教育网广东"
    },
    {
      "id": 4,
      "value": "教育网湖北"
    },
    {
      "id": 5,
      "value": "电信-山东"
    },
    {
      "id": 6,
      "value": "电信-江苏"
    },
    {
      "id": 7,
      "value": "电信-江西"
    },
    {
      "id": 8,
      "value": "电信-河南"
    },
    {
      "id": 9,
      "value": "电信云南"
    },
    {
      "id": 10,
      "value": "电信四川"
    }
  ]
`)

var groupTagsData = []byte(`
  [
    {
      "id": 1,
      "name": "内网"
    },
    {
      "id": 2,
      "name": "外网"
    }
  ]
`)
