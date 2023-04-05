package user

import (
	"context"
	"errors"

	database "github.com/PacodiazDG/Backend-blog/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func GetBasicInfo(ID primitive.ObjectID, data interface{}) error {
	collection := *database.Database.Collection("Users")
	err := collection.FindOne(context.Background(), bson.M{"_id": ID}).Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func DelateAccount(ID primitive.ObjectID) error {
	collection := *database.Database.Collection("Users")
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": ID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("id not found")
	}
	return nil
}

func EmailIsAvailable(Email string) bool {
	var result bson.M
	collection := *database.Database.Collection("Users")
	collection.FindOne(context.TODO(), bson.M{"Email": Email}).Decode(&result)
	return len(result) != 0
}

func UpdateUserInfo(id primitive.ObjectID, data *UserStrcture) (*mongo.UpdateResult, error) {
	collection := *database.Database.Collection("Users")
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": data})
	if result.ModifiedCount == 0 {
		return nil, errors.New("no document was updated")
	}
	if err != nil {
		return nil, errors.New("error in database")
	}
	return result, nil
}

func IPAddrUser(data *IpAddrUser) (bool, error) {
	collection := *database.Database.Collection("IPAddrUser")
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func RemoveIPAddrUser(uudi string) error {
	collection := *database.Database.Collection("IPAddrUser")
	removed, err := collection.DeleteOne(context.TODO(), bson.M{"Uuidtoken": uudi})
	if removed.DeletedCount == 0 {
		return errors.New("error apparently the uuid does not exist")
	}
	if err != nil {
		return errors.New("error")
	}
	return nil
}

func GetIpaddes(Userid primitive.ObjectID) ([]IpAddrUser, error) {
	var ipaddress []IpAddrUser
	collection := *database.Database.Collection("IPAddrUser")
	results, err := collection.Find(context.TODO(), bson.M{"IDuser": Userid})
	if err != nil {
		return nil, errors.New("error in database")
	}
	if err = results.All(context.Background(), &ipaddress); err != nil {
		return nil, errors.New("error in database")
	}
	return ipaddress, nil
}
