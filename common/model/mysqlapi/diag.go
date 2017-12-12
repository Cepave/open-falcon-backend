package mysqlapi

// HealthView is the model for responsed JSON data in ``/health`
// example:
// {
//    "rdb":{
//       "dsn":"DSN....",
//       "open_connections":10,
//       "ping_result":0,
//       "ping_message":"",
//       "<db1_name>": {
//       	"dsn":"DSN....",
//       	"open_connections":10,
//       	"ping_result":0,
//       	"ping_message":"",
//		 },
//       "<db2_name>": {
//       	"dsn":"DSN....",
//       	"open_connections":10,
//       	"ping_result":0,
//       	"ping_message":"",
//		 }
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
	Rdb  *AllRdbHealth `json:"rdb"`
	Http *Http         `json:"http"`
	Nqm  *Nqm          `json:"nqm"`
}

type AllRdbHealth struct {
	// Deprecated; old information on portal database
	Dsn string `json:"dsn"`
	// Deprecated; old information on portal database
	OpenConnections int `json:"open_connections"`
	// Deprecated; old information on portal database
	PingResult int `json:"ping_result"`
	// Deprecated; old information on portal database
	PingMessage string `json:"ping_message"`

	Portal *Rdb `json:"portal"`
	Graph  *Rdb `json:"graph"`
	Boss   *Rdb `json:"boss"`
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
