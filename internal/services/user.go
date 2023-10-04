package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
	"tutorial-auth/internal/mongodb"
	"tutorial-auth/internal/mongodb/models"
)

var UserAlreadyExistsError = fmt.Errorf("user already exists")

type UserService struct {
	logger      *zap.Logger
	mongoClient *mongodb.MongoDB
	dbClient    *sqlx.DB
	collection  string
}

func NewUserService(logger *zap.Logger, mongoClient *mongodb.MongoDB, db *sqlx.DB) *UserService {
	return &UserService{
		logger:      logger,
		mongoClient: mongoClient,
		dbClient:    db,
		collection:  "users",
	}
}

func (us *UserService) GetByGuid(guid string) (*models.User, error) {
	var user *models.User
	collection := us.mongoClient.GetCollection(us.collection)
	err := collection.FindOne(context.TODO(), bson.M{"guid": guid}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) GetByLogin(login string) (*models.User, error) {
	var user *models.User
	collection := us.mongoClient.GetCollection(us.collection)
	err := collection.FindOne(context.TODO(), bson.M{"login": login}).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserService) Register(login string, password string) (*models.User, error) {
	//var user *models.User
	collection := us.mongoClient.GetCollection(us.collection)

	existedUser, err := us.GetByLogin(login)
	us.logger.Info("checking if user already exists", zap.String("login", login), zap.Error(err))
	if err != nil {
		return nil, err
	}
	if existedUser != nil {
		return nil, UserAlreadyExistsError
	}

	userGUID := uuid.New().String()
	timeNow := time.Now().Local()
	newUser := &models.User{
		GUID:      userGUID,
		Login:     login,
		LoginType: models.LoginType{ID: 1, Name: "email"},
		Name:      "",
		LastName:  "",
		CreatedAt: timeNow,
	}
	insertResult, err := collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return nil, err
	}

	sql := `INSERT INTO passwords (user_id, password, expires_at) VALUES ($1, $2, $3)`
	hashedPass, err := us.HashPassword(password)
	if err != nil {
		us.logger.Error("failed to hashing password", zap.Error(err))
		errDelete := us.DeleteByID(insertResult.InsertedID.(primitive.ObjectID))
		if errDelete != nil {
			return nil, errDelete
		}
		return nil, err
	}
	_, err = us.dbClient.Exec(sql, userGUID, hashedPass, timeNow)
	if err != nil {
		us.logger.Error("failed to inserting password", zap.Error(err))
		errDelete := us.DeleteByID(insertResult.InsertedID.(primitive.ObjectID))
		if errDelete != nil {
			return nil, errDelete
		}
		return nil, err
	}
	return newUser, nil
}

func (us *UserService) DeleteByID(id primitive.ObjectID) error {
	collection := us.mongoClient.GetCollection(us.collection)
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (us *UserService) CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
