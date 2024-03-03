package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NYTimesNews struct {
	NewsID          string `bson:"newsId,omitempty"`
	Category        string `bson:"category,omitempty"`
	OriginalNewsURL string `bson:"originalNewsURL,omitempty"`
	CreatedNewsURL  string `bson:"createdNewsURL,omitempty"`
}

func (n *NYTimesNews) Save(ctx context.Context, collection *mongo.Collection) error {
	_, err := collection.InsertOne(ctx, n)
	return err
}

func FindByID(ctx context.Context, collection *mongo.Collection, id string) (*NYTimesNews, error) {
	var nyTimesNews NYTimesNews
	err := collection.FindOne(ctx, bson.M{"newsId": id}).Decode(&nyTimesNews)
	if err != nil {
		return nil, err
	}
	return &nyTimesNews, nil
}

func FindAll(ctx context.Context, collection *mongo.Collection) ([]*NYTimesNews, error) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var nyTimesNews []*NYTimesNews
	for cursor.Next(ctx) {
		var news NYTimesNews
		err := cursor.Decode(&news)
		if err != nil {
			return nil, err
		}
		nyTimesNews = append(nyTimesNews, &news)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return nyTimesNews, nil
}
