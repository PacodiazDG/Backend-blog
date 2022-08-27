package admin

import (
	"context"
	"errors"
	"regexp"

	"github.com/PacodiazDG/Backend-blog/Api/v1/User"
	database "github.com/PacodiazDG/Backend-blog/Database"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func ListUsers(next int64, username string) ([]User.BasicInfo, error) {
	var Listofusers []User.BasicInfo
	collection := *database.Database.Collection("Users")
	options := options.Find()
	options.SetSort(bson.M{"_id": -1})
	options.SetSkip(next)
	options.SetLimit(10)
	results, err := collection.Find(context.TODO(), bson.M{"Username": bson.M{"$regex": `(?i)` + regexp.QuoteMeta(username)}}, options)
	if err != nil {
		return nil, errors.New("error in database")
	}
	if err = results.All(context.Background(), &Listofusers); err != nil {
		return nil, errors.New("error in database")
	}
	return Listofusers, nil
}
