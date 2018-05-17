package json

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
)

func Marshal(msg proto.Message) ([]byte, error) {
	return json.Marshal(msg)
}

func Unmarshal(data []byte, msg proto.Message) error {
	return json.Unmarshal(data, msg)
}
