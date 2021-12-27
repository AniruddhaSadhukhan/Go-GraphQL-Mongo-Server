package models

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var dbSession *mongo.Client
var dbName string

func InitializeDB() {
	logger.Log.Info("Initializing DB")
	GetDbSession()
}

func GetDbSession() *mongo.Client {
	once.Do(func() {
		if dbSession == nil {
			dbSession = newDatabaseSession(config.ConfigManager.Database)
		}
	})
	return dbSession
}

func newDatabaseSession(db config.Database) *mongo.Client {
	logger.Log.Info("Creating new DB session")

	dbName = db.Name
	credential := options.Credential{
		Username:      db.Username,
		Password:      db.Password,
		AuthMechanism: "SCRAM-SHA-1",
	}

	mongoURI := "mongodb+srv://" + db.Host
	dbPort := strings.Trim(db.Port, " ")
	if len(dbPort) != 0 {
		mongoURI = mongoURI + ":" + dbPort
	}

	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(credential)

	if db.InsecureSkipVerify {
		clientOptions.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	session, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Log.Error("Error connecting to DB: " + err.Error())
	}

	return session
}

func getCollection(collectionName string) *mongo.Collection {
	return GetDbSession().Database(dbName).Collection(collectionName)
}

func Insert(collectionName string, document interface{}, ctx context.Context) error {

	_, err := getCollection(collectionName).InsertOne(ctx, document)
	if err != nil {
		logger.Log.Error("Error inserting document: " + err.Error())
	}
	return err

}

func InsertMany(collectionName string, documents []interface{}, ctx context.Context) error {

	_, err := getCollection(collectionName).InsertMany(ctx, documents)
	if err != nil {
		logger.Log.Error("Error inserting documents: " + err.Error())
	}
	return err

}

func IsExist(collectionName string, filter interface{}, ctx context.Context) bool {

	count, err := getCollection(collectionName).CountDocuments(ctx, filter)
	if err != nil {
		logger.Log.Error("Error checking document existence: " + err.Error())
	}
	return count > 0

}

// Generic Find One Document from MongoDB
func FindOne(collectionName string, filter interface{}, projection interface{}, resultPointer interface{}, ctx context.Context) error {

	err := getCollection(collectionName).FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(resultPointer)
	if err != nil {
		logger.Log.Error("Error finding document: " + err.Error())
	}
	return err

}

// Generic Find All Documents from MongoDB
func FindAll(collectionName string, filter interface{}, projection interface{}, resultSlicePointer interface{}, ctx context.Context) error {

	cursor, err := getCollection(collectionName).Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		logger.Log.Error("Error finding documents: " + err.Error())
		return err
	}

	err = cursor.All(ctx, resultSlicePointer)
	if err != nil {
		logger.Log.Error("Error decoding documents: " + err.Error())
	}
	return err

}

// Generic Aggregate Documents from MongoDB
func Aggregate(collectionName string, pipeline []bson.M, resultSlicePointer interface{}, ctx context.Context) error {

	cursor, err := getCollection(collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		logger.Log.Error("Error aggregating documents: " + err.Error())
		return err
	}

	err = cursor.All(ctx, resultSlicePointer)
	if err != nil {
		logger.Log.Error("Error decoding documents: " + err.Error())
	}
	return err

}

//Update with options
func UpdateWithOptions(collectionName string, filter interface{}, update interface{}, options *options.UpdateOptions, ctx context.Context) error {

	res, err := getCollection(collectionName).UpdateOne(ctx, filter, update, options)
	if err != nil {
		logger.Log.Error("Error updating document: " + err.Error())
	}

	// Check if anything got modified or upserted
	if res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("no document modified")
	}

	return err

}

func Update(collectionName string, filter interface{}, update interface{}, ctx context.Context) error {

	return UpdateWithOptions(collectionName, filter, update, nil, ctx)

}

func Upsert(collectionName string, filter interface{}, update interface{}, ctx context.Context) error {

	return UpdateWithOptions(collectionName, filter, update, options.Update().SetUpsert(true), ctx)

}

func UpdateMany(collectionName string, filter interface{}, update interface{}, ctx context.Context) error {

	_, err := getCollection(collectionName).UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Log.Error("Error updating documents: " + err.Error())
	}
	return err

}

func Delete(collectionName string, filter interface{}, ctx context.Context) error {

	_, err := getCollection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		logger.Log.Error("Error deleting document: " + err.Error())
	}
	return err

}

func DeleteMany(collectionName string, filter interface{}, ctx context.Context) error {

	_, err := getCollection(collectionName).DeleteMany(ctx, filter)
	if err != nil {
		logger.Log.Error("Error deleting documents: " + err.Error())
	}
	return err

}

func FindDistinct(collectionName string, fieldName string, filter interface{}, ctx context.Context) ([]string, error) {

	if filter == nil {
		filter = bson.M{}
	}

	queryResults, err := getCollection(collectionName).Distinct(ctx, fieldName, filter)
	if err != nil {
		logger.Log.Error("Error finding distinct: " + err.Error())
		return nil, err
	}

	result := make([]string, len(queryResults))
	for i, v := range queryResults {
		result[i] = fmt.Sprint(v)
	}

	return result, nil

}
