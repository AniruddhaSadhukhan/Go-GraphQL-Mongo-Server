package models

type contextKey string

const (
	PermissionDenied = "permission denied"

	// Context Keys
	UserContextKey = contextKey("User")

	// Users
	InternalUser = "__INTERNAL__"
	GuestUser    = "__GUEST__"

	//Collection Names
	UserCollection            = "users"
	SchemaMigrationCollection = "schema_migrations"
	TokenCollection           = "tokens"
)
