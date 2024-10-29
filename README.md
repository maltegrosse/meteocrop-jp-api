# MeteoCropDB JP API Service
## About MeteoCropDB Japan
This database stores daily weather data of 1980-2008 (from 1976 onward in some areas) at the AMeDAS point (about 850 sites) nationwide. In addition to basic elements such as temperature, wind speed and precipitation, it is a major feature of this database that it houses elements that have a significant influence on crop production such as solar radiation, humidity and transpiration demand. Although solar radiation and humidity are not observed at the AMeDAS point, it is presumed by utilizing the data of the sunshine hours and the meteorological offices in the neighborhood. By considering the difference in characteristics depending on the type of sunshine tower etc., stable solar radiation estimation value can be provided for a long time. As for the transpiration demand, it is calculated using the heat balance model. We also stored daily weather data from 1961 - 2008 for meteorological offices throughout the country (about 150 locations). 

further information: http://meteocrop.dc.affrc.go.jp/explan.php

## Run
`docker pull maltegrosse/meteocrop-jp-api:latest`

`docker run -it  maltegrosse/meteocrop-jp-api:latest`

Environmentvariables:
- `AMEDAS_STATION_URL`, default: https://meteocrop.rad.naro.go.jp/real/csv/ame.csv
- `AMEDAS_WEATHER_URL`, default: https://meteocrop.rad.naro.go.jp/real/download.php?kind=5&id=
- `API_PORT`, default: 8080
	
## Usage
`[GET] http://localhost:8080/weather?latitude=35.183334&longitude=136.899994&from_date=1991-01-01&to_date=2021-12-31&output=json`

Parameters:
- latitude: float
- longitude: float
- from_date: string date formate YYYY-MM-DD
- todate: string date formate YYYY-MM-DD
- output: optional. (json, csv, wtd)

The api chooses automatically the closest AMEDAS station.

`[GET] http://localhost:8080/ping`

Response as 'pong' for readiness probes

## Unified Code for Units of Measure (UCUM)
- average_air_temperature (°C)
- maximum_air_temperature (°C)
- minimum_air_temperature (°C)
- precipitation (mm)
- atmospheric_pressure (hPa)
- vapor_pressure (hPa)
- vapor_pressure_deficit (hPa)
- average_relative_humidity (%)
- minimum_relative_humidity (%)
- average_wind_speed (m/s)
- maximum_wind_speed (m/s)
- sunshine_duration (h)
- solar_irradiance_estimate (Wm²)
- downward_longwave_irradiance (Wm²)
- potential_evaporation (mm)
- evapotranspiration_FAO (mm)
- water_temperature_tw0 (°C)
- water_temperature_twinf (°C)

## License
**[MIT license](http://opensource.org/licenses/mit-license.php)**

Copyright 2022 © Malte Grosse.