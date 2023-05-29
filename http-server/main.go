package main

import (
	"context"
	"log"
	"time"	
	"strconv"

	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/kitex_gen/rpc"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/kitex_gen/rpc/imservice"
	"github.com/TikTokTechImmersion/assignment_demo_2023/http-server/proto_gen/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/hertz-contrib/cors"

)

var cli imservice.Client

func main() {
	r, err := etcd.NewEtcdResolver([]string{"etcd:2379"})
	if err != nil {
		log.Fatal(err)
	}
	cli = imservice.MustNewClient("demo.rpc.server",
		client.WithResolver(r),
		client.WithRPCTimeout(5*time.Second),
		client.WithHostPorts("rpc-server:8888"),
	)

	h := server.Default(server.WithHostPorts("0.0.0.0:8080"))
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
    AllowHeaders:     []string{"Origin", "Content-Type"}, // Include "Content-Type" header
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
	},
		MaxAge: 12 * time.Hour,
	}))

	h.GET("/ping", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{"message": "pong"})
	})

	h.POST("/api/send", sendMessage)
	h.GET("/api/pull", pullMessage)
	h.GET("/api/pull/messages", pullMessageClient)

	h.Spin() //start the server to listen for http requests
}

func sendMessage(ctx context.Context, c *app.RequestContext) {

	var req api.SendRequest 

	err := c.Bind(&req) 

	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body: %v", err)
		return
	} 

	message := &rpc.Message{
		Chat:   req.Chat,
		Text:   req.Text,
		Sender: req.Sender,
	}

	resp, err := cli.Send(ctx, &rpc.SendRequest{
		Message: message,
	})


	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
	} 
	
	c.String(consts.StatusOK, resp.Msg)

}

func pullMessage(ctx context.Context, c *app.RequestContext) {
	var req api.PullRequest
	err := c.Bind(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, "Failed to parse request body: %v", err)
		return
	}

	resp, err := cli.Pull(ctx, &rpc.PullRequest{
		Chat:    req.Chat,
		Cursor:  req.Cursor,
		Limit:   req.Limit,
		Reverse: &req.Reverse,
	})

	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
		return
	} 

	messages := make([]*api.Message, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, &api.Message{
			Chat:     msg.Chat,
			Text:     msg.Text,
			Sender:   msg.Sender,
			SendTime: msg.SendTime,
		})
	}


	c.JSON(consts.StatusOK, &api.PullResponse{
		Messages:   messages,
		HasMore:    resp.GetHasMore(),
		NextCursor: resp.GetNextCursor(),
	})
}

func pullMessageClient(ctx context.Context, c *app.RequestContext) {
	var req api.PullRequest

	// Read from request parameters
	req.Chat = c.Query("chat")
	cursorStr := c.Query("cursor")
	limitStr := c.Query("limit")
	reverseStr := c.Query("reverse")

	// Parse cursor as int64
	cursor, err := strconv.ParseInt(cursorStr, 10, 64)
	if err != nil {
		c.String(consts.StatusBadRequest, "Invalid value for 'cursor' parameter")
		return
	}
	req.Cursor = cursor

	// Parse limit as int
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.String(consts.StatusBadRequest, "Invalid value for 'limit' parameter")
		return
	}
	req.Limit = int32(limit)

	// Parse reverse as boolean
	reverse, err := strconv.ParseBool(reverseStr)
	if err != nil {
		c.String(consts.StatusBadRequest, "Invalid value for 'reverse' parameter")
		return
	}
	req.Reverse = reverse

	resp, err := cli.Pull(ctx, &rpc.PullRequest{
		Chat:    req.Chat,
		Cursor:  req.Cursor,
		Limit:   req.Limit,
		Reverse: &req.Reverse,
	})

	if err != nil {
		c.String(consts.StatusInternalServerError, err.Error())
		return
	} else if resp.Code != 0 {
		c.String(consts.StatusInternalServerError, resp.Msg)
		return
	} 

	messages := make([]*api.Message, 0, len(resp.Messages))
	for _, msg := range resp.Messages {
		messages = append(messages, &api.Message{
			Chat:     msg.Chat,
			Text:     msg.Text,
			Sender:   msg.Sender,
			SendTime: msg.SendTime,
		})
	}


	c.JSON(consts.StatusOK, &api.PullResponse{
		Messages:   messages,
		HasMore:    resp.GetHasMore(),
		NextCursor: resp.GetNextCursor(),
	})
}
