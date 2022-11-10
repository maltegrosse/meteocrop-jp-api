package main

import (
	"net/http"
	"os"
	"github.com/maltegrosse/meteocrop-jp-api/log"
	"github.com/maltegrosse/meteocrop-jp-api/models"
	"github.com/maltegrosse/meteocrop-jp-api/routers"
)



func main() {
	log.Logger.Debug("loading AMEDAS stations")
	amedasUrl := getEnv("AMEDAS_STATION_URL", "https://meteocrop.dc.affrc.go.jp/csv/amedas.csv")
	amedasWeatherUrl := getEnv("AMEDAS_WEATHER_URL", "https://meteocrop.dc.affrc.go.jp/real/download.php?kind=5&id=")
	models.Stations = models.AmedasStations{}
	err := models.Stations.GetStations(amedasUrl, amedasWeatherUrl)
	if err != nil {
		panic(err)
	}

	apiPort := getEnv("API_PORT", "8080")
	log.Logger.Debug("starting MeteoCrop jp api on port " + apiPort)
	apiConnection := ":" + apiPort
	err = http.ListenAndServe(apiConnection, routers.Routes())
	if err != nil {
		panic(err)
	}

}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
