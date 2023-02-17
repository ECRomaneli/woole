package payload

func (auth *Auth) HTTPUrl() string {
	return "http://" + auth.Url + ":" + auth.HttpPort
}

func (auth *Auth) HTTPSUrl() string {
	if auth.HttpsPort == "" {
		return ""
	}
	return "https://" + auth.Url + ":" + auth.HttpsPort
}

func (auth *Auth) TunnelUrl() string {
	if auth.HttpsPort == "" {
		return "http://" + auth.Url + ":" + auth.TunnelPort
	}

	return "https://" + auth.Url + ":" + auth.TunnelPort
}
