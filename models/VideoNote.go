package models

//VideoNote This object represents a video message (available in Telegram apps as of v.4.0).
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumb        *PhotoSize `json:"thumb"`
	FileSize     int        `json:"file_size"`
	Duration     int        `json:"duration"`
	Length       int        `json:"length"`
}
