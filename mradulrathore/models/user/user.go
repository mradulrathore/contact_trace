package user

import (
	"context"
	"mradulrathore/contactTrace/connectionHelper"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gookit/color.v1"
)

// User is for storing user information
type User struct {
	ID          primitive.ObjectID `bson:"Id"`
	Name        string             `bson:"Name"`
	BirthDate   time.Time          `bson:"birthdate"`
	PhoneNumber string             `bson:"phoneNumber"`
	Email       string             `bson:"email"`
	CreatedOn   time.Time          `bson:"timestamp"`
}

//CreateUser
func CreateUser(user User) error {
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.USER)
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//CreateMany
func CreateMany(list []User) error {
	//Map struct slice to interface slice as InsertMany accepts interface slice as parameter
	insertableList := make([]interface{}, len(list))
	for i, v := range list {
		insertableList[i] = v
	}
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.USER)
	//Perform InsertMany operation & validate against the error.
	_, err = collection.InsertMany(context.TODO(), insertableList)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//GetUsersByID
func GetUsersByID(id string) (User, error) {
	result := User{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "Id", Value: id}}
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return result, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.USER)
	//Perform FindOne operation & validate against the error.
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	//Return result without any error.
	return result, nil
}

//GetAllUsers
func GetAllUsers() ([]User, error) {
	//Define filter query for fetching specific document from collection
	filter := bson.D{{}} //bson.D{{}} specifies 'all documents'
	users := []User{}
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return users, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.USER)
	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return users, findError
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := User{}
		err := cur.Decode(&t)
		if err != nil {
			return users, err
		}
		users = append(users, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}
	return users, nil
}

//DeleteOne
func DeleteOne(id string) error {
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "Id", Value: id}}
	//Get MongoDB connection using connectionhelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.USER)

	//Perform DeleteOne operation & validate against the error.
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//DeleteAllUsers
func DeleteAllUsers() error {
	//Define filter query for fetching specific document from collection
	selector := bson.D{{}} // bson.D{{}} specifies 'all documents'
	//Get MongoDB connection using connectionhelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.USER)
	//Perform DeleteMany operation & validate against the error.
	_, err = collection.DeleteMany(context.TODO(), selector)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//PrintList
func PrintList(users []User) {
	for i, v := range users {
		color.Yellow.Printf("%d %s %s %s %s %s\n", i+1, v.ID, v.Name, v.BirthDate, v.Email, v.PhoneNumber)
	}
}

func getDOB(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}
