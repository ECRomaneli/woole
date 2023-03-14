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

func (session *SessionDetails) FromConfig(config *app.Config) *SessionDetails {
	auth := app.GetSession()

	session.ClientID = auth.ClientId
	session.HTTP = auth.HttpUrl()
	session.HTTPS = auth.HttpsUrl()
	session.Proxying = config.CustomUrl.String()
	session.Dashboard = config.DashboardUrl.String()
	session.Tunnel = config.TunnelUrl.String()
	session.MaxRecords = config.MaxRecords

	return session
}
