package model

type ClientConfig struct {
	Telemetry                bool   `json:"telemetry"`
	TelemetryID              string `json:"telemetryid"`
	EnablePublicSharedBoards bool   `json:"enablePublicSharedBoards"`
}
