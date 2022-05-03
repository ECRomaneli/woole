package payload

type Auth struct {
	Name   string `json:"name"`
	Http   string `json:"http"`
	Https  string `json:"https"`
	Bearer string `json:"bearer"`
}
