package msg

import "encoding/json"

type Message struct {
	Msg          string
	OwnerAddress string
}

func New(ownerAddress, msg string) *Message {
	return &Message{
		Msg:          msg,
		OwnerAddress: ownerAddress,
	}
}

func (m *Message) ParseIntoByte() ([]byte, error) {
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return msgBytes, nil
}
