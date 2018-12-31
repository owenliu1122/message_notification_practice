package notice

// UserMessage is a message record that notifies the user.
type UserMessage struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	WeChat  string `json:"wechat"`
}
