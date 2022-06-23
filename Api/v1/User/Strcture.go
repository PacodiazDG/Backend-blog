package User

import (
	"time"

	"github.com/PacodiazDG/Backend-blog/Modules/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserStrcture
type UserStrcture struct {
	Email       string    `bson:"Email,omitempty"`
	Image       string    `bson:"Image,omitempty"`
	Password    string    `json:"Password,omitempty" bson:"Password,omitempty"`
	FirstName   string    `bson:"FirstName,omitempty"`
	LastName    string    `bson:"LastName,omitempty"`
	Permissions string    `bson:"Permissions,omitempty"`
	Username    string    `bson:"Username,omitempty"`
	Banned      bool      `bson:"Banned"`
	Created_at  time.Time `bson:"Created_at,omitempty"`
	ID          string    `bson:"_id,omitempty"`
}

// BasicInfo
type BasicInfo struct {
	Username    string `bson:"Username,omitempty"`
	FirstName   string `bson:"FirstName,omitempty"`
	LastName    string `bson:"LastName,omitempty"`
	Permissions string `bson:"Permissions,omitempty"`
	Email       string `bson:"Email,omitempty"`
	ID          string `bson:"_id,omitempty"`
	Image       string `bson:"Image,omitempty"`
	Banned      bool   `bson:"Banned,omitempty"`
}

//IpAddrUser
type IpAddrUser struct {
	IDuser    primitive.ObjectID `bson:"IDuser"`
	Date      time.Time          `bson:"Date"`
	IpAddrs   string             `bson:"IpAddrs"`
	DateOut   time.Time          `bson:"DateOut"`
	UserAgent string             `bson:"UserAgent"`
	Uuidtoken string             `bson:"Uuidtoken"`
}

//LoginRequestStruct
type LoginRequestStruct struct {
	ID       uint64 `json:"id"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

//Data
type Data []UserStrcture

//RecoveryAccountStrcture
type RecoveryAccountStrcture struct {
	Name     string
	UrlToken string
}

//IsValidUserInfo
func IsValidUserInfo(Info *UserStrcture) bool {
	if !validation.IsValidEmail(Info.Email) {
		return false
	} else if len(Info.Password) > 10 {
		return false
	}
	return true
}
