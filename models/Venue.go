package models

//Venue This object represents a venue.
type Venue struct {
	Location       *Location `json:"location"`
	Title          string    `json:"title"`
	Address        string    `json:"address"`
	FoursquereID   string    `json:"foursquare_id"`
	FoursquereType string    `json:"foursquare_type"`
}
