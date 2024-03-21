package models

import "time"

type User struct {
	ID        string    `bson:"id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password"`
	StartDate time.Time `bson:"startDate" json:"startDate"`
	Level     string    `bson:"level,omitempty" json:"level"`
}
