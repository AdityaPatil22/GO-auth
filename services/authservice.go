package services

import (
	"GO-temp-backend/config"
	"GO-temp-backend/models"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func InitServices() {
	userCollection = config.MongoClient.Database("meme_gen").Collection("users")
}
func HashPassowrd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func GenerateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func SignupUser(user models.User) (string, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("User Already Exists")
	}

	hashedPassword, err := HashPassowrd(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID)
	idStr := id.Hex()
	return GenerateJWT(idStr)
}

func LoginUser(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := VerifyPassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	return GenerateJWT(user.ID.Hex())
}
