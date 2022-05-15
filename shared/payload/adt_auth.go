package payload

type Auth struct {
	ClientID   string
	URL        string
	HttpPort   string
	HttpsPort  string
	TunnelPort string
	Bearer     string
}

func (this *Auth) HTTPUrl() string {
	return "http://" + this.URL + ":" + this.HttpPort
}

func (this *Auth) HTTPSUrl() string {
	if this.HttpsPort == "" {
		return ""
	}
	return "https://" + this.URL + ":" + this.HttpsPort
}

func (this *Auth) TunnelUrl() string {
	if this.HttpsPort == "" {
		return "http://" + this.URL + ":" + this.TunnelPort
	}

	return "https://" + this.URL + ":" + this.TunnelPort
}
