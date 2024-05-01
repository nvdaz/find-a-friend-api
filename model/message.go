package model

type Message struct {
	Id         string `json:"id"`
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
}
