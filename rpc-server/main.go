package main

import (
	"log"
	"context"
	rpc "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc/imservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"time"
)


func main() {
	ctx := context.Background()

	if err_db := db_client.InitClient(ctx,"username:password@tcp(IMDB)/im_service"); err_db != nil {
		log.Println("Error initializing MySQL client:", err_db)
		return
	}
	
	retries := 5
	var connectionErr error
	for i := 1; i <= retries; i++ { 
		connectionErr = db_client.TestConnection(ctx) 
		if connectionErr == nil {
			break
		}
		log.Printf("Error connecting to database (attempt %d/%d): %v\n", i, retries, connectionErr)
		time.Sleep(5 * time.Second)
	}

	if connectionErr != nil {
		log.Println("Failed to establish database connection:", connectionErr)
		return
	}

	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"}) // r should not be reused.
	if err != nil {
		log.Fatal(err)
	}

	svr := rpc.NewServer(new(IMServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: "demo.rpc.server",
	}))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}

}



