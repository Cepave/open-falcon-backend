package mysqlapi

// HealthView is the model for responsed JSON data in ``/health`
// example:
// {
//    "rdb":{
//       "dsn":"DSN....",
//       "open_connections":10,
//       "ping_result":0,
//       "ping_message":""
//    },
//    "http":{
//       "listening":":6040"
//    },
//    "nqm":{
//       "heartbeat":{
//          "count":21387
//       }
//    }
// }
type HealthView struct {
	Rdb  *Rdb  `json:"rdb"`
	Http *Http `json:"http"`
	Nqm  *Nqm  `json:"nqm"`
}

type Rdb struct {
	Dsn             string `json:"dsn"`
	OpenConnections int    `json:"open_connections"`
	PingResult      int    `json:"ping_result"`
	PingMessage     string `json:"ping_message"`
}

type Http struct {
	Listening string `json:"listening"`
}

type Nqm struct {
	Heartbeat *Heartbeat `json:"heartbeat"`
}

type Heartbeat struct {
	Count uint64 `json:"count"`
}

// :~)
