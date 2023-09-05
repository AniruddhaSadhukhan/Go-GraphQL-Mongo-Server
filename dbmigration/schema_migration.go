package dbmigration

import (
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"go-graphql-mongo-server/models"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file driver
)

// RunDbSchemaMigration runs the database schema migration
// Give the relative path prefix to the main package as input
// Eg: "." or "./" or "", "./.." or "./../" or "../" or "..", "./../.." etc.
func RunDbSchemaMigration(relativePathPrefixToMainDir string) error {

	dbName := config.Store.Database.Name
	mongoDriver, err := mongodb.WithInstance(models.GetDbSession(), &mongodb.Config{
		DatabaseName:         dbName,
		MigrationsCollection: models.SchemaMigrationCollection,
		TransactionMode:      false,
	})
	if err != nil {
		logger.Log.Error("Error creating mongo driver: " + err.Error())
		return err
	}

	if len(relativePathPrefixToMainDir) > 0 && !strings.HasSuffix(relativePathPrefixToMainDir, "/") {
		relativePathPrefixToMainDir = relativePathPrefixToMainDir + "/"
	}
	migrationPath := "file://" + relativePathPrefixToMainDir + "resources/schema_migrations"

	// Read migrations from resources/schema_migrations and connect to database
	migrateInstance, err := migrate.NewWithDatabaseInstance(
		migrationPath,
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
