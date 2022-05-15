package dashboard

import (
	"woole/cmd/client/app"
)

type Info struct {
	ClientID   string `json:"clientId"`
	HTTP       string `json:"http"`
	HTTPS      string `json:"https"`
	Proxying   string `json:"proxying"`
	Dashboard  string `json:"dashboard"`
	Tunnel     string `json:"tunnel"`
	MaxRecords int    `json:"maxRecords"`
	Bearer     string `json:"bearer"`
}

func (this *Info) FromConfig() *Info {
	config := app.ReadConfig()
	auth := app.GetAuth()

	this.ClientID = auth.ClientID
	this.HTTP = auth.HTTPUrl()
	this.HTTPS = auth.HTTPSUrl()
	this.Proxying = config.CustomHost
	this.Dashboard = config.DashboardURL()
	this.Tunnel = auth.TunnelUrl()
	this.MaxRecords = config.MaxRecords
	this.Bearer = auth.Bearer

	return this
}
