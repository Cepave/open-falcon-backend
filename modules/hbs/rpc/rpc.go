package rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Cepave/open-falcon-backend/modules/hbs/g"
)

type Hbs int
type Agent int
type NqmAgent int

func Start() {
	nqmAgentHbsService.Start()
	agentHeartbeatService.Start()

	addr := g.Config().Listen

	server := rpc.NewServer()
	// server.Register(new(filter.Filter))
	server.Register(new(Agent))
	server.Register(new(Hbs))
	server.Register(new(NqmAgent))

	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln("listen error:", e)
	} else {
		log.Println("listening", addr)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("listener accept fail:", err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
func Stop() {
	nqmAgentHbsService.Stop()
}
