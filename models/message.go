package models

// Message contains info
type Message struct {
	UserID    int64
	QueueName string
	Content   string
	Index     int64
	Recievers []int64
}

// AddReceiver adds reciever for message
func (m *Message) AddReceiver(clientID int64) {
	m.Recievers = append(m.Recievers, clientID)
}
