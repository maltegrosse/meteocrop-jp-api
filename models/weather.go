package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type WeatherObservation struct {
	Date      time.Time `csv:"date" json:"date"`
	DayOfYear int       `csv:"doy" json:"-"`
	Year      int       `csv:"year" json:"-"`
	Month     int       `csv:"month" json:"-"`
	Day       int       `csv:"day" json:"-"`

	AverageAirTemperature      float32 `csv:"T" json:"average_air_temperature"`
	MaximumAirTemperature      float32 `csv:"Tmax" json:"maximum_air_temperature"`
	MinimumAirTemperature      float32 `csv:"Tmin" json:"minimum_air_temperature"`
	Precipitation              float32 `csv:"Pr" json:"precipitation"`
	AtmosphericPressure        float32 `csv:"P" json:"atmospheric_pressure"`
	VaporPressure              float32 `csv:"e" json:"vapor_pressure"`
	VaporPressureDeficit       float32 `csv:"VPD" json:"vapor_pressure_deficit"`
	AverageRelativeHumidity    float32 `csv:"RH" json:"average_relative_humidity"`
	MinimumRelativeHumidity    float32 `csv:"RHmin" json:"minimum_relative_humidity"`
	AverageWindSpeed           float32 `csv:"u" json:"average_wind_speedValues"`
	MaximumWindSpeed           float32 `csv:"u10max" json:"maximum_wind_speed"`
	SunshineDuration           float32 `csv:"N" json:"sunshine_duration"`
	SolarIrradianceEstimate    float32 `csv:"Sd" json:"solar_irradiance_estimate"`
	DownwardLongwaveIrradiance float32 `csv:"Ld" json:"downward_longwave_irradiance"`
	PotentialEvaporation       float32 `csv:"EP" json:"potential_evaporation"`
	EvapotranspirationFAO      float32 `csv:"ET0" json:"evapotranspiration_FAO"`
	WaterTemperatureTw0        float32 `csv:"Tw0" json:"water_temperature_tw0"`
	WaterTemperatureTwinf      float32 `csv:"Twinf" json:"water_temperature_twinf"`
}
type HistoricalWeatherObservations struct {
	Station      *AmedasStation        `json:"station"`
	Latitude     float64               `json:"latitude"`
	Longitude    float64               `json:"longitude"`
	FromDate     time.Time             `json:"from_date"`
	ToDate       time.Time             `json:"to_date"`
	Observations []*WeatherObservation `json:"observations"`
}
// manually change header, see https://github.com/gocarina/gocsv/issues/128
type weatherObservationCsv struct {
	Date      time.Time `csv:"date" `
	DayOfYear int       `csv:"doy" `
	Year      int       `csv:"year" `
	Month     int       `csv:"month" `
	Day       int       `csv:"day"`

	AverageAirTemperature      float32 `csv:"average_air_temperature"`
	MaximumAirTemperature      float32 `csv:"maximum_air_temperature"`
	MinimumAirTemperature      float32 `csv:"minimum_air_temperature"`
	Precipitation              float32 `csv:"precipitation"`
	AtmosphericPressure        float32 `csv:"atmospheric_pressure"`
	VaporPressure              float32 `csv:"vapor_pressure"`
	VaporPressureDeficit       float32 `csv:"vapor_pressure_deficit"`
	AverageRelativeHumidity    float32 `csv:"average_relative_humidity"`
	MinimumRelativeHumidity    float32 `csv:"minimum_relative_humidity"`
	AverageWindSpeed           float32 `csv:"average_wind_speedValues"`
	MaximumWindSpeed           float32 `csv:"maximum_wind_speed"`
	SunshineDuration           float32 `csv:"sunshine_duration"`
	SolarIrradianceEstimate    float32 `csv:"solar_irradiance_estimate"`
	DownwardLongwaveIrradiance float32 `csv:"downward_longwave_irradiance"`
	PotentialEvaporation       float32 `csv:"potential_evaporation"`
	EvapotranspirationFAO      float32 `csv:"evapotranspiration_FAO"`
	WaterTemperatureTw0        float32 `csv:"water_temperature_tw0"`
	WaterTemperatureTwinf      float32 `csv:"water_temperature_twinf"`
}

func (w *HistoricalWeatherObservations) ConvertToCsv() (string, error) {

	newObservation := []*weatherObservationCsv{}
	for _, o := range w.Observations {
		newOb := weatherObservationCsv{
			Date:                       o.Date,
			DayOfYear:                  o.DayOfYear,
			AverageAirTemperature:      o.AverageAirTemperature,
			MaximumAirTemperature:      o.MaximumAirTemperature,
			MinimumAirTemperature:      o.MinimumAirTemperature,
			Precipitation:              o.Precipitation,
			AtmosphericPressure:        o.AtmosphericPressure,
			VaporPressure:              o.VaporPressure,
			VaporPressureDeficit:       o.VaporPressureDeficit,
			AverageRelativeHumidity:    o.AverageRelativeHumidity,
			MinimumRelativeHumidity:    o.MinimumRelativeHumidity,
			AverageWindSpeed:           o.AverageWindSpeed,
			MaximumWindSpeed:           o.MaximumWindSpeed,
			SunshineDuration:           o.SunshineDuration,
			SolarIrradianceEstimate:    o.SolarIrradianceEstimate,
			DownwardLongwaveIrradiance: o.DownwardLongwaveIrradiance,
			PotentialEvaporation:       o.PotentialEvaporation,
			EvapotranspirationFAO:      o.EvapotranspirationFAO,
			WaterTemperatureTw0:        o.WaterTemperatureTw0,
			WaterTemperatureTwinf:      o.WaterTemperatureTwinf,
		}
		newObservation = append(newObservation, &newOb)

	}
	
	return gocsv.MarshalString(&newObservation)
}
func (w *HistoricalWeatherObservations) ConvertToWtd() string {
	
	var sb strings.Builder
	// write header
	sb.WriteString("@  DATE  SRAD  TMAX  TMIN  RAIN  TAVE \n")
	for _, o := range w.Observations { 
		sb.WriteString(fmt.Sprint(o.Year))
		sb.WriteString(fmt.Sprintf("%03d", o.DayOfYear))
		//todo which srad?
		sb.WriteString(fmt.Sprintf("%6.1f", o.SolarIrradianceEstimate))
		sb.WriteString(fmt.Sprintf("%6.1f", o.MaximumAirTemperature))
		sb.WriteString(fmt.Sprintf("%6.1f", o.MinimumAirTemperature))
		sb.WriteString(fmt.Sprintf("%6.1f", o.Precipitation))
		sb.WriteString(fmt.Sprintf("%6.1f", o.AverageAirTemperature))
		sb.WriteString(" \n")
	}

	return sb.String()
}