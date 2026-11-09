// Package contract holds the formats crow-core shares with the Python CLI, the Python
// sidecars and crow-deck. Spec: logics/specs/core_contract and logics/specs/req_037_run_contract.
//
// The package does no I/O beyond decoding the bytes it is given, and imports only the
// standard library and yaml.v3. crow-deck may import it; everything else in crow-core
// lives under internal/.
//
// Optional fields are pointers with omitempty and nullable fields are pointers without it,
// so a decode followed by an encode gives back the same JSON.
package contract

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Version is the contract version: the v of every sidecar message and the contract field
// of crow blocks and questions files.
const Version = 1

// Error codes of a failed reply.
const (
	ErrBadRequest          = "bad_request"
	ErrUnknownType         = "unknown_type"
	ErrUnsupportedContract = "unsupported_contract"
	ErrTimeout             = "timeout"
	ErrInternal            = "internal"
)

// Request is one line core writes to a sidecar's stdin.
type Request struct {
	V    int             `json:"v"`
	ID   string          `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Reply is one line a sidecar writes to its stdout. Data is set when OK, Error when not.
type Reply struct {
	V     int             `json:"v"`
	ID    string          `json:"id"`
	Type  string          `json:"type"`
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data,omitempty"`
	Error *Error          `json:"error,omitempty"`
}

// Error is the error of a failed reply.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HelloRequest struct {
	Contract int    `json:"contract"`
	Core     string `json:"core"`
}

type HelloReply struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Contract int      `json:"contract"`
	Messages []string `json:"messages"`
}

type HealthRequest struct{}

type HealthReply struct {
	Status string  `json:"status"`
	Detail *string `json:"detail,omitempty"`
}

type ShutdownRequest struct {
	GraceS int `json:"grace_s"`
}

type ShutdownReply struct{}

// Message pairs a message type with its request and reply data types.
type Message struct {
	Type    string
	Request func() any
	Reply   func() any
}

// Messages lists every message of the contract. Each sidecar answers the lifecycle three.
var Messages = []Message{
	{"hello", func() any { return &HelloRequest{} }, func() any { return &HelloReply{} }},
	{"health", func() any { return &HealthRequest{} }, func() any { return &HealthReply{} }},
	{"shutdown", func() any { return &ShutdownRequest{} }, func() any { return &ShutdownReply{} }},
}

// LookupMessage returns the message of that type.
func LookupMessage(t string) (Message, bool) {
	for _, m := range Messages {
		if m.Type == t {
			return m, true
		}
	}
	return Message{}, false
}

// Decode reads exactly one JSON value into v and refuses unknown fields.
func Decode(b []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return fmt.Errorf("trailing data after the JSON value")
	}
	return nil
}

// DecodeRequestData decodes the data of a request into its typed struct.
func DecodeRequestData(r Request) (any, error) {
	m, ok := LookupMessage(r.Type)
	if !ok {
		return nil, fmt.Errorf("%s: unknown message type: %s", ErrUnknownType, r.Type)
	}
	v := m.Request()
	if err := Decode(r.Data, v); err != nil {
		return nil, fmt.Errorf("%s: %s data: %w", ErrBadRequest, r.Type, err)
	}
	return v, nil
}

// DecodeReplyData decodes the data of a successful reply into its typed struct.
func DecodeReplyData(r Reply) (any, error) {
	m, ok := LookupMessage(r.Type)
	if !ok {
		return nil, fmt.Errorf("unknown message type: %s", r.Type)
	}
	if !r.OK {
		return nil, fmt.Errorf("%s reply failed: %s: %s", r.Type, r.Error.Code, r.Error.Message)
	}
	v := m.Reply()
	if err := Decode(r.Data, v); err != nil {
		return nil, fmt.Errorf("%s reply data: %w", r.Type, err)
	}
	return v, nil
}
