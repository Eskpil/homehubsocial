package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var client *mongo.Client = nil
var db *mongo.Database = nil

func Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	host := "localhost"
	if os.Getenv("ENV") == "docker" {
		host = "host.docker.internal"
	}

	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Infof("%v", evt.Command)
		},
	}

	_ = cmdMonitor

	clientOps := options.Client().SetHosts(
		[]string{host},
	).SetAuth(options.Credential{
		Username: "user",
		Password: "secret",
	})
	//.SetMonitor(cmdMonitor)

	client, err = mongo.NewClient(clientOps)
	if err != nil {
		return fmt.Errorf("failed to create mongo client: %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect with mongo: %v", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping primary in the database: %v", err)
	}

	db = client.Database("homehubsocial")
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	return err
}

func Collection(collection string) *mongo.Collection {
	return db.Collection(collection)
}
