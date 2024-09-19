// database/generic_Collection.go

package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GenericCollection[S any, T any] struct {
	collection *mongo.Collection
}

func NewGenericCollection[S any, T any](collection *mongo.Collection) *GenericCollection[S, T] {
	return &GenericCollection[S, T]{collection: collection}
}

// TenantFilter creates a BSON filter to match documents by tenant
func TenantFilter(tenant string) bson.D {
	return bson.D{{Key: "tenant", Value: tenant}}
}

func (r *GenericCollection[S, T]) GetAll(ctx context.Context, tenant string) ([]T, error) {
	var filter interface{}

	if tenant == "" {
		filter = bson.D{}
	} else {
		filter = TenantFilter(tenant)
	}
	// filter := TenantFilter(tenant)
	cursor, err := r.collection.Find(ctx, filter)
	log.Println(r.collection.Name())

	log.Println(cursor)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GenericCollection[S, T]) GetByID(ctx context.Context, id string, tenant string) (T, error) {
	var result T
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	// Combine tenant and ID filters
	filter := bson.M{
		"_id":    objectId,
		"tenant": tenant,
	}

	err = r.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (r *GenericCollection[S, T]) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (T, error) {
	var result T

	// . In MongoDB, when you use a dynamic key in a BSON query,
	// you need to ensure that the key is correctly interpreted as a field name.
	// Use bson.D instead of bson.M
	invalid := r.collection.FindOne(ctx, bson.D{{Key: key, Value: value},
		{Key: "tenant", Value: tenant},
	}).Decode(&result)
	log.Println(invalid, "invalid")
	return result, invalid
}

func (r *GenericCollection[S, T]) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]T, error) {
	var results []T
	log.Println(key)
	log.Println(value)
	log.Println(tenant)

	cursor, err := r.collection.Find(ctx, bson.D{{Key: key, Value: value},
		{Key: "tenant", Value: tenant},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result T
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	log.Println(results)
	return results, nil
}

func (r *GenericCollection[S, T]) Create(ctx context.Context, entity S) (string, error) {
	// Insert the entity into the collection
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		// Handle duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrDuplicateKey // Assuming ErrDuplicateKey is a pre-defined error
		}
		// Return other errors as is
		return "", err
	}
	print(result)
	// Extract the inserted ID and convert it to string
	id, ok := result.InsertedID.(primitive.ObjectID)
	print(ok)
	if !ok {
		return "", DBError
	}
	return id.Hex(), nil
}

func (r *GenericCollection[S, T]) Update(ctx context.Context, id string, entity S) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, invalid := r.collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": entity})

	return invalid
}

func (r *GenericCollection[S, T]) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, invalid := r.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	return invalid
}
