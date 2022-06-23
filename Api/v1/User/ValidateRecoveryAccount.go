package User

import (
	"context"
	"net/http"

	database "github.com/PacodiazDG/Backend-blog/Database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ValidateRecoveryAccount s
func ValidateRecoveryAccount(c *gin.Context) {
	collection := *database.Database.Collection("RecoveryAccount")
	Token := c.Param("Token")
	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"Token": Token}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		panic(err)
	}
	if result["Used"] == true {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	update := bson.M{"$set": bson.M{"Used": true}}
	res, err := collection.UpdateOne(context.Background(), bson.M{"_id": result["_id"]}, update)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusAccepted, bson.M{"": res})
}
