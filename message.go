package signalr

import (
	"context"
	"encoding/json"
	"fmt"
)

var (
	textMessage   = 1
	statusStarted = 1
)

// Message represents a message sent from the server to the persistent websocket
// connection.
type Message struct {
	// message id, present for all non-KeepAlive messages
	C string

	// an array containing actual data
	M []ClientMsg

	// indicates that the transport was initialized (a.k.a. init message)
	S int

	// groups token – an encrypted string representing group membership
	G string

	// other miscellaneous variables that sometimes are sent by the server
	I int `json:"I,string"`
	E string
	R json.RawMessage
	H json.RawMessage // could be bool or string depending on a message type
	D json.RawMessage
	T json.RawMessage
}

// ClientMsg represents a message sent to the Hubs API from the client.
type ClientMsg struct {
	// invocation identifier – allows to match up responses with requests
	I int

	// the name of the hub
	H string

	// the name of the method
	M string

	// arguments (an array, can be empty if the method does not have any
	// parameters)
	A []json.RawMessage

	// state – a dictionary containing additional custom data (optional)
	S *json.RawMessage `json:",omitempty"`
}

// ServerMsg represents a message sent to the Hubs API from the server.
type ServerMsg struct {
	// invocation Id (always present)
	I int

	// the value returned by the server method (present if the method is not
	// void)
	R *json.RawMessage `json:",omitempty"`

	// error message
	E *string `json:",omitempty"`

	// true if this is a hub error
	H *bool `json:",omitempty"`

	// an object containing additional error data (can only be present for
	// hub errors)
	D *json.RawMessage `json:",omitempty"`

	// stack trace (if detailed error reporting (i.e. the
	// HubConfiguration.EnableDetailedErrors property) is turned on on the
	// server)
	T *json.RawMessage `json:",omitempty"`

	// state – a dictionary containing additional custom data (optional)
	S *json.RawMessage `json:",omitempty"`
}

func readMessage(ctx context.Context, conn WebsocketConn, msg *Message, state *State) error {
	t, p, err := conn.ReadMessage(ctx)
	if err != nil {
		return fmt.Errorf("message read failed: %w", err)
	}

	if t != textMessage {
		return fmt.Errorf("unexpected websocket control type: %d", t)
	}

	if err := json.Unmarshal(p, msg); err != nil {
		return err
	}

	// Update the groups token.
	if msg.G != "" {
		state.GroupsToken = msg.G
	}

	// Update the current message ID.
	if msg.C != "" {
		state.MessageID = msg.C
	}

	return nil
}

type negotiateResponse struct {
	URL                     string  `json:"Url"`
	ConnectionToken         string  `json:"ConnectionToken"`
	ConnectionID            string  `json:"ConnectionId"`
	KeepAliveTimeout        float64 `json:"KeepAliveTimeout"`
	DisconnectTimeout       float64 `json:"DisconnectTimeout"`
	ConnectionTimeout       float64 `json:"ConnectionTimeout"`
	TryWebSockets           bool    `json:"TryWebSockets"`
	ProtocolVersion         string  `json:"ProtocolVersion"`
	TransportConnectTimeout float64 `json:"TransportConnectTimeout"`
	LongPollDelay           float64 `json:"LongPollDelay"`
}

type startResponse struct {
	Response string `json:"Response"`
}