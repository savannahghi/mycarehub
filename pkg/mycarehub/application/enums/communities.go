package enums

import (
	"fmt"
	"io"
	"strconv"
)

// MessageType is the type of message of a given getsream message
type MessageType string

const (
	// MessageTypeRegular is the type of a regular message
	MessageTypeRegular MessageType = "regular"
	// MessageTypeError is the type of an error message
	MessageTypeError MessageType = "error"
	// MessageTypeReply is the type of a reply message
	MessageTypeReply MessageType = "reply"
	// MessageTypeSystem is the type of a system message
	MessageTypeSystem MessageType = "system"
	// MessageTypeEphemeral is the type of an ephemeral message
	MessageTypeEphemeral MessageType = "ephemeral"
)

// AllMessageType is a valid list of message types
var AllMessageType = []MessageType{
	MessageTypeRegular,
	MessageTypeError,
	MessageTypeReply,
	MessageTypeSystem,
	MessageTypeEphemeral,
}

// IsValid returns true if the MessageType is a valid value
func (e MessageType) IsValid() bool {
	switch e {
	case MessageTypeRegular, MessageTypeError, MessageTypeReply, MessageTypeSystem, MessageTypeEphemeral:
		return true
	}
	return false
}

// String is the string representation of the MessageType
func (e MessageType) String() string {
	return string(e)
}

// UnmarshalGQL parses the input string as a MessageType
func (e *MessageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageType", str)
	}
	return nil
}

// MarshalGQL serializes the MessageType
func (e MessageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
