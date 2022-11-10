package models

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
)

var Stations AmedasStations

type AmedasStation struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Location  string  `json:"location"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	URL       string  `json:"-"`
}
type tmpAmedasStation struct {
	Titen string  `csv:"titen"`
	Tiiki string  `csv:"tiiki"`
	Lon   float64 `csv:"lon"`
	Lat   float64 `csv:"lat"`
	Html  string  `csv:"html"`
	K     int     `csv:"k"`
}
type AmedasStations struct {
	URL      string
	Stations []AmedasStation
}

func (a *AmedasStations) GetStations(url string, weatherUrl string) error {
	a.URL = url
	client := http.Client{}
	req, err := http.NewRequest("GET", a.URL, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New(fmt.Sprint("error calling url", a.URL, "http-code", res.StatusCode))
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		return gocsv.LazyCSVReader(in)
	})
	tmpStations := []*tmpAmedasStation{}
	if err := gocsv.UnmarshalBytes(respBody, &tmpStations); err != nil {
		return err
	}
	stations := make([]AmedasStation, len(tmpStations))
	for idx, tmpStation := range tmpStations {
		//	fmt.Println(tmpStation.K,tmpStation.Html)
		stations[idx].Latitude = tmpStation.Lat
		stations[idx].Longitude = tmpStation.Lon
		stations[idx].Region = tmpStation.Titen
		stations[idx].Location = tmpStation.Tiiki
		name, id, err := a.handleHtmlString(tmpStation.Html)
		if err != nil {
			return err
		}
		stations[idx].Name = name
		stations[idx].ID = id
		stations[idx].URL = weatherUrl + fmt.Sprint(id)

	}
	a.Stations = stations
	return nil
}
func (a *AmedasStations) handleHtmlString(html string) (name string, id int, err error) {
	left := "<b>"
	right := "</b>"
	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))
	matches := rx.FindAllStringSubmatch(html, -1)

	if len(matches) > 0 && len(matches[0]) > 0 {
		name := matches[0][1]
		id, err := strconv.Atoi(matches[1][1])
		if err != nil {
			return name, id, err
		}
		return name, id, nil
	}

	return name, id, errors.New("could not identify station id or name")
}
func (a *AmedasStations) getDistance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	p := 0.017453292519943295
	v := 0.5 - math.Cos((lat2-lat1)*p)/2 + math.Cos(lat1*p)*math.Cos(lat2*p)*(1-math.Cos((lon2-lon1)*p))/2
	return 2 * 6371 * math.Asin(math.Sqrt(v))

}
func (a *AmedasStations) getClosestStation(lat float64, lon float64) (station AmedasStation, err error) {
	latestStation := -1
	distance := math.MaxFloat64
	for i, s := range a.Stations {
		dist := a.getDistance(lat, lon, s.Latitude, s.Longitude)
		if dist < distance {
			latestStation = i
			distance = dist
		}
	}

	if latestStation < 0 {
		return station, errors.New("could not found station by given lat/lon")
	}
	return a.Stations[latestStation], nil
}
func (a *AmedasStation) GetWeatherData(fromDate time.Time, toDate time.Time) (h HistoricalWeatherObservations, err error) {
	
	client := http.Client{}
	req, err := http.NewRequest("GET", a.URL, nil)
	if err != nil {
		return h, err
	}

	res, err := client.Do(req)
	if err != nil {
		return h, err
	}

	if res.StatusCode != 200 {
		return h, errors.New(fmt.Sprint("error calling url", a.URL, "http-code", res.StatusCode))
	}
	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return h, err
	}
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comment = ' '
		r.TrimLeadingSpace = true
		return r
	})
	respBody = bytes.TrimSpace(respBody)

	tmpObservation := []*WeatherObservation{}
	h.Station = a
	if err := gocsv.UnmarshalBytes(respBody, &tmpObservation); err != nil {
		return h, err
	}
	observations := []*WeatherObservation{}
	for _, ob := range tmpObservation {
		ob.Date = time.Date(ob.Year, time.Month(ob.Month), ob.Day, 0, 0, 0, 0, time.UTC)
		ob.DayOfYear = ob.Date.YearDay()
		if ob.Date.After(fromDate.AddDate(0, 0, -1)) && ob.Date.Before(toDate.AddDate(0, 0, 1)) {
			observations = append(observations, ob)
		}
	}

	h.Observations = observations
	return h, nil
}
func (a *AmedasStations) GetWeatherObservation(lat float64, lon float64, fromDate time.Time, toDate time.Time) (h HistoricalWeatherObservations, err error) {
	err = a.validateLatLon(lat, lon)
	if err != nil {
		return h, err
	}
	station, err := a.getClosestStation(lat, lon)
	if err != nil {
		return h, err
	}
	return station.GetWeatherData(fromDate, toDate)
}
func (a *AmedasStations) validateLatLon(lat float64, lon float64) error {
	if lat > 90 || lat < -90 || lon > 180 || lon < -180 {
		return errors.New("wrong latitude and longitude")
	}

	return nil
}
