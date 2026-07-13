package api

import "encoding/json"

type Message struct {
	Text string `json:"text"`
}

// NewMessage returns a JSON message body.
func NewMessage(text string) []byte {
	m := Message{
		Text: text,
	}

	v, err := json.Marshal(m)
	if err != nil {
		return []byte(`{"text":"failed to build message"}`)
	}

	return v
}
