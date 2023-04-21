package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/kpotier/carburant/pkg/gas"
)

const degToRad = 0.017453292519943295

type Result struct {
	Distance float32 `json:"distance"`
	gas.Data
}

func (a *API) getStations(w http.ResponseWriter, r *http.Request) {
	err := a.parseForm(w, r)
	if err != nil {
		a.log.Printf("parseForm: %s", err)
		w.Write([]byte(`{"error": "POST_ERR"}`))
		return
	}
	lng, lat, err := a.parseLngLat(r)
	if err != nil {
		a.log.Printf("parseLngLat: %s", err)
		w.Write([]byte(`{"error": "POST_ERR"}`))
		return
	}

	limStr := r.PostFormValue("lim")
	lim, err := strconv.Atoi(limStr)
	if err != nil {
		a.log.Printf("lim: %s", err)
		w.Write([]byte(`{"error": "POST_ERR"}`))
		return
	}
	g := r.PostFormValue("gas")

	a.mux.RLock()
	defer a.mux.RUnlock()

	var list = make([]Result, lim)
	var sort = make([]int, 0, lim)

	for i := 0; i < lim; i++ {
		var distance float64 = math.MaxFloat64
		var id int
		for j, d := range a.data {
			// Does it has the selected gas?
			if _, ok := d.Gas[gas.GasType(g)]; !ok {
				continue
			}

			// Already in the list?
			var found bool
			for _, k := range sort {
				if k == j {
					found = true
					break
				}
			}
			if found {
				continue
			}

			dist := calcDistance(d.Coords[1], d.Coords[0], lng, lat)
			if dist < distance {
				distance = dist
				id = j
			}
		}
		sort = append(sort, id)
		list[i] = Result{endDistance(distance), a.data[id]}
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(list)
	if err != nil {
		a.log.Println(err)
	}
}

func (a *API) getHistory(w http.ResponseWriter, r *http.Request) {
	err := a.parseForm(w, r)
	if err != nil {
		a.log.Printf("parseForm: %s", err)
		w.Write([]byte(`{"error": "POST_ERR"}`))
		return
	}

	id := r.PostFormValue("id")
	gas := r.PostFormValue("gas")

	a.muxDB.RLock()
	defer a.muxDB.RUnlock()
	enc := json.NewEncoder(w)
	err = enc.Encode(a.db[id+gas])
	if err != nil {
		a.log.Println(err)
	}
}

func (a *API) getFavorites(w http.ResponseWriter, r *http.Request) {
	err := a.parseForm(w, r)
	if err != nil {
		a.log.Printf("parseForm: %s", err)
		w.Write([]byte(`{"error": "POST_ERR"}`))
		return
	}
	lng, lat, err := a.parseLngLat(r)
	if err != nil {
		a.log.Printf("parseLngLat: %s", err)
		w.Write([]byte(`{"error": "POST_ERR"}`))
		return
	}

	listStr := r.PostFormValue("list")
	listSt := strings.Split(listStr, ",")
	var listInt = make([]int, len(listSt))
	for i, l := range listSt {
		listInt[i], err = strconv.Atoi(l)
		if err != nil {
			a.log.Printf("list: %s", err)
			w.Write([]byte(`{"error": "POST_ERR"}`))
			return
		}
	}

	a.mux.RLock()
	defer a.mux.RUnlock()

	var list = make([]Result, 0, len(listInt))

	for _, d := range a.data {
		// In the list?
		var found bool
		for _, k := range listInt {
			if k == d.ID {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		dist := calcDistance(d.Coords[1], d.Coords[0], lng, lat)
		list = append(list, Result{endDistance(dist), d})
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(list)
	if err != nil {
		a.log.Println(err)
	}
}

func calcDistance(lng1, lat1, lng2, lat2 float64) float64 {
	// Faster but a little bit more inaccurate than the Haversine
	// formula. As we only deal with short distances, this formula
	// should not lead to problems.
	x := degToRad * (lng1 - lng2) * math.Cos(0.5*degToRad*(lat2+lat1))
	y := degToRad * (lat1 - lat2)
	return x*x + y*y
}

func endDistance(dist float64) float32 {
	return float32(math.Sqrt(dist) * 6371.)
}

func (a *API) parseForm(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("content-type", "application/json; charset=utf-8")
	if r.Method != "POST" {
		return fmt.Errorf("getStations: method is %s, want %s", r.Method, "POST")
	}
	return r.ParseForm()
}

func (a *API) parseLngLat(r *http.Request) (float64, float64, error) {
	lngStr := r.PostFormValue("lng")
	latStr := r.PostFormValue("lat")

	lng, err := strconv.ParseFloat(lngStr, 32)
	if err != nil {
		return 0, 0, err
	}
	lat, err := strconv.ParseFloat(latStr, 32)
	if err != nil {
		return 0, 0, err
	}
	return lng, lat, nil
}
