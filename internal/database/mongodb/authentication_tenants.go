package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type TenantUser struct {
	Email    string `json:"email"`
	Password string `json:"password"` // Ensure passwords are hashed and never stored as plain text
	Tenant   string `json:"tenant"`
}

type TenantUserResponse struct {
	ID       string `bson:"_id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"` // Ensure passwords are hashed and never stored as plain text
	Tenant   string `json:"tenant"`
}

type LoginCollection struct {
	*GenericCollection[TenantUser, TenantUserResponse]
}

func NewLoginCollection(collection *mongo.Collection) *LoginCollection {
	return &LoginCollection{
		GenericCollection: NewGenericCollection[TenantUser, TenantUserResponse](collection),
	}
}

var _ Collection[TenantUser, TenantUserResponse] = (*LoginCollection)(nil)

func (r *LoginCollection) Create(ctx context.Context, entity TenantUser) (string, error) {
	tenantUser := TenantUser{
		Email:    entity.Email,
		Password: entity.Password,
		Tenant:   entity.Tenant,
	}
	log.Println(tenantUser)
	return r.GenericCollection.Create(ctx, tenantUser)
}

func (r *LoginCollection) GetByKeyValue(ctx context.Context, key string, value string, tenant string) (TenantUserResponse, error) {
	return r.GenericCollection.GetByKeyValue(ctx, key, value, tenant)
}

func (r *LoginCollection) Delete(ctx context.Context, id string) error {
	return r.GenericCollection.Delete(ctx, id)
}
