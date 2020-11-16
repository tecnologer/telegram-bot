package models

//Animation This object represents an animation file (GIF or H.264/MPEG-4 AVC video without sound).
type Animation struct {
	*Document
	Width    int `json:"width"`
	Height   int `json:"height"`
	Duration int `json:"duration"`
}
