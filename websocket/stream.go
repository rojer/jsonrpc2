// Package websocket provides WebSocket transport support for JSON-RPC
// 2.0.
package websocket

import (
	"io"

	"github.com/gorilla/websocket"
)

// A ObjectStream is a mgrpc.ObjectStream that uses a WebSocket to
// send and receive JSON-RPC 2.0 objects.
type ObjectStream struct {
	conn *websocket.Conn
}

// NewObjectStream creates a new mgrpc.ObjectStream for sending and
// receiving JSON-RPC 2.0 objects over a WebSocket.
func NewObjectStream(conn *websocket.Conn) ObjectStream {
	return ObjectStream{conn: conn}
}

// WriteObject implements mgrpc.ObjectStream.
func (t ObjectStream) WriteObject(obj interface{}) error {
	return t.conn.WriteJSON(obj)
}

// ReadObject implements mgrpc.ObjectStream.
func (t ObjectStream) ReadObject(v interface{}) error {
	err := t.conn.ReadJSON(v)
	if e, ok := err.(*websocket.CloseError); ok {
		if e.Code == websocket.CloseAbnormalClosure && e.Text == io.ErrUnexpectedEOF.Error() {
			// Suppress a noisy (but harmless) log message by
			// unwrapping this error.
			err = io.ErrUnexpectedEOF
		}
	}
	return err
}

// Close implements mgrpc.ObjectStream.
func (t ObjectStream) Close() error {
	return t.conn.Close()
}
