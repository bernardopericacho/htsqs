package aws

import (
	"encoding/json"
	"strconv"
	"time"
)

// SqsJSONMessage is the SQS message response in JSON format
// Message is processed as json.RawMessage for later use
type SqsJSONMessage struct {
	MessageID string          `json:"MessageId"`
	Type      string          `json:"Type"`
	Timestamp time.Time       `json:"Timestamp"`
	Message   json.RawMessage `json:"Message"`
	TopicArn  string          `json:"TopicArn"`
}

// UnmarshalJSON allows Node to implement the json.Unmarshaler interface and
// unmarshal itself from a json structure
func (sjm *SqsJSONMessage) UnmarshalJSON(b []byte) (err error) {
	type Alias SqsJSONMessage
	v := &Alias{}
	err = json.Unmarshal(b, v)

	m, err := strconv.Unquote(string(v.Message))

	*sjm = (SqsJSONMessage)(*v)
	sjm.Message = json.RawMessage(m)

	return err
}
