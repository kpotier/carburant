package gas

import (
	"time"
)

// Data represents a service-station. It contains the coordinates, identifier,
// opening and closing hours.
type Data struct {
	// ID serves as a unique identifier for a service station, allowing it to be
	// distinguished from other locations or entities within a system or
	// database.
	ID int `json:"id"`
	// Coords refer to the geographic coordinates of a service station, which
	// consist of two values indicating its longitude and latitude respectively.
	Coords [2]float64 `json:"coords"`
	// Address refers to a location identifier that follows a specific format
	// consisting of the street number and street name.
	AddressRd string `json:"address_rd"`
	// Same as AddressRd but it contains the zip code, and city name.
	AddressCP string `json:"address_cp"`

	// Automate2424 indicates whether a service station is open 24/7 due to the
	// presence of an automated system, which allows customers to access
	// services even when the station is not staffed.
	Automate2424 bool `json:"automate_2424"`
	// Horaires refer to the specific opening and closing hours of a service
	// station on each day of the week.
	Horaires [7][][2]DataTime `json:"horaires"`
	// Services refer to a collection of different offerings that a business or
	// establishment provides to its customers, which may include products,
	// support, maintenance, or other types of assistance.
	Services []string `json:"services"`

	// Gas data structure contains information about the various types of gases
	// available at a service station, including the cost of each gas.
	Gas map[GasType]DataGas `json:"gas"`
}

// DateTime refers to a specific time of day that includes only the hour and
// minute values.
type DataTime struct {
	Hour    int
	Minutes int
}

// DataGas is a type of data that contains information about a particular type
// of gas, including the time when its price was last updated and its cost per
// liter.
type DataGas struct {
	Date   time.Time `json:"date"`
	Amount int       `json:"amount"`
}

// GasType is the name of the gas.
type GasType string

// A record of gas names that can be categorized by properties or applications.
const (
	GasGazole GasType = "Gazole"
	GasSP95   GasType = "SP95"
	GasE85    GasType = "E85"
	GasGPLc   GasType = "GPLc"
	GasE10    GasType = "E10"
	GasSP98   GasType = "SP98"
)
