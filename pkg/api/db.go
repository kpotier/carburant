package api

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/kpotier/carburant/pkg/gas"
)

// MaxHistory variable determines the maximum number of data points that will be
// saved for each service station in the database.
var MaxHistory = 30

func (a *API) dbUpdate() error {
	// Put the prices into db
	a.muxDB.Lock()
	for _, d := range a.data {
		id := fmt.Sprint(d.ID)
		for gn, g := range d.Gas {
			res, ok := a.db[id+string(gn)]
			var last gas.DataGas
			if !ok {
				a.db[id+string(gn)] = make([]gas.DataGas, 0, 1)
			} else if len(res) > 0 {
				last = res[len(res)-1]
			}
			if last.Amount != g.Amount && last.Date != g.Date {
				a.db[id+string(gn)] = append(a.db[id+string(gn)], g)
				if len(a.db[id+string(gn)]) > MaxHistory {
					a.db[id+string(gn)] = a.db[id+string(gn)][1:]
				}
			}
		}
	}
	a.muxDB.Unlock()

	// Encode
	f, err := os.Create(a.pathDB)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	a.muxDB.RLock()
	defer a.muxDB.RUnlock()
	err = enc.Encode(a.db)
	return err
}

func (a *API) dbGet() error {
	_, err := os.Stat(a.pathDB)
	if err != nil {
		if os.IsNotExist(err) {
			a.db = make(map[string][]gas.DataGas)
			return nil
		}
		return err
	}

	f, err := os.Open(a.pathDB)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	a.muxDB.Lock()
	defer a.muxDB.Unlock()
	err = dec.Decode(&a.db)
	return err
}
