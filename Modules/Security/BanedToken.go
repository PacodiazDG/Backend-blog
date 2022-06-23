package Security

import (
	"context"
	"encoding/hex"

	database "github.com/PacodiazDG/Backend-blog/Database"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/sha3"
)

func IsBanedToken(Token string) bool {
	h := sha3.New224()
	h.Write([]byte(Token))
	sha3hash := hex.EncodeToString(h.Sum(nil))
	collection := *database.Database.Collection("Users")
	err := collection.FindOne(context.TODO(), bson.M{"Email": sha3hash})
	return err != nil
}
