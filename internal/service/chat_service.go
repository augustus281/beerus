package service

import (
	"io"
	"log"
	"sync"

	"github.com/google/uuid"

	"github.com/augustus281/beerus/proto/gen/go/beerus"
)

type ChatService struct {
	beerus.UnimplementedChatServiceServer
	users map[string]beerus.ChatService_ChatServer
	mu    sync.Mutex
}

func NewChatService() *ChatService {
	return &ChatService{
		users: make(map[string]beerus.ChatService_ChatServer),
	}
}

func (c *ChatService) Chat(stream beerus.ChatService_ChatServer) error {
	clientID := uuid.New().String()

	// Register the user
	c.mu.Lock()
	c.users[clientID] = stream
	c.mu.Unlock()

	defer func() {
		// Unregister the user when the stream is closed
		c.mu.Lock()
		delete(c.users, clientID)
		c.mu.Unlock()
	}()

	for {
		// Receive a message from the client
		req, err := stream.Recv()
		if err == io.EOF {
			// Client closed the stream
			return nil
		}
		if err != nil {
			return err
		}

		// Broadcast the message to all connected users
		c.mu.Lock()
		for addr, userStream := range c.users {
			if addr == clientID {
				continue
			}
			err := userStream.Send(&beerus.ChatResponse{
				User: req.Name,
				Text: req.Text,
			})
			if err != nil {
				log.Printf("Error sending message to %s: %v\n", addr, err)
			}
		}
		c.mu.Unlock()
	}
}
