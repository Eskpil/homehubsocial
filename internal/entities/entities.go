package entities

import (
	"context"
	"github.com/eskpil/homehubsocial/pkg/database"
	"github.com/eskpil/homehubsocial/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func CreateOne(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var body struct {
		EntityId   string `json:"entity_id"`   // entity id.
		EntityName string `json:"entity_name"` // entity name.
		UserId     string `json:"user_id"`

		State      string            `json:"state"`      // the new state of the entity
		Attributes map[string]string `json:"attributes"` // the attributes of the new state
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	opts := options.FindOne().SetProjection(bson.D{{"history", 0}})

	var entity models.Entity
	if err := database.Collection("entities").FindOne(
		ctx,
		bson.D{
			{"_id", body.EntityId},
		},
		opts,
	).Decode(&entity); err != nil {
		if err == mongo.ErrNoDocuments {
			entity.Id = body.EntityId
			entity.Name = body.EntityName
			entity.UserId = body.UserId
			entity.HistoryId = uuid.New().String()

			if _, err := database.Collection("entities").InsertOne(ctx, entity); err != nil {
				log.Errorf("failed to insert a new entity: %v", err)
				return err
			}
		} else {
			log.Errorf("failed to find the entity: %v", err)
			return err
		}
	}

	if body.EntityName != entity.Name {
		if _, err := database.Collection("entities").UpdateOne(
			ctx,
			bson.D{
				{"_id", entity.Id},
			},
			bson.D{
				{"$set", bson.D{
					{"name", body.EntityName},
				}},
			},
		); err != nil {
			log.Errorf("failed to update entity model with the new state: %v", err)

			return err
		}

		entity.Name = body.EntityName
	}

	state := models.EntityState{
		Id:        uuid.New().String(),
		HistoryId: entity.HistoryId,

		State:      body.State,
		Attributes: body.Attributes,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := database.Collection("histories").InsertOne(ctx, state); err != nil {
		log.Errorf("failed to append to the history: %v", err)
		return err
	}

	return c.JSON(http.StatusOK, entity)
}

func GetUserRanged(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	entityId := c.Param("entityId")

	opts := options.FindOne().SetProjection(bson.D{{"history", 0}})

	var entity models.Entity
	if err := database.Collection("entities").FindOne(
		ctx,
		bson.D{
			{"_id", entityId},
			{"user_id", userId},
		},
		opts,
	).Decode(&entity); err != nil {
		log.Errorf("could not find entity: %s/%s", userId, entityId)
		return err
	}

	cursor, err := database.Collection("histories").Find(
		ctx,
		bson.D{
			{"history_id", entity.HistoryId},
			{"created_at", bson.D{
				{"$gt", time.Now().Add(-time.Hour)},
			}},
		},
	)

	if err != nil {
		log.Errorf("failed to get state history of entity: %v", err)
		return err
	}

	for cursor.Next(ctx) {
		var state models.EntityState
		if err := cursor.Decode(&state); err != nil {
			log.Errorf("failed to decode document: %v", err)
			return err
		}

		entity.History = append(entity.History, state)
	}

	return c.JSON(http.StatusOK, entity)
}
