package blog

import (
	"context"
	"errors"

	"github.com/PacodiazDG/Backend-blog/database"
	"github.com/PacodiazDG/Backend-blog/modules/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Queryconf struct {
	Collection string
}

// Get the feed
func (v *Queryconf) GetFeed(next int64, query bson.M, limit int64) ([]FeedStrcture, error) {
	findOptions := options.Find()
	collection := *database.Database.Collection(v.Collection)
	findOptions.SetSort(bson.M{"_id": -1})
	findOptions.SetLimit(limit)
	findOptions.SetSkip(next * 11)
	var results []FeedStrcture
	cursor, err := collection.Find(context.Background(), query, findOptions)
	if err != nil {
		logs.WriteLogs(err, logs.LowError)
		return []FeedStrcture{}, errors.New("internal server error")
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		logs.WriteLogs(err, logs.MediumError)
		return []FeedStrcture{}, errors.New("internal server error")
	}
	return results, nil
}

// Insert a post from a structure
func (v *Queryconf) ModelInsertPost(result *PostStruct) (*mongo.InsertOneResult, error) {
	collection := *database.Database.Collection(v.Collection)
	info, err := collection.InsertOne(context.Background(), result)
	if err != nil {
		logs.WriteLogs(err, logs.LowError)
		return nil, err
	}
	return info, nil
}

// Gets a post from an id
func (v *Queryconf) ModelGetArticle(objectId primitive.ObjectID) (PostStruct, error) {
	collection := *database.Database.Collection(v.Collection)
	var result PostStruct
	err := collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return PostStruct{}, errors.New("post not found")
		}
		return PostStruct{}, errors.New("post not found")
	}
	return result, nil
}

// Update a post from a structure
func (v *Queryconf) ModelUpdate(dataInsert *PostStruct, objectId primitive.ObjectID) (bool, error) {
	collection := *database.Database.Collection(v.Collection)
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": dataInsert})
	if err != nil {
		return false, err
	}
	return true, nil
}

// Delete a specific post
func (v *Queryconf) DelatePost(id primitive.ObjectID) error {
	collection := *database.Database.Collection(v.Collection)
	cursor, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return errors.New("internal server error")
	}
	if cursor.DeletedCount == 0 {
		return errors.New("there is no post with that id")
	}
	return nil
}

// Adds the number of visits from a specific post
func (v *Queryconf) Addviews(id primitive.ObjectID) error {
	_, err := (*database.Database.Collection(v.Collection)).UpdateOne(context.Background(), bson.M{
		"_id": id}, bson.M{"$inc": bson.M{"Views": 1}})
	if err != nil {
		return err
	}
	return nil
}

// Get the last ten posts most viewed
func (v *Queryconf) GetTOP() ([]PostStruct, error) {
	collection := *database.Database.Collection(v.Collection)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"Views": -1})
	findOptions.SetLimit(50)
	var results []PostStruct
	cursor, err := collection.Find(context.Background(), bson.M{"Password": "", "Visible": true}, findOptions)
	if err != nil {
		return []PostStruct{}, err
	}
	if err = cursor.All(context.Background(), &results); err != nil {
		return []PostStruct{}, err
	}
	return results, nil
}
