package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/render"

	"github.com/maltegrosse/meteocrop-jp-api/models"
)

func WeatherObservations(w http.ResponseWriter, r *http.Request) {

	lat, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "invalid latitude "+err.Error(), 400)
		return
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil {
		http.Error(w, "invalid longitude "+err.Error(), 400)
		return
	}

	from, err := time.Parse("2006-1-2", r.URL.Query().Get("from_date"))
	if err != nil {
		http.Error(w, "invalid from date "+err.Error(), 400)
		return
	}
	to, err := time.Parse("2006-1-2", r.URL.Query().Get("to_date"))
	if err != nil {
		http.Error(w, "invalid to date "+err.Error(), 400)
		return
	}

	h, err := models.Stations.GetWeatherObservation(lat, lon, from, to)
	if err != nil {
		http.Error(w, "error getting data "+err.Error(), 500)
		return
	}
	h.Latitude = lat
	h.Longitude = lon
	h.FromDate = from
	h.ToDate = to
	output := r.URL.Query().Get("output")
	filename := fmt.Sprint(lat, "_", lon, "-", from.Format("2006-01-02"), "-", to.Format("2006-01-02"))
	switch {
	case output == "csv":
		w.Header().Add("Content-Disposition", `attachment; filename="`+filename+`.csv"`)
		w.Header().Set("Content-Type", "text/csv")
		csvContent, err := h.ConvertToCsv()
		if err != nil {
			http.Error(w, "error exporting to csv"+err.Error(), 500)
			return
		}
		w.Write([]byte(csvContent))

	case output == "wtd":
		w.Header().Add("Content-Disposition", `attachment; filename="`+filename+`.csv"`)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(h.ConvertToWtd()))
	default:
		render.JSON(w, r, &h)
	}

}
