package delayer

// 消息结构
type Message struct {
	ID    string
	Topic string
	Body  string
}

// 效验
func (p *Message) Valid() bool {
	if p.ID == "" || p.Topic == "" || p.Body == "" {
		return false
	}
	return true
}
