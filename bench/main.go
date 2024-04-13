package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/shashank-priyadarshi/bench/common"
	grpc2 "github.com/shashank-priyadarshi/bench/grpc"
	vegeta "github.com/tsenart/vegeta/lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	port_          = "809"
	redis_         = "localhost:6379"
	redis_password = ""
	nats_          = nats.DefaultURL
	subject        = "test"
	nats_key       = "nats_"
)

var (
	wg             sync.WaitGroup
	grpc_clients   = make(map[string]common.Benchmarking_BidirectionalClient)
	nats_publisher = make(map[string]*nats.Conn)
	rdbClient      = redis.NewClient(&redis.Options{
		Addr:     redis_,
		Password: redis_password,
		DB:       0,
	})
	stats []statistic
)

type statistic struct {
	name      string
	timestamp int64
}

func init() {
	startGRPC(&wg, rdbClient)
	startNATS(&wg, rdbClient)
	go startServer(&wg)
}

func startGRPC(wg *sync.WaitGroup, rdbClient *redis.Client) {

	// one server and one client
	port := fmt.Sprintf("localhost:%s%d", port_, 1)

	go func(port string) {
		wg.Add(1)
		defer wg.Done()
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to start server %v", err)
		}

		grpcServer := grpc.NewServer()
		common.RegisterBenchmarkingServer(grpcServer, &grpc2.Server{RDB: rdbClient, Server: port})
		log.Printf("server started at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to start: %v", err)
		}
	}(port)

	time.Sleep(2 * time.Second)

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := common.NewBenchmarkingClient(conn)
	stream, err := client.Bidirectional(context.Background())
	if err != nil {
		fmt.Println("error while fetching stream for grpc server: ", port, " , ", err)
	} else {
		grpc_clients[fmt.Sprintf("grpc_%s", port)] = stream
	}

}

func startNATS(wg *sync.WaitGroup, rdbClient *redis.Client) {
	// one subscriber, three publishers
	nc, _ := nats.Connect(nats_)

	nats_publisher[fmt.Sprintf("pub_%d", 1)] = nc

	nc.Subscribe(subject, func(msg *nats.Msg) {
		// write received message time to redis
		receivedTime := time.Now().UnixMilli()

		var payload common.RequestMessage
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			log.Printf("Error processing data: %v", err)
			return
		}

		go rdbClient.Set(context.Background(), fmt.Sprintf("%s_%s", nats_key, payload.Name), receivedTime-payload.Time, 0)
		stats = append(stats, statistic{name: fmt.Sprintf("%s_%s", nats_key, payload.Name), timestamp: receivedTime - payload.Time})
	})
}

func startServer(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	r := gin.Default()
	r.POST("/grpc", func(c *gin.Context) {
		payload := common.RequestMessage{Payload: make([]byte, 100_000)}
		payload.Name = fmt.Sprintf("message_%v", c.Query("count"))
		for _, val := range grpc_clients {
			payload.Time = time.Now().UnixMilli()
			_, _ = json.Marshal(&payload)
			go func(payload *common.RequestMessage) {
				err := val.Send(payload)
				if err != nil {
					fmt.Println("error while sending stream: ", err)
				}
			}(&payload)
		}
	})

	r.POST("/nats", func(c *gin.Context) {
		payload := common.RequestMessage{Payload: make([]byte, 100_000)}
		payload.Name = fmt.Sprintf("message_%v", c.Query("count"))
		for _, conn := range nats_publisher {
			payload.Time = time.Now().UnixMilli()
			bytes, _ := json.Marshal(&payload)
			go func(conn *nats.Conn, bytes []byte) {
				err := conn.Publish(subject, bytes)
				if err != nil {
					fmt.Println("error publishing message to nats: ", err)
				}
			}(conn, bytes)
			fmt.Println("published value to subject: ", subject)
		}
	})

	r.Run("localhost:9090")
}

func main() {
	startBenchmarking()
	wg.Wait()
	fmt.Printf("%+v\n", stats)
}

func startBenchmarking() {
	frequencystring := os.Getenv("FREQUENCY")
	durationstring := os.Getenv("DURATION")

	frequency, err := strconv.Atoi(frequencystring)
	if err != nil {
		fmt.Println("error in frequency")
		return
	}

	durationint, err := strconv.Atoi(durationstring)
	if err != nil {
		fmt.Println("error in frequency")
		return
	}

	endpoints := []string{"http://localhost:9090/nats"}
	//endpoints := []string{"http://localhost:9090/grpc"}
	//endpoints := []string{"http://localhost:9090/grpc", "http://localhost:9090/nats"}
	rate := vegeta.Rate{Freq: frequency, Per: time.Second}      // change the rate here
	duration := time.Duration(int64(durationint)) * time.Second // change the duration here

	for _, endpoint := range endpoints {
		func(endpoint string) {

			attacker := vegeta.NewAttacker()

			attacker.Attack(func(target *vegeta.Target) error {
				target.URL = fmt.Sprintf("%s?count=%v", endpoint, uuid.New().String())
				target.Method = "POST"
				return nil
			}, rate, duration, endpoint)

		}(endpoint)
	}
}

// EVAL "return redis.call('keys', ARGV[1])" 0 grpc_message_*
//func getNumbers() {
//	rdbClient.Eval(context.Background(), "return redis.call('keys', ARGV[1])", 0, "grpc_*").Slice()
//}
