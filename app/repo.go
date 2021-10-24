package app

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoRepository struct {
	client *mongo.Client
}

func (r *MongoRepository) init(){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	r.client = client
	fmt.Println("Connected to MongoDB!")
}

func (r *MongoRepository) getClient() *mongo.Client {
	if r.client == nil {
		r.init()
	}
	return r.client
}

func (r *MongoRepository) InsertOrReplaceVereniging(v *Vereniging){
	if v.Id == "" {
		log.Fatal("Vereniging zonder ID kan niet worden opgeslagen")
	}
	filter := bson.D{{
		"id", v.Id,
	}}

	collection := r.getClient().Database("nbg").Collection("verenigingen")
	upsert := true
	after := options.After
	opts := options.FindOneAndReplaceOptions{
		Upsert: &upsert,
		ReturnDocument: &after,
	}

	result := new(Vereniging)
	collection.FindOneAndReplace(context.TODO(), filter, v, &opts).Decode(result)
}

func (r *MongoRepository) InsertOrReplaceGebeurtenis(g *Gebeurtenis) {
	if g.Id == "" {
		log.Fatal("Gebeurtenis zonder ID kan niet worden opgeslagen")
	}

	filter := bson.D{{
		"id", g.Id,
	}}

	collection := r.getClient().Database("nbg").Collection("gebeurtenissen")

	upsert := true
	after := options.After
	opts := options.FindOneAndReplaceOptions{
		Upsert: &upsert,
		ReturnDocument: &after,
	}

	result := new(Gebeurtenis)
	collection.FindOneAndReplace(context.TODO(), filter, g, &opts).Decode(result)
}

func (r *MongoRepository) GetVerenigingById(id string) *Vereniging{
	collection := r.getClient().Database("nbg").Collection("verenigingen")

	filter := bson.D{{
		"id", id,
	}}

	result := new(Vereniging)
	err := collection.FindOne(context.TODO(), filter).Decode(result)
	if err != nil {
		return nil
	}
	return result
}


