package storage

import (
	"context"

	"github.com/ssonit/oms-orders/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DbName   = "orders"
	CollName = "orders"
)

type store struct {
	db *mongo.Client
}

func NewStore(db *mongo.Client) *store {
	return &store{
		db: db,
	}
}

func (s *store) Create(ctx context.Context, o model.OrdersCreation) (primitive.ObjectID, error) {
	col := s.db.Database(DbName).Collection(CollName)

	newOrder, err := col.InsertOne(ctx, o)

	id := newOrder.InsertedID.(primitive.ObjectID)
	return id, err
}

func (s *store) GetItem(ctx context.Context, filter map[string]interface{}) (*model.OrderItem, error) {
	col := s.db.Database(DbName).Collection(CollName)

	_id, _ := primitive.ObjectIDFromHex(filter["_id"].(string))

	var order model.OrderItem
	err := col.FindOne(ctx, bson.M{
		"_id":        _id,
		"customerId": filter["customerId"],
	}).Decode(&order)

	return &order, err
}
