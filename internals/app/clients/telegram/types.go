package telegram

type MessageWrapper struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Id      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type Message struct {
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type User struct {
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}
