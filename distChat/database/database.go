package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	chatpb "github.com/Distributed-Messaging/distChat/chatpb"
)

type Message struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	User  string             `bson:"user,omitempty"`
	Text  string             `bson:"text,omitempty"`
	Time  int64              `bson:"time,omitempty"`
	Group string             `bson:"group,omitempty"`
}

type Group struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	IPs   []string           `bson:"IPs,omitempty"`
	Name  string             `bson:"Name,omitempty"`
	Owner string             `bson:"Owner,omitempty"`
}

type Account struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Password string             `bson:"IPs,omitempty"`
	Name     string             `bson:"Name,omitempty"`
}

func StoreAccount(account Account, collection *mongo.Collection) {
	//store a single account
	collection.InsertOne(context.Background(), account)
}

func GetOneAccountByName(name string, collection *mongo.Collection) (Account, error) {
	//returns the account with Name: name
	var a Account
	err := collection.FindOne(context.Background(), Account{Name: name}).Decode(&a)
	if err != nil {
		fmt.Println("could not read Accounts")
		return a, errors.New("could not find Account")
	}
	return a, nil
}

func StoreMessage(req *chatpb.ChatRequest, collection *mongo.Collection) {
	//constructs a database message entry from a chatrequest and stores it.
	message := req.GetMsg()
	messagetostore := Message{
		User:  message.GetUser(),
		Text:  message.GetText(),
		Time:  message.GetTime(),
		Group: message.GetGroup(),
	}

	collection.InsertOne(context.Background(), messagetostore)
}

func GetAllMessagesInGroup(name string, collection *mongo.Collection) []Message {
	//returns all messages with Group: name
	var messages []Message
	cursor, err := collection.Find(context.Background(), Message{Group: name})
	if err != nil {
		fmt.Println("could not read messages in group")
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var message Message
		cursor.Decode(&message)
		messages = append(messages, message)
	}

	return messages
}

func StoreGroup(group Group, collection *mongo.Collection) {
	//store a single group
	collection.InsertOne(context.Background(), group)
}

func GetAllGroups(collection *mongo.Collection) []Group {
	//returns all groups in the collection
	var groups []Group
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("could not read Groups")
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var group Group
		cursor.Decode(&group)
		groups = append(groups, group)
	}

	return groups
}

func GetOwnedGroups(owner string, collection *mongo.Collection) ([]Group, error) {
	//returns all groups with Owner: owner
	var groups []Group
	cursor, err := collection.Find(context.Background(), Group{Owner: owner})
	if err != nil {
		fmt.Println("could not read OwnedGroups")
		return groups, errors.New("could not read OwnedGroups")
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var group Group
		cursor.Decode(&group)
		groups = append(groups, group)
	}

	return groups, nil
}

func GetOneGroup(group string, collection *mongo.Collection) (Group, error) {
	//returns a single group struct with the corresponding name
	var g Group
	err := collection.FindOne(context.Background(), Group{Name: group}).Decode(&g)
	if err != nil {
		fmt.Println("could not read Groups")
		return g, errors.New("could not find group")
	}

	return g, nil
}
