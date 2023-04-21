// gas package includes data structures for service stations, such as their ID,
// address, and coordinates. It also holds information on available gases and
// prices, opening hours, automated services, and offerings provided by each
// station.
package gas

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DataURL is the URL to retrieve the price of each ful in each service station
// in France. It must points to a JSON-encoded file.
var DataURL = "https://www.data.gouv.fr/fr/datasets/r/b0561905-7b5e-4f38-be50-df05708acb80"

// Fetch retrieves and decodes the JSON file available at the URL DataURL. It
// returns the list of stations and the price of each fuel.
func Fetch() (d []Data, err error) {
	r, err := http.Get(DataURL)
	if err != nil {
		return
	}
	if r.StatusCode != 200 {
		err = fmt.Errorf("r.StatusCode is %d, want 200", r.StatusCode)
		return
	}
	if r.Header.Get("content-type") != "application/json; charset=utf-8" {
		err = fmt.Errorf("content-type is not `application/json; charset=utf-8` but `%s`",
			r.Header.Get("content-type"))
		return
	}

	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&d)
	return
}
