package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"tutorial-auth/internal/config"
)

type MongoDB struct {
	Client    *mongo.Client
	log       *zap.Logger
	cfg       *config.MongoDbConnectionConfig
	connected bool
}

func NewMongoDB(logger *zap.Logger, cfg *config.MongoDbConnectionConfig) *MongoDB {
	return &MongoDB{
		Client:    nil,
		log:       logger,
		cfg:       cfg,
		connected: false,
	}
}

func (m *MongoDB) Connect() error {
	var err error
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?ssl=false&authSource=admin", m.cfg.Username, m.cfg.Password, m.cfg.Host, m.cfg.Port, m.cfg.Database)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	m.Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	var result bson.M
	if err = m.Client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		return err
	}

	m.connected = true
	return nil
}

func (m *MongoDB) IsConnected() bool {
	return m.connected
}

func (m *MongoDB) Disconnect() {
	if m.Client != nil {
		m.Client.Disconnect(context.TODO())
	}
}

func (m *MongoDB) GetDB() *mongo.Database {
	return m.Client.Database(m.cfg.Database)
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.GetDB().Collection(name)
}
