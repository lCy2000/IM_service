package validation

import (
	"testing"
	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateSendRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  *rpc.SendRequest
		expected error
	}{
		{
			name: "Valid request",
			request: &rpc.SendRequest{
				Message: &rpc.Message{
					Chat:   "john:doe",
					Text:   "Hello",
					Sender: "john",
				},
			},
			expected: nil,
		},
		{
			name: "Request with missing chat",
			request: &rpc.SendRequest{
				Message: &rpc.Message{
					Chat:   "",
					Text:   "Hello",
					Sender: "john",
				},
			},
			expected: errors.New("Please ensure that the fields chat, text, and sender are specified."),
		},
		{
			name: "Request with missing text",
			request: &rpc.SendRequest{
				Message: &rpc.Message{
					Chat:   "john:doe",
					Text:   "",
					Sender: "john",
				},
			},
			expected: errors.New("Please ensure that the fields chat, text, and sender are specified."),
		},
		{
			name: "Request with missing sender",
			request: &rpc.SendRequest{
				Message: &rpc.Message{
					Chat:   "john:doe",
					Text:   "Hello",
					Sender: "",
				},
			},
			expected: errors.New("Please ensure that the fields chat, text, and sender are specified."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendRequest(tt.request)
			if !assert.Equal(t, tt.expected, err){
				t.Errorf("Validation failed, expected: %v, got: %v", tt.expected, err)
			}
		})
	}
}

func TestValidateChatFormat(t *testing.T) {
	chat := "john:doe"

	err := ValidateChatFormat(chat)
	if err != nil {
		t.Errorf("Validation failed, expected no error, got: %v", err)
	}
}

func TestValidateSenderInChat(t *testing.T) {
	chat := "john:doe"
	sender := "john"

	err := ValidateSenderInChat(chat, sender)
	if err != nil {
		t.Errorf("Validation failed, expected no error, got: %v", err)
	}
}

func TestValidatePullRequest(t *testing.T) {
	req := &rpc.PullRequest{
		Chat: "john:doe",
	}

	err := ValidatePullRequest(req)
	if err != nil {
		t.Errorf("Validation failed, expected no error, got: %v", err)
	}
}
