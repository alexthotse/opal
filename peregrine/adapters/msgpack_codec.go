package adapters

import (
	"github.com/vmihailenco/msgpack/v5"
	"connectrpc.com/connect"
)

const MsgPackCodecName = "msgpack"

// msgPackCodec implements connect.Codec
type msgPackCodec struct{}

func NewMsgPackCodec() connect.Codec {
	return &msgPackCodec{}
}

func (c *msgPackCodec) Name() string {
	return MsgPackCodecName
}

func (c *msgPackCodec) Marshal(message any) ([]byte, error) {
	return msgpack.Marshal(message)
}

func (c *msgPackCodec) Unmarshal(data []byte, message any) error {
	return msgpack.Unmarshal(data, message)
}

func (c *msgPackCodec) MarshalAppend(dst []byte, message any) ([]byte, error) {
	b, err := msgpack.Marshal(message)
	if err != nil {
		return dst, err
	}
	return append(dst, b...), nil
}

func (c *msgPackCodec) IsBinary() bool {
	return true
}
