package main

import (
	"context"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/db"
	validate "github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/validation"
	"time"
	"strings"
)

type IMServiceImpl struct{}

func createTimestamp() int64 {
	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()
	return int64(unixTimestamp)
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()
	resp.Code = 0
	resp.Msg = "Insert values successfully!"

	if err := validate.ValidateSendRequest(req); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, err
	}

	chat := req.Message.Chat
	chat = strings.ReplaceAll(chat, " ", "")
	text := req.Message.Text
	sender := req.Message.Sender
	sendtime := createTimestamp()
	
	if err := validate.ValidateChatFormat(chat); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, err
	}

	if err := validate.ValidateSenderInChat(chat, sender); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, err
	}	

	if err := db.InsertMessage(db_client,chat, sender, text, sendtime);err != nil{
		resp.Code = 500
		resp.Msg = "Failed to insert values"
		return resp, err
	}

	return resp, nil
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Msg = "Success"
	resp.Code = 0
	

	err := validate.ValidatePullRequest(req)
	if err != nil{
		resp.Code = 500
		resp.Msg = err.Error()
		return resp,err
	}

	chat := req.Chat
	chat = strings.ReplaceAll(chat, " ", "") //strip off a whitespace characters in chat
	cursor := req.Cursor
	limit := req.Limit
	
	if err := validate.ValidateChatFormat(chat); err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, err
	}

	//set default limit to be 10 (if limit is not specified, default 0)
	if req.Limit == 0{
		limit = 10
	}
	reverse := req.Reverse
	reverseSetting := "ASC"
	if reverse != nil && *reverse {
		reverseSetting = "DESC"
	}

	rows,err := db.GetMessages(db_client,reverseSetting, chat, cursor)
	if err != nil {
		resp.Code = 500
		resp.Msg = err.Error()
		return resp, err
	}

	messages,nextcursor, hasMore, err := db.GetProperties(rows, limit)
	resp.Messages = messages	
	resp.NextCursor = &nextcursor	
	resp.HasMore = &hasMore

	defer rows.Close()
	
	return resp, nil
}


