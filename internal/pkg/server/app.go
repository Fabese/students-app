package server

import (
	"context"
	"errors"
	"github.com/Fabese/students-app/internal/pkg/models"
	"github.com/Fabese/students-app/internal/pkg/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func Serve() {
	const (
		timeout = 120 * time.Minute
	)

	var (
		opt         = options.Client()
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		opts        = []*options.ClientOptions{
			opt.ApplyURI("mongodb://localhost:27017"),
		}
	)
	defer cancel()
	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	coll := client.Database("course").Collection("students")
	students := repositories.New(coll)
	router.GET("/:email", func(c *gin.Context) {
		email := c.Param("email")
		user, err := students.Select(ctx, email)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				c.JSON(http.StatusNotFound, gin.H{"error": "the user does not exist"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, user)
	})
	router.POST("/create", func(c *gin.Context) {
		var newUser models.User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al decodificar JSON"})
			return
		}
		if err := students.Create(ctx, newUser); err != nil {
			return
		}

		c.JSON(http.StatusCreated, newUser)
	})
	router.Run(":8080")
}
