package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type FeedEmail struct {
	Tenant           string    `json:"tenant"`
	MasterEmail      string    `json:"from_email_address"`
	ReceivedDateTime time.Time `json:"received_date_time"`
	ToEmailAddress   string    `json:"to_email_address"`
	Subject          string    `json:"subject"`
	BodyContent      string    `json:"body_content"`
	ShipmentId       string    `json:"shipment_id"`
}

type FeedEmailResponse struct {
	ID               string    `bson:"_id"`
	Tenant           string    `json:"tenant"`
	MasterEmail      string    `json:"master_email"`
	ReceivedDateTime time.Time `json:"received_date_time"`
	ToEmailAddress   string    `json:"to_email_address"`
	Subject          string    `json:"subject"`
	BodyContent      string    `json:"body_content"`
	ShipmentId       string    `json:"shipment_id"`
}

type FeedEmailCollection struct {
	*GenericCollection[FeedEmail, FeedEmailResponse]
}

func NewFeedMessageCollection(collection *mongo.Collection) *FeedEmailCollection {
	return &FeedEmailCollection{
		GenericCollection: NewGenericCollection[FeedEmail, FeedEmailResponse](collection),
	}
}

var _ Collection[FeedEmail, FeedEmailResponse] = (*FeedEmailCollection)(nil)

func (r *FeedEmailCollection) Create(ctx context.Context, entity FeedEmail) (string, error) {
	feedEmail := FeedEmail{
		Tenant:           entity.Tenant,
		MasterEmail:      entity.MasterEmail,
		ReceivedDateTime: entity.ReceivedDateTime,
		ToEmailAddress:   entity.ToEmailAddress,
		Subject:          entity.Subject,
		BodyContent:      entity.BodyContent,
		ShipmentId:       entity.ShipmentId,
	}
	return r.GenericCollection.Create(ctx, feedEmail)
}

func (r *FeedEmailCollection) GetAll(ctx context.Context, tenant string) ([]FeedEmailResponse, error) {
	return r.GenericCollection.GetAll(ctx, tenant)
}

func (r *FeedEmailCollection) GetAllByKeyValue(ctx context.Context, key string, value string, tenant string) ([]FeedEmailResponse, error) {
	return r.GenericCollection.GetAllByKeyValue(ctx, key, value, tenant)
}
