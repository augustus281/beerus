package wrapper

import (
	"io"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	underlyinConnection *websocket.Conn

	// CloseTimeout is the amount of time to wait after sending a closure status.
	CloseTimeout time.Duration
}

func New(conn *websocket.Conn) *Connection {
	return &Connection{
		underlyinConnection: conn,
		CloseTimeout:        time.Second * 30,
	}
}

func (c *Connection) SendMessage(data []byte) error {
	return c.underlyinConnection.WriteMessage(websocket.TextMessage, data)
}

func (c *Connection) SendClose() error {
	return c.underlyinConnection.WriteControl(websocket.CloseMessage, nil, time.Now().Add(c.CloseTimeout))
}

func (c *Connection) ReceiveMessage() ([]byte, error) {
	_, data, err := c.underlyinConnection.ReadMessage()
	if err != nil && isClosedConnectionError(err) {
		return nil, io.EOF
	}
	return data, err
}

func (c *Connection) Close() error {
	err := c.underlyinConnection.Close()
	if err != nil && isClosedConnectionError(err) {
		return err
	}
	return nil
}

func isClosedConnectionError(err error) bool {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived, websocket.CloseGoingAway) {
		return true
	}
	return strings.Contains(err.Error(), "use of closed network connection")
}
