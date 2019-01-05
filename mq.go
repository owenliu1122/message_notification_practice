package notice

// Job 消息队列发布载体
type Job struct {
	Message interface{} // Message struct
	Q       string      // Queue name
	D       int         // Delay time
}

type JobInterface interface {
	Queue() string
	Delay() int
	Marshal() ([]byte, error)
}
