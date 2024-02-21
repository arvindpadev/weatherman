package server

type OwmWeather struct {
	Main string `json:"main"`
}

type OwmMain struct {
	FeelsLike float32 `json:"feels_like"`
}

type OwmResponse struct {
	Weather []OwmWeather `json:"weather"`
	Main    OwmMain      `json:"main"`
}
