package model

const (
	MsgTypeMail   = "mail"
	MsgTypePhone  = "phone"
	MsgTypeWeChat = "wechat"
)

type UserMsg struct {
	ID      uint64
	Name    string
	Content string
	Type    string
	Phone   string
	Email   string
	WeChat  string
}
