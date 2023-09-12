package telemetry

import (
	"go-graphql-mongo-server/common"
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/shirou/gopsutil/v3/host"
)

type Event struct {
	EventType string    `json:"eventType"`
	TimeStamp time.Time `json:"timeStamp"`
	Platform  Platform  `json:"platform"`
	User      User      `json:"user"`
	Device    Device    `json:"device"`
	EventData EventData `json:"eventData"`
}

type Platform struct {
	PlatformName  string `json:"platformName"`
	Env           string `json:"env"`
	ComponentName string `json:"componentName"`
}

type User struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Device struct {
	DeviceType string `json:"deviceType"`
	OS         string `json:"os"`
	Device     string `json:"device"`
	OsVersion  string `json:"osVersion"`
}

type EventData struct {
	Operation string         `json:"operation"`
	Name      string         `json:"name"`
	Args      map[string]any `json:"args"`
	IsFailed  bool           `json:"isFailed,omitempty"`
	Error     string         `json:"error,omitempty"`
}

func LogGraphQlCall(params graphql.ResolveParams, graphQlError error) {

	operation := params.Info.Operation.GetOperation()
	name := params.Info.FieldName
	username := common.GetUserName(params)

	logger.Log.Infof("[GraphQl] %v '%v' called by %v", operation, name, username)
	if graphQlError != nil {
		logger.Log.Errorf("[GraphQl] %v '%v' error: %v", operation, name, graphQlError)
	}

	if !config.Store.ProductionMode || config.Store.TelemetryURL == "" {
		return
	}

	deviceInfo, _ := host.Info()

	telemetryEvent := Event{
		EventType: "BACKEND_GRAPHQL_CALL",
		TimeStamp: time.Now(),
		Platform: Platform{
			PlatformName:  config.Store.PlatformName,
			Env:           config.Store.Env,
			ComponentName: config.Store.ComponentName,
		},
		User: User{
			ID:   username,
			Type: common.GetUserType(params),
		},
		Device: Device{
			DeviceType: "Server",
			OS:         deviceInfo.OS,
			Device:     deviceInfo.Hostname,
			OsVersion:  deviceInfo.Platform + " - " + deviceInfo.PlatformVersion,
		},
		EventData: EventData{
			Operation: operation,
			Name:      name,
			Args:      params.Args,
		},
	}

	if graphQlError != nil {
		telemetryEvent.EventData.IsFailed = true
		telemetryEvent.EventData.Error = graphQlError.Error()
	}

	sendTelemetry(telemetryEvent)
}

func sendTelemetry(telemetryEvent Event) {

	client := common.NewHTTPClient(config.Store.TelemetryURL, config.Store.Auth.SecretToken, config.Store.InsecureSkipVerify)

	query := `
		mutation AddTelemetry($events: [TelemetryEventInput]!) {
			addTelemetryEvents(events: $events)
		}
	`

	variables := map[string]any{
		"events": []Event{telemetryEvent},
	}

	var response struct {
		AddTelemetryEvents bool `json:"addTelemetryEvents"`
	}

	err := client.ExecuteGraphQl(query, variables, &response)
	if err != nil {
		logger.Log.Errorf("error while sending telemetry: %v", err)
	}
}
