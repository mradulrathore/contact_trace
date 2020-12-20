package contacts

import (
	"context"
	"mradulrathore/contactTrace/connectionHelper"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Contact is for storing contact details
type Contact struct {
	UserIDOne string    `bson:"userIdOne"`
	UserIDTwo string    `bson:"userIdTwo"`
	Timestamp time.Time `bson:"timestamp"`
}

//CreateContact for creating contacts
func CreateContact(contact Contact) error {
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.CONTACT)
	//Perform InsertOne operation & validate against the error.
	_, err = collection.InsertOne(context.TODO(), contact)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//CreateManyContact - Insert multiple contacts at once in the collection.
func CreateManyContact(list []Contact) error {
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
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.CONTACT)
	//Perform InsertMany operation & validate against the error.
	_, err = collection.InsertMany(context.TODO(), insertableList)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}

//GetContacts - Get  contact of given user id  for collection
func GetContacts(id string) (Contact, error) {
	result := Contact{}
	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "userIdOne", Value: id}}
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return result, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.CONTACT)
	//Perform FindOne operation & validate against the error.
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	//Return result without any error.

	return result, nil
}

//GetAllContacts - Get All contacts of given user id
func GetAllContacts(id string) ([]Contact, error) {
	filter := bson.D{primitive.E{Key: "userIdOne", Value: id}}
	contacts := []Contact{}
	//Get MongoDB connection using connectionHelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return contacts, err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.CONTACT)
	//Perform Find operation & validate against the error.
	cur, findError := collection.Find(context.TODO(), filter)
	if findError != nil {
		return contacts, findError
	}
	//Map result to slice
	for cur.Next(context.TODO()) {
		t := Contact{}
		err := cur.Decode(&t)
		if err != nil {
			return contacts, err
		}
		contacts = append(contacts, t)
	}
	// once exhausted, close the cursor
	cur.Close(context.TODO())
	if len(contacts) == 0 {
		return contacts, mongo.ErrNoDocuments
	}
	return contacts, nil
}

//DeleteAllContacts ...
func DeleteAllContacts() error {
	//Define filter query for fetching specific document from collection
	selector := bson.D{{}} // bson.D{{}} specifies 'all documents'
	//Get MongoDB connection using connectionhelper.
	client, err := connectionHelper.GetMongoClient()
	if err != nil {
		return err
	}
	//Create a handle to the respective collection in the database.
	collection := client.Database(connectionHelper.DB).Collection(connectionHelper.CONTACT)
	//Perform DeleteMany operation & validate against the error.
	_, err = collection.DeleteMany(context.TODO(), selector)
	if err != nil {
		return err
	}
	//Return success without any error.
	return nil
}
