package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/Fabese/students-app/internal/pkg/models"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type (
	Generator interface {
		ID() string
	}
	Students struct {
		db  *mongo.Collection
		svc *http.Client
		gen Generator
	}
)
type generator string

func (generator) ID() string {
	return ulid.Make().String()
}

func New(db *mongo.Collection) *Students {
	var generate generator
	return &Students{
		db:  db,
		gen: generate,
	}
}
func (s *Students) Select(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.D{{"email", email}}
	err := s.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Students) Create(ctx context.Context, in models.User) error {
	if in.ID == "" {
		in.ID = s.gen.ID()
	}
	user, err := s.Select(ctx, in.Email)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	}
	if user != nil {
		return fmt.Errorf("the email '%s' has already been registered\n", in.Email)
	}
	if _, err := s.db.InsertOne(ctx, in); err != nil {
		return err
	}
	return nil
}

func Update() {}

func Delete() {}
