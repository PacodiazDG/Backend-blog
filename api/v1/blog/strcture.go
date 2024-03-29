package blog

import (
	"errors"
	"time"
)

// Complete structure of the stories
type StoryStruct struct {
	Title         string    `bson:"Title,omitempty"`
	Content       string    `bson:"Content,omitempty"`
	Visible       bool      `bson:"Visible"`
	Tags          []string  `bson:"Tags,omitempty"`
	Date          time.Time `bson:"Date,omitempty" `
	Imagen        string    `bson:"Imagen,omitempty"`
	Status        string    `bson:"Status,omitempty"`
	Author        string    `bson:"Author,omitempty"`
	Description   string    `bson:"Description,omitempty"`
	Password      string    `bson:"Password"`
	UrlImageFound []string  `bson:"UrlImageFound"`
	Views         int64     `bson:"Views"`
	ID            string    `bson:"_id,omitempty" `
	Referal       string    `bson:"Referal,omitempty"`
	Folder        string    `bson:"Folder,omitempty"`
}

// Feed structure
type FeedStrcture struct {
	Title       string    `bson:"Title,omitempty" `
	Author      string    `bson:"Author,omitempty"  `
	ID          string    `bson:"_id,omitempty" `
	Date        time.Time `bson:"Date,omitempty"`
	Imagen      string    `bson:"Imagen,omitempty"`
	Description string    `bson:"Description,omitempty"`
	Visible     bool      `bson:"Visible,omitempty"`
}

// valid if the stories' structure is valid
func IsValidStruct(result *StoryStruct) error {
	if len(result.Title) <= 5 {
		return errors.New("the title of content is very short ")
	} else if len(result.Content) <= 20 {
		return errors.New("the title of content is very short ")
	} else if len(result.Description) <= 10 {
		return errors.New("the lenght of description is very short")
	}
	return nil
}

const DefaultLimit int64 = 10
