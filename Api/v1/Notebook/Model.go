package notebook

import (
	"context"
	"fmt"

	database "github.com/PacodiazDG/Backend-blog/Database"
)

func modelInsertNote(result *NoteStruct) (bool, error) {
	collection := *database.Database.Collection("Notes")
	info, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		return false, err
	}
	fmt.Printf("info.InsertedID: %v\n", info.InsertedID)
	return true, nil
}
func GetNote(result *NoteStruct) (bool, error) {
	collection := *database.Database.Collection("Notes")
	info, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		return false, err
	}
	fmt.Printf("info.InsertedID: %v\n", info.InsertedID)
	return true, nil
}
