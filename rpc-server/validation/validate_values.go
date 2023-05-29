package validation

import (
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"errors"
	"strings"
)

//For send request
func ValidateSendRequest(req *rpc.SendRequest) error {
	if req.Message == nil || req.Message.Chat == "" || req.Message.Text == "" || req.Message.Sender == "" {
		return errors.New("Please ensure that the fields chat, text, and sender are specified.")
	}
	return nil
}

func ValidateChatFormat(chat string) error {
	chatLowercase := strings.ToLower(chat)
	chatParts := strings.Split(chatLowercase, ":")
	if len(chatParts) != 2 {
		return errors.New("Please ensure you specify the sender correctly in the format of a:b")
	}
	return nil
}

func ValidateSenderInChat(chat string, sender string) error {
	chatLowercase := strings.ToLower(chat)
	chatParts := strings.Split(chatLowercase, ":")

	if chatParts[0] != sender && chatParts[1] != sender{
		return errors.New("Please ensure the sender is part of the chat room.")
	}
	return nil
}

//For pull request
func ValidatePullRequest(req *rpc.PullRequest) error {
	if req.Chat == "" {
		return errors.New("Please ensure that the the chat field is specified.")
	}
	return nil
}

