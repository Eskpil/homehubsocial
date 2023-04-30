package main

import (
	"context"
	"github.com/eskpil/homehubsocial/internal/entities"
	"github.com/eskpil/homehubsocial/internal/events"
	"github.com/eskpil/homehubsocial/internal/patients"
	"github.com/eskpil/homehubsocial/pkg/database"
	"github.com/eskpil/homehubsocial/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func seedTestPatients() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if count, err := database.Collection("patients").CountDocuments(ctx, bson.D{{}}); err != nil {
		log.Fatalf("Failed to count patients when seeding")
	} else if 0 >= count {
		// continue
	} else {
		log.Infof("There is already %v test patient(s) seeded. Ignoring", count)
		return
	}

	{
		patient := models.Patient{
			Id:   "1",
			Name: "Edith",
			Age:  82,
		}

		log.Infof("Trying to seed: %v", patient.Name)

		if _, err := database.Collection("patients").InsertOne(ctx, patient); err != nil {
			log.Errorf("Failed to seed test patient: Edith: %v", err)
		}
	}
}

func main() {
	e := echo.New()

	if err := database.Connect(); err != nil {
		log.Error("Failed to connect with the database: %v", err)
	}

	seedTestPatients()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.GET("/patients/", patients.GetMany)
	e.GET("/patients/:id/", patients.GetOne)

	e.POST("/entities/", entities.CreateOne)

	e.GET("/user/:userId/entities/:entityId/", entities.GetUserRanged)

	e.GET("/events/", events.Get)
	e.POST("/events/", events.Create)

	e.Start("0.0.0.0:8080")
}
