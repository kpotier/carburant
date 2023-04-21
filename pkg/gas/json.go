package gas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// JSONFields contain the useful JSON data for each station. This is used to
// decode the JSON file on the fly. A post-processing will then be done from
// this data.
type JSONFields struct {
	Fields struct {
		ID   int        `json:"id"`
		Geom [2]float64 `json:"geom"`

		Adresse string `json:"adresse"`
		// CP is the zip code.
		CP    string `json:"cp"`
		Ville string `json:"ville"`

		// HorairesAutomate2424 is set to "Oui" when there is an automate
		// available 24/7.
		HorairesAutomate2424 string `json:"horaires_automate_24_24"`

		// Horaires are the shop/counter opening hours. The format is as follows:
		//
		// {
		//   "@automate-24-24": "1",
		//   "jour": [
		//    {
		//     "@id": "1",
		//     "@nom": "Lundi",
		//     "@ferme": "",
		//     "horaire": {"@ouverture": "08.30", "@fermeture": "19.00"}
		//    },
		//    {
		//     "@id": "2",
		//     "@nom": "Lundi",
		//     "@ferme": "1",
		//     "horaire": [{"@ouverture": "08.30", "@fermeture": "19.00"}, {"@ouverture": "08.30", "@fermeture": "19.00"}]
		//    },
		//   ]
		// }
		Horaires string `json:"horaires"`

		// ServiceService is a "//" separated list of services.
		ServicesService string `json:"services_service"`
		// CarburantsDisponibles is a ";" separated list of available fuels.
		CarburantsDisponibles string `json:"carburants_disponibles"`

		// GazoleMaj is the date of the latest update for the "Gazole" fuel. The
		// format is as follows: 2006-01-02 15:04:05.
		GazoleMaj string `json:"gazole_maj"`
		// GazolePrix is the price of the "Gazole" fuel in Euros. The format is
		// as follows: 1.02.
		GazolePrix string `json:"gazole_prix"`
		SP95Maj    string `json:"sp95_maj"`
		SP95Prix   string `json:"sp95_prix"`
		E85Maj     string `json:"e85_maj"`
		E85Prix    string `json:"e85_prix"`
		GPLcMaj    string `json:"gplc_maj"`
		GPLcPrix   string `json:"gplc_prix"`
		E10Maj     string `json:"e10_maj"`
		E10Prix    string `json:"e10_prix"`
		SP98Maj    string `json:"sp98_maj"`
		SP98Prix   string `json:"sp98_prix"`
	} `json:"fields"`
}

// JSONHoraires hold the shop/counter opening hours. This is used to decode the
// JSON file on the fly. A post-processing will then be done from this data.
type JSONHoraires struct {
	// Automate2424 is set to "1" when there is an automate available 24/7.
	Automate2424 string `json:"@automate-24-24"`
	// Jour contain the opening hours for each day.
	Jour []struct {
		// ID of the day. For example: 1 = Monday.
		ID string `json:"@id"`
		// Ferme is set to "1" if it the shop/counter is closed.
		Ferme string `json:"@ferme"`
		// Horaire hold the opening hours for this day.
		Horaire JSONHoraire `json:"horaire"`
	} `json:"jour"`
}

// JSONHoraire hold the shop/counter opening hours for a specific day. This is
// used to decode the JSON file on the fly. A post-processing will then be done
// from this data.
type JSONHoraire []JSONOuverture

// JSONOuvertures contain the opening and closing hour of a single time slot.
// This is used to decode the JSON file on the fly. A post-processing will then
// be done from this data.
type JSONOuverture struct {
	// Ouverture is the opening hour. The format is as follows: 15.04.
	Ouverture string `json:"@ouverture"`
	// Fermeture is the closing hour. The format is as follows: 15.04.
	Fermeture string `json:"@fermeture"`
}

func (h *JSONHoraire) UnmarshalJSON(data []byte) error {
	// See format of Horaires in JSONFields struct.
	var err error
	var hNew []JSONOuverture
	if data[0] == '{' {
		var h1 JSONOuverture
		err = json.Unmarshal(data, &h1)
		hNew = append(hNew, h1)
	} else {
		err = json.Unmarshal(data, &hNew)
	}
	*h = hNew
	return err
}

func (d *Data) UnmarshalJSON(data []byte) error {
	// On the fly decoding
	var fields JSONFields
	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}

	// Generate address, ID and coords from the decoded JSON file.
	d.ID = fields.Fields.ID
	d.Coords = fields.Fields.Geom
	d.AddressRd = fields.Fields.Adresse
	d.AddressCP = fields.Fields.CP + " " + fields.Fields.Ville

	// shop/counter opening and closing hours.
	var horaires JSONHoraires
	if fields.Fields.HorairesAutomate2424 == "Oui" {
		d.Automate2424 = true
	}
	if json.Valid([]byte(fields.Fields.Horaires)) {
		err = json.Unmarshal([]byte(fields.Fields.Horaires), &horaires)
		if err != nil {
			return err
		}
	}
	for _, j := range horaires.Jour {
		if j.Ferme == "1" {
			continue
		}
		id, err := strconv.Atoi(j.ID)
		if err != nil {
			return err
		}
		if id < 1 || id > 7 {
			return fmt.Errorf("invalid id: %d", id)
		}
		d.Horaires[id-1] = make([][2]DataTime, len(j.Horaire))
		for i, h := range j.Horaire {
			o, err := toDataTime(h.Ouverture)
			if err != nil {
				return err
			}
			c, err := toDataTime(h.Fermeture)
			if err != nil {
				return err
			}
			d.Horaires[id-1][i] = [2]DataTime{o, c}
		}
	}
	d.Services = strings.Split(fields.Fields.ServicesService, "//")

	if fields.Fields.CarburantsDisponibles != "" {
		gas := strings.Split(fields.Fields.CarburantsDisponibles, ";")
		d.Gas = make(map[GasType]DataGas, len(gas))
		for _, g := range gas {
			var amount string
			var lastUpdate string
			switch g {
			case string(GasE10):
				lastUpdate = fields.Fields.E10Maj
				amount = fields.Fields.E10Prix
			case string(GasE85):
				lastUpdate = fields.Fields.E85Maj
				amount = fields.Fields.E85Prix
			case string(GasGPLc):
				lastUpdate = fields.Fields.GPLcMaj
				amount = fields.Fields.GPLcPrix
			case string(GasGazole):
				lastUpdate = fields.Fields.GazoleMaj
				amount = fields.Fields.GazolePrix
			case string(GasSP95):
				lastUpdate = fields.Fields.SP95Maj
				amount = fields.Fields.SP95Prix
			case string(GasSP98):
				lastUpdate = fields.Fields.SP98Maj
				amount = fields.Fields.SP98Prix
			default:
				return fmt.Errorf("unknown gas %s", g)
			}

			idx := strings.IndexByte(amount, '.')
			if idx < 1 || idx == len(amount)-1 {
				return fmt.Errorf("out of range `.` in amount: %s", amount)
			}
			amountInt, err := strconv.Atoi(amount[:idx])
			if err != nil {
				return err
			}
			amountInt *= 1000
			minor, err := strconv.Atoi(amount[idx+1:])
			if err != nil {
				return err
			}
			amountInt += minor

			date, err := time.Parse("2006-01-02 15:04:05", lastUpdate)
			if err != nil {
				return err
			}
			d.Gas[GasType(g)] = DataGas{date, amountInt}
		}
	}

	return nil
}

func toDataTime(data string) (d DataTime, err error) {
	idx := strings.IndexByte(data, '.')
	if idx < 0 || idx == len(data)-1 {
		err = fmt.Errorf("cannot find `.` or is located at the end of the string for `%s`", data)
		return
	}
	d.Hour, err = strconv.Atoi(data[:idx])
	if err != nil {
		return
	}
	d.Minutes, err = strconv.Atoi(data[idx+1:])
	if err != nil {
		return
	}
	return
}
