package models

//PassportData Contains information about Telegram Passport data shared with the bot by the user.
type PassportData struct {
	Data        []*EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials        `json:"credentials"`
}
