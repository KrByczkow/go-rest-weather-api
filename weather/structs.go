package weather

import "time"

type DayWeather struct {
	Weather []struct {
		Timestamp                  time.Time `json:"timestamp"`
		SourceId                   int       `json:"source_id"`
		Precipitation              float64   `json:"precipitation"`
		PressureMsl                float64   `json:"pressure_msl"`
		Sunshine                   float64   `json:"sunshine"`
		Temperature                float64   `json:"temperature"`
		WindDirection              int       `json:"wind_direction"`
		WindSpeed                  float64   `json:"wind_speed"`
		CloudCover                 int       `json:"cloud_cover"`
		DewPoint                   float64   `json:"dew_point"`
		RelativeHumidity           *int      `json:"relative_humidity"`
		Visibility                 int       `json:"visibility"`
		WindGustDirection          *int      `json:"wind_gust_direction"`
		WindGustSpeed              float64   `json:"wind_gust_speed"`
		Condition                  string    `json:"condition"`
		PrecipitationProbability   *int      `json:"precipitation_probability"`
		PrecipitationProbability6H *int      `json:"precipitation_probability_6h"`
		Solar                      *float64  `json:"solar"`
		Icon                       string    `json:"icon"`
	} `json:"weather"`
}

type CurrentWeather struct {
	Weather struct {
		Timestamp           time.Time `json:"timestamp"`
		CloudCover          int       `json:"cloud_cover"`
		Condition           string    `json:"condition"`
		DewPoint            float64   `json:"dew_point"`
		Solar10             float64   `json:"solar_10"`
		Solar30             float64   `json:"solar_30"`
		Solar60             float64   `json:"solar_60"`
		Precipitation10     float64   `json:"precipitation_10"`
		Precipitation30     float64   `json:"precipitation_30"`
		Precipitation60     float64   `json:"precipitation_60"`
		PressureMsl         float64   `json:"pressure_msl"`
		RelativeHumidity    int       `json:"relative_humidity"`
		Visibility          int       `json:"visibility"`
		WindDirection10     int       `json:"wind_direction_10"`
		WindDirection30     int       `json:"wind_direction_30"`
		WindDirection60     int       `json:"wind_direction_60"`
		WindSpeed10         float64   `json:"wind_speed_10"`
		WindSpeed30         float64   `json:"wind_speed_30"`
		WindSpeed60         float64   `json:"wind_speed_60"`
		WindGustDirection10 float64   `json:"wind_gust_direction_10"`
		WindGustDirection30 float64   `json:"wind_gust_direction_30"`
		WindGustDirection60 float64   `json:"wind_gust_direction_60"`
		WindGustSpeed10     float64   `json:"wind_gust_speed_10"`
		WindGustSpeed30     float64   `json:"wind_gust_speed_30"`
		WindGustSpeed60     float64   `json:"wind_gust_speed_60"`
		Sunshine30          float64   `json:"sunshine_30"`
		Sunshine60          float64   `json:"sunshine_60"`
		Temperature         float64   `json:"temperature"`
		Icon                string    `json:"icon"`
	} `json:"weather"`
}
