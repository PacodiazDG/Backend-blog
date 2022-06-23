package notebook

import "time"

type NoteStruct struct {
	Title   string    `bson:"Title,omitempty"`
	Content string    `bson:"Content,omitempty"`
	Tags    []string  `bson:"Tags,omitempty"`
	Date    time.Time `bson:"Date,omitempty" `
	Author  string    `bson:"Author,omitempty"`
	ID      string    `bson:"_id,omitempty" `
}
