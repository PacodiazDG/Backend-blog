package admin

import (
	"context"
	"errors"
	"regexp"

	"github.com/PacodiazDG/Backend-blog/api/v1/user"
	database "github.com/PacodiazDG/Backend-blog/database"
	logs "github.com/PacodiazDG/Backend-blog/modules/logs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func ListUsers(next int64, username string) ([]user.BasicInfo, error) {
	var Listofusers []user.BasicInfo
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
		logs.WriteLogs(err, logs.LowError)
		return nil, errors.New("error in database")
	}
	return Listofusers, nil
}
