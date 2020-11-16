package models

//Video This object represents a video file.
type Video struct {
	*Document
	Width    int `json:"width"`
	Height   int `json:"height"`
	Duration int `json:"duration"`
}
