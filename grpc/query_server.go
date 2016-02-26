package grpc

import (
	"encoding/json"
	"fmt"
	cmodel "github.com/Cepave/common/model"
	"github.com/Cepave/fe/g"
	"github.com/Cepave/fe/graph"
	pb "github.com/Cepave/query/grpc/proto/owlquery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

//this function will generate query string obj for QueryRRDtool
func getq(q *pb.QueryInput) cmodel.GraphQueryParam {
	request := cmodel.GraphQueryParam{
		Start:     int64(q.StartTs),
		End:       int64(q.EndTs),
		ConsolFun: q.Consolfun,
		Endpoint:  q.Endpoint,
		Counter:   q.Counter,
	}
	return request
}
func (s *server) Query(ctx context.Context, in *pb.QueryInput) (*pb.QueryReply, error) {
	result, _ := graph.QueryOne(getq(in))
	//genreate json string and send back to client
	res, _ := json.Marshal(result)
	return &pb.QueryReply{Result: string(res)}, nil
}

func Start() {
	port := fmt.Sprintf(":%v", g.Config().Grpc.Port)
	log.Printf("start grpc in port %v ..", port)
	//queryrrd(1452806153, 1452827753, "AVERAGE", "docker-agent", "cpu.idle")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOwlQueryServer(s, &server{})
	s.Serve(lis)
}
