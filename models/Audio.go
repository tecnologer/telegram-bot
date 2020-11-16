package models

//Audio This object represents an audio file to be treated as music by the Telegram clients.
type Audio struct {
	*Document
	Duration  int    `json:"duration"`
	Performer string `json:"performer"`
	Title     string `json:"title"`
}
