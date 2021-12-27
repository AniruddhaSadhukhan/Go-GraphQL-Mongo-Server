package dbmigration

import (
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file driver
)

func RunDbSchemaMigration() error {

	dbName := config.ConfigManager.Database.Name
	mongoDriver, err := mongodb.WithInstance(models.GetDbSession(), &mongodb.Config{
		DatabaseName:         dbName,
		MigrationsCollection: models.SchemaMigrationCollection,
		TransactionMode:      false,
	})
	if err != nil {
		logger.Log.Error("Error creating mongo driver: " + err.Error())
		return err
	}

	// Read migrations from resources/schema_migrations and connect to database
	migrateInstance, err := migrate.NewWithDatabaseInstance(
		"file://resources/schema_migrations",
		dbName,
		mongoDriver,
	)
	if err != nil {
		logger.Log.Error("Error reading migration files: " + err.Error())
		return err
	}

	// Run migrations all the way up
	err = migrateInstance.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Log.Error("Error running migrations: " + err.Error())
		return err
	}

	logger.Log.Info("Successfully ran schema migration")

	return nil

}
