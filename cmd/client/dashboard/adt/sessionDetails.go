package adt

import (
	"woole/cmd/client/app"
)

type SessionDetails struct {
	ClientID   string `json:"clientId"`
	HTTP       string `json:"http"`
	HTTPS      string `json:"https"`
	Proxying   string `json:"proxying"`
	Dashboard  string `json:"dashboard"`
	Tunnel     string `json:"tunnel"`
	MaxRecords int    `json:"maxRecords"`
}

func NewSessionDetails() *SessionDetails {
	session := app.GetSessionWhenAvailable()
	config := app.ReadConfig()

	sessionDetails := &SessionDetails{
		ClientID:   session.ClientId,
		HTTP:       session.HttpUrl(),
		HTTPS:      session.HttpsUrl(),
		Proxying:   config.CustomUrl.String(),
		Dashboard:  config.DashboardUrl.String(),
		MaxRecords: config.MaxRecords,
	}

	if !config.IsStandalone {
		sessionDetails.Tunnel = config.TunnelUrl.String()
	}

	return sessionDetails
}
