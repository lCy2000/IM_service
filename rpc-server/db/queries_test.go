package db

import (
	"testing"
	"math/rand"
	"time"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertMessage(t *testing.T){
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock SQL client: %v", err)
	}
	defer mockdb.Close()

	client := MySQLClient{db: mockdb}

	mock.ExpectBegin()
	mock.ExpectExec("Insert messages").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	rand.Seed(time.Now().UnixNano())
	
	testCases := []map[string]interface{}{
    {
			"testData": map[string]interface{}{
				"chat":     "joe:doe",
				"sender":   "john",
				"message":  "hello",
				"sendtime": rand.Int63(),
			},
			"expected": nil,
    },
		{
			"testData": map[string]interface{}{
				"chat":     "",
				"sender":   "john",
				"message":  "hello",
				"sendtime": rand.Int63(),
			},
			"expected": nil,
    },
	}
	for _, tt := range testCases {
		testData := tt["testData"].(map[string]interface{})
		expected := tt["expected"]
	
		chat := testData["chat"].(string)
		sender := testData["sender"].(string)
		message := testData["message"].(string)
		sendTime := testData["sendtime"].(int64)
	
		err := InsertMessage(client, chat, sender, message, sendTime)
		fmt.Println(err)
		if expected == nil {
			assert.NoError(t, err)
		} else {
			t.Fatalf("Failed to create mock SQL client: %v", err)
			assert.Error(t, err)
		}
	}
}
