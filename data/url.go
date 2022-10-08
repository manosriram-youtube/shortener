package data

import "time"

type UrlShortener struct {
	Src  string `json:"src" bson:"src"`
	Dest string `json:"dest" bson:"dest"`
	// todo: update time.Time type to something accurate
	Created_at time.Time
	Expires_at time.Time
	Hits       int
}

type ShortenerRequest struct {
	Src string `json:"src" bson:"src"`
}
