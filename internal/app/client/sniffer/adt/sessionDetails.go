package adt

import (
	"woole/internal/app/client/app"
)

type SessionDetails struct {
	ClientID   string `json:"clientId"`
	HTTP       string `json:"http"`
	HTTPS      string `json:"https"`
	Proxying   string `json:"proxying"`
	Sniffer    string `json:"sniffer"`
	Tunnel     string `json:"tunnel"`
	MaxRecords int    `json:"maxRecords"`
	ExpireAt   string `json:"expireAt"`
}

func NewSessionDetails() *SessionDetails {
	session := app.GetSessionWhenAvailable()
	config := app.ReadConfig()

	sessionDetails := &SessionDetails{
		ClientID:   session.ClientId,
		HTTP:       session.HttpUrl(),
		HTTPS:      session.HttpsUrl(),
		Sniffer:    config.SnifferUrl.String(),
		MaxRecords: config.MaxRecords,
		ExpireAt:   app.ExpireAt(),
	}

	if !config.IsStandalone {
		sessionDetails.Tunnel = config.TunnelUrl.String()
	}

	return sessionDetails
}
