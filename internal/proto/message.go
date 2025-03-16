package proto

const (
	MsgTypeHello   = "hello"
	MsgTypeSysInfo = "sysinfo"
)

type Message struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}
