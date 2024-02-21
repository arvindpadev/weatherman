package server

type FeelsLike string

const (
	Hot      FeelsLike = "Hot"
	Cold     FeelsLike = "Cold"
	Moderate FeelsLike = "Moderate"
)

type WeatherCondition struct {
	Feeling   FeelsLike `json:"feeling"`
	Condition string    `json:"condition"`
}
