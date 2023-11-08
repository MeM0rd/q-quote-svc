package quote

type Quote struct {
	Id     int    `json:"id"`
	UserId int    `json:"user_id"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}
