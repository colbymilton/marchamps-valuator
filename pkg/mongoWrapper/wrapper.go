package mongoWrapper

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

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

	cli, err := mongo.Connect(defaultContext(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}
	mdb.db = cli.Database(database)

	return mdb
}

func defaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	return ctx
}

func findBsonField(thing interface{}, fieldName string) (interface{}, error) {
	rv := reflect.ValueOf(thing)
	rt := reflect.TypeOf(thing)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		rt = rt.Elem()
	}
	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			fv := rv.Field(i)
			ft := rt.Field(i)

			if strings.Contains(string(ft.Tag), "bson:\""+fieldName+"\"") {
				return fv.Interface(), nil
			}
		}
	}

	return nil, fmt.Errorf("could not find the %v field", fieldName)
}

func (mdb *MongoDB) Ping() error {
	return mdb.db.Client().Ping(defaultContext(), nil)
}

func (mdb *MongoDB) EmptyCollection(coll string) {
	mdb.db.Collection(coll).Drop(defaultContext())
}

func (mdb *MongoDB) GetCollectionSize(coll string) int {
	result := mdb.db.RunCommand(defaultContext(), bson.M{"collStats": coll})
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

// CreateMany will insert multiple documents into the database
// if a document with a matching "_id" already exists, it is ignored
func CreateMany[T any](mdb *MongoDB, coll string, things []*T) error {
	collection := mdb.db.Collection(coll)

	for _, thing := range things {
		_, err := collection.InsertOne(defaultContext(), thing)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return nil
}

// ReplaceMany
func ReplaceMany[T any](mdb *MongoDB, coll string, filter bson.D, things []*T) error {
	for _, thing := range things {
		if err := ReplaceOne[T](mdb, coll, filter, thing); err != nil {
			return err
		}
	}
	return nil
}

// ReplaceManyID
func ReplaceManyID[T any](mdb *MongoDB, coll string, things []*T) error {
	for _, thing := range things {
		if err := ReplaceOneID[T](mdb, coll, thing); err != nil {
			return err
		}
	}
	return nil
}

// ReplaceOne will replace a document in the database with the given filter
func ReplaceOne[T any](mdb *MongoDB, coll string, filter bson.D, thing *T) error {
	collection := mdb.db.Collection(coll)

	opts := options.Replace().SetUpsert(true)
	_, err := collection.ReplaceOne(defaultContext(), filter, thing, opts)
	return err
}

// ReplaceOne will replace a document in the database that has a matching "_id"
func ReplaceOneID[T any](mdb *MongoDB, coll string, thing *T) error {
	id, err := findBsonField(thing, "_id")
	if err != nil {
		return err
	}

	filter := BuildEqualsFilter("_id", id)

	return ReplaceOne[T](mdb, coll, filter, thing)
}

// GetMany returns a slice of T objects that
// are retrieved from the specified database, collection, sort, and filter
func GetMany[T any](mdb *MongoDB, coll string, filter bson.D, sort bson.M) ([]*T, error) {
	things := make([]*T, 0)
	collection := mdb.db.Collection(coll)
	opts := options.Find().SetSort(sort)
	cur, err := collection.Find(defaultContext(), filter, opts)
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
	cur, err := collection.Find(defaultContext(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(defaultContext())
	for cur.Next(defaultContext()) {
		var thing *T
		err := cur.Decode(&thing)
		if err != nil {
			return nil, err
		}
		return thing, nil
	}

	return nil, nil
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
