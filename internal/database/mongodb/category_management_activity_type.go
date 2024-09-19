package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryManagementActivityType struct {
	ActivityType string `json:"activity_type"`
	Tenant       string `json:"tenant"`
}

type CategoryManagementActivityTypeResponse struct {
	ID           string `bson:"_id"`
	ActivityType string `json:"activity_type"`
	Tenant       string `json:"tenant"`
}

type ActivityTypeCollection struct {
	*GenericCollection[CategoryManagementActivityType, CategoryManagementActivityTypeResponse]
}

func NewActivityTypeCollection(collection *mongo.Collection) *ActivityTypeCollection {
	return &ActivityTypeCollection{
		GenericCollection: NewGenericCollection[CategoryManagementActivityType, CategoryManagementActivityTypeResponse](collection),
	}
}

var _ Collection[CategoryManagementActivityType, CategoryManagementActivityTypeResponse] = (*ActivityTypeCollection)(nil)

func (r *ActivityTypeCollection) Create(ctx context.Context, entity CategoryManagementActivityType) (string, error) {
	activityType := CategoryManagementActivityType{
		ActivityType: entity.ActivityType,
		Tenant:       entity.Tenant,
	}
	return r.GenericCollection.Create(ctx, activityType)
}

// func (r *ActivityTypeCollection) Update(ctx context.Context, id string, entity ActivityType) error {

// 	activityType := ActivityType{
// 		ActivityType: entity.ActivityType,
// 		Tenant: entity.Tenant,
// 	}
// 	return r.GenericCollection.Update(ctx, id, activityType)
// }

func (r *ActivityTypeCollection) GetAll(ctx context.Context, tenant string) ([]CategoryManagementActivityTypeResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *ActivityTypeCollection) GetByID(ctx context.Context, id string, tenant string) (CategoryManagementActivityTypeResponse, error) {
	return r.GenericCollection.GetByID(ctx, id, tenant)
}

func (r *ActivityTypeCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
