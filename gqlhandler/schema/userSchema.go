package schema

import "github.com/graphql-go/graphql"

var UserSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"dob": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"address": &graphql.Field{
				Type: graphql.NewNonNull(UserAddressSchema),
			},
			"isVerified": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"remarks": &graphql.Field{
				Type: graphql.String,
			},
			"subscription": &graphql.Field{
				Type: SubscriptionTypeEnum,
			},
		},
	},
)

var UserAddressSchema = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Address",
		Fields: graphql.Fields{
			"block": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"street": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"city": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)

var SubscriptionTypeEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "SubscriptionType",
	Values: graphql.EnumValueConfigMap{
		"FREE": &graphql.EnumValueConfig{
			Value: "Free",
		},
		"PAID": &graphql.EnumValueConfig{
			Value: "Paid",
		},
	},
})

var UserSchemaInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "UserInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"dob": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"address": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(UserAddressSchemaInput),
			},
			"isVerified": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
			"remarks": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"subscription": &graphql.InputObjectFieldConfig{
				Type: SubscriptionTypeEnum,
				DefaultValue: "Free",
			},
		},
	},
)

var UserAddressSchemaInput = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "AddressInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"block": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"street": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"city": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	},
)
