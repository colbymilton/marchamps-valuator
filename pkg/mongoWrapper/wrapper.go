package mongoWrapper

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BsonNoneD = bson.D{}
var BsonNoneM = bson.M{}

type MongoDB struct {
	db *mongo.Database
}

func NewMongoDB(uri, database string) *MongoDB {
	mdb := &MongoDB{}

	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}
	mdb.db = cli.Database(database)

	return mdb
}

// CreateMany with insert multiple documents into the database
// if a document with a matching "_id" already exists, it is ignored
func CreateMany[T any](mdb *MongoDB, coll string, things []*T) error {
	collection := mdb.db.Collection(coll)

	for _, thing := range things {
		_, err := collection.InsertOne(context.Background(), thing)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return nil
}

// GetMany returns a slice of T objects that
// are retrieved from the specified database, collection, sort, and filter
func GetMany[T any](mdb *MongoDB, coll string, filter bson.D, sort bson.M) ([]*T, error) {
	things := make([]*T, 0)
	collection := mdb.db.Collection(coll)
	opts := options.Find().SetSort(sort)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var next *T
		err := cur.Decode(&next)
		if err != nil {
			return nil, err
		}
		things = append(things, next)
	}

	return things, nil
}

// GetAll returns a slice of T objects that
// are retrieved from the specified database and collection without filters or sorting
func GetAll[T any](mdb *MongoDB, coll string) ([]*T, error) {
	return GetMany[T](mdb, coll, BsonNoneD, BsonNoneM)
}

// GetOne returns a limit of 1 thing
// returns nil, nil if there isn't one thing to return
func GetOne[T any](mdb *MongoDB, coll string, filter bson.D, sort bson.M) (*T, error) {
	collection := mdb.db.Collection(coll)
	opts := options.Find().SetSort(sort).SetLimit(1)
	cur, err := collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var thing *T
		err := cur.Decode(&thing)
		if err != nil {
			return nil, err
		}
		return thing, nil
	}

	return nil, nil
}

func EmptyCollection(mdb *MongoDB, coll string) {
	mdb.db.Collection(coll).Drop(context.Background())
}

func GetCollectionSize(mdb *MongoDB, coll string) int {
	result := mdb.db.RunCommand(context.Background(), bson.M{"collStats": coll})
	var document bson.M
	if err := result.Decode(&document); err != nil {
		return -1
	}
	if i, ok := document["count"].(int32); !ok {
		return -2
	} else {
		return int(i)
	}
}

// BuildEqualsFilter returns the bson.D representation of simple equals filter
// with given key/value
func BuildEqualsFilter(key string, val interface{}) bson.D {
	return bson.D{{Key: key, Value: val}}
}

// BuildAndFilter takes in a set of existing filters and returns a filter
// that "ands" them all together
func BuildAndFilter(parts []bson.D) bson.D {
	a := bson.A{}
	for _, part := range parts {
		a = append(a, part)
	}
	return bson.D{{Key: "$and", Value: a}}
}
