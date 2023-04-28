package main

import (
	"flag"
	"log"

	"github.com/ryanreadbooks/go-grpc-example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 创建客户端
func InitClient(target string) (pb.CellphoneServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not dial to %s, failure connection: %v\n", target, err)
	}
	return pb.NewCellphoneServiceClient(conn), conn
}

func main() {
	target := flag.String("target", "", "the target of the grpc server")
	flag.Parse()

	log.Printf("dialing to %s\n", *target)
	client, conn := InitClient(*target)
	defer conn.Close()
	log.Println(client)
	// TODO
}
