package structures

import "time"

type Room struct {
	ID           uint64 `json:"id" bson:"_id"`
	Block        string `json:"block" bson:"block"`
	Number       int64  `json:"number" bson:"number"` // cabinet number
	SEF          string `json:"sef" bson:"sef"`       // search engines friendly to make search easier and more presentable
	ScheduleHash string `json:"hash" bson:"hash"`
	Deleted      bool
	Updated      time.Time `json:"updated" bson:"updated"`
}

type File struct {
	ID       uint64 `json:"id" bson:"_id"`
	Hash     string
	Name     string
	Path     string
	FileType string
}
