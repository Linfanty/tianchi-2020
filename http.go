package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func installHttpHandlers(player *Player) {
	http.HandleFunc("/ready", func(writer http.ResponseWriter, request *http.Request) {
		ready := player.Ready()
		statusCode := http.StatusOK
		if !ready {
			statusCode = http.StatusInternalServerError
		}
		http.Error(writer, "", statusCode)
	})

	// If want to reuse without restarting, call this API each time before start.
	http.HandleFunc("/reset", func(writer http.ResponseWriter, request *http.Request) {
		if err := player.Reset(); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Error(writer, "", http.StatusOK)
	})

	http.HandleFunc("/p1_start", func(writer http.ResponseWriter, request *http.Request) {
		var (
			delay = apiDelay
		)
		query := request.URL.Query()
		if v := query.Get("delay"); v != "" {
			if d, err := time.ParseDuration(v); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				delay = d
			}
		}

		// Parse pilots from req body.
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var pilots struct {
			Pilots []string `json:"pilots"`
		}
		if err := json.Unmarshal(bs, &pilots); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if delay > 0 {
			time.Sleep(delay)
		}

		res, err := player.P1(pilots.Pilots)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if bs, err := json.Marshal(res); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, string(bs), http.StatusOK)
		}
	})

	http.HandleFunc("/p2_start", func(writer http.ResponseWriter, request *http.Request) {
		var (
			delay = apiDelay
		)
		query := request.URL.Query()
		if v := query.Get("delay"); v != "" {
			if d, err := time.ParseDuration(v); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			} else {
				delay = d
			}
		}

		// Parse apps and dependencies from req body.
		bs, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var param pParams
		if err := json.Unmarshal(bs, &param); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if delay > 0 {
			time.Sleep(delay)
		}

		res, err := player.P2(param)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if bs, err := json.Marshal(res); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(writer, string(bs), http.StatusOK)
		}
	})
}
