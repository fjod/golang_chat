package chat

import (
	DB "chat/db"
	Mocks "chat/db/mocks"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"sync"
	"testing"
)

func TestHub(t *testing.T) {
	sentMsg := make([]string, 0, 2)
	var wg sync.WaitGroup
	h := NewHub()
	dbService := Mocks.NewRepoOperations(t)
	dbService.On("Append", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
		msg := args.Get(0)
		sentMsg = append(sentMsg, fmt.Sprint(msg))
		wg.Done()
	})

	go h.Run(DB.NewService(dbService))

	h.broadcast <- []byte("test1")
	h.broadcast <- []byte("test2")

	wg.Add(2)
	wg.Wait() // wait for hub to broadcast messages
	assert.True(t, containsString(sentMsg, "test1"))
	assert.True(t, containsString(sentMsg, "test2"))
}

func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}
