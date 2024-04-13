package grpc

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/shashank-priyadarshi/bench/common"
	"io"
	"time"
)

const (
	grpc_key = "grpc_"
)

type Server struct {
	common.BenchmarkingServer
	RDB    *redis.Client
	Server string
}

func (s *Server) Bidirectional(server common.Benchmarking_BidirectionalServer) error {

	for {
		req, err := server.Recv()
		if req == nil {
			fmt.Println("nil request received: ", req, " ,", err, " ,server: ", s.Server)
			//return err
			continue
		}
		receivedTime := time.Now().UnixMilli()

		if err == io.EOF {
			fmt.Println("stream ended")
			return nil
		}

		if err != nil {
			fmt.Println("error while receiving stream: ", err)
		}

		if err = s.RDB.Set(context.Background(), fmt.Sprintf("%s%s", grpc_key, req.Name), receivedTime-req.Time, 0).Err(); err != nil {
			fmt.Println("error while putting data to redis: ", err)
		}

	}
	//return nil
}
