package notice

// UserMessage is a message record that notifies the user.
type UserMessage struct {
	ID uint64 `json:"id"`
	//RecordID    uint64 `json:"record_id"`
	Name        string `json:"name"`
	Content     string `json:"content"`
	NoticeType  string `json:"notice_type"`
	Destination string `json:"destination"`
	//Phone   string `json:"phone"`
	//Email   string `json:"email"`
	//WeChat  string `json:"wechat"`
}
