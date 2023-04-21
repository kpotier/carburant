package api

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/kpotier/carburant/pkg/gas"
)

func TestStart(t *testing.T) {
	gas.DataURL = "http://127.0.0.1:8081/"
	Refresh = 1 * time.Second
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("content-type", "application/json; charset=utf-8")
			w.Write([]byte(CONTENT))
		})
		http.ListenAndServe(":8081", nil)
	}()

	time.Sleep(2 * time.Second)

	a, err := Start(log.Default(), filepath.Join(t.TempDir(), "db.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer a.Stop()

	go func() {
		http.ListenAndServe(":8080", a.Handler()) // or just a
	}()

	time.Sleep(2 * time.Second)

	r, err := http.PostForm("http://127.0.0.1:8080/stations", url.Values{"lng": {"7.186"}, "lat": {"43.73"}, "gas": {"E10"}, "lim": {"1"}})
	if err != nil {
		t.Fatal(err)
	}

	str, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(str) != WANT {
		t.Fatalf("bad body stations: got %s, want %s", string(str), WANT)
	}

	r, err = http.PostForm("http://127.0.0.1:8080/history", url.Values{"id": {"57430004"}, "gas": {"E10"}})
	if err != nil {
		t.Fatal(err)
	}

	str, err = io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(str) != WANTHISTORY {
		t.Fatalf("bad body history: got %s, want %s", string(str), WANTHISTORY)
	}

	r, err = http.PostForm("http://127.0.0.1:8080/favorites", url.Values{"lng": {"7.186"}, "lat": {"43.73"}, "list": {"57430004,13380002"}})
	if err != nil {
		t.Fatal(err)
	}

	str, err = io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(str) != WANTFAVORITES {
		t.Fatalf("bad body favorites: got %s, want %s", string(str), WANTFAVORITES)
	}
}

const CONTENT string = `[{"datasetid": "prix-des-carburants-en-france-flux-instantane-v2", "recordid": "5da5a0625e55e658d40138ba11ff01c89aa41236", "fields": {"services": "{\"service\": [\"Boutique alimentaire\", \"Boutique non alimentaire\", \"Restauration \u00e0 emporter\", \"Carburant additiv\u00e9\", \"Vente de gaz domestique (Butane, Propane)\", \"Vente d'additifs carburants\", \"DAB (Distributeur automatique de billets)\"]}", "region": "Grand Est", "gazole_maj": "2023-04-06 06:47:40", "id": 57430004, "horaires": "{\"@automate-24-24\": \"\", \"jour\": [{\"@id\": \"1\", \"@nom\": \"Lundi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"07.00\", \"@fermeture\": \"13.00\"}}, {\"@id\": \"2\", \"@nom\": \"Mardi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"07.00\", \"@fermeture\": \"13.00\"}}, {\"@id\": \"3\", \"@nom\": \"Mercredi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"07.00\", \"@fermeture\": \"13.00\"}}, {\"@id\": \"4\", \"@nom\": \"Jeudi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"07.00\", \"@fermeture\": \"13.00\"}}, {\"@id\": \"5\", \"@nom\": \"Vendredi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"07.00\", \"@fermeture\": \"13.00\"}}, {\"@id\": \"6\", \"@nom\": \"Samedi\", \"@ferme\": \"1\"}, {\"@id\": \"7\", \"@nom\": \"Dimanche\", \"@ferme\": \"1\"}]}", "gazole_prix": "1.883", "pop": "R", "horaires_automate_24_24": "Non", "e10_maj": "2023-04-06 06:47:41", "cp": "57430", "prix": "[{\"@nom\": \"Gazole\", \"@id\": \"1\", \"@maj\": \"2023-04-06 06:47:40\", \"@valeur\": \"1.883\"}, {\"@nom\": \"E10\", \"@id\": \"5\", \"@maj\": \"2023-04-06 06:47:41\", \"@valeur\": \"1.956\"}, {\"@nom\": \"SP98\", \"@id\": \"6\", \"@maj\": \"2023-04-06 06:47:41\", \"@valeur\": \"2.066\"}]", "ville": "Sarralbe", "code_region": "44", "code_departement": "57", "latitude": "4898700", "departement": "Moselle", "carburants_indisponibles": "SP95;E85;GPLc", "geom": [48.987, 7.028], "sp98_maj": "2023-04-06 06:47:41", "sp98_prix": "2.066", "services_service": "Boutique alimentaire//Boutique non alimentaire//Restauration \u00e0 emporter//Carburant additiv\u00e9//Vente de gaz domestique (Butane, Propane)//Vente d'additifs carburants//DAB (Distributeur automatique de billets)", "e10_prix": "1.956", "carburants_disponibles": "Gazole;E10;SP98", "adresse": "97 Route de Strasbourg", "longitude": 702800.0}, "geometry": {"type": "Point", "coordinates": [7.028, 48.987]}, "record_timestamp": "2023-04-10T19:04:00.926+02:00"},{"datasetid": "prix-des-carburants-en-france-flux-instantane-v2", "recordid": "387a57a51a82a211922f51dbf22772eeea2eebb2", "fields": {"region": "Provence-Alpes-C\u00f4te d'Azur", "gazole_maj": "2023-04-08 04:02:00", "id": 13380002, "horaires": "{\"@automate-24-24\": \"\", \"jour\": [{\"@id\": \"1\", \"@nom\": \"Lundi\", \"@ferme\": \"\"}, {\"@id\": \"2\", \"@nom\": \"Mardi\", \"@ferme\": \"\"}, {\"@id\": \"3\", \"@nom\": \"Mercredi\", \"@ferme\": \"\"}, {\"@id\": \"4\", \"@nom\": \"Jeudi\", \"@ferme\": \"\"}, {\"@id\": \"5\", \"@nom\": \"Vendredi\", \"@ferme\": \"\"}, {\"@id\": \"6\", \"@nom\": \"Samedi\", \"@ferme\": \"\"}, {\"@id\": \"7\", \"@nom\": \"Dimanche\", \"@ferme\": \"\"}]}", "gazole_prix": "1.863", "pop": "R", "horaires_automate_24_24": "Non", "e10_maj": "2023-04-08 04:03:00", "cp": "13380", "prix": "[{\"@nom\": \"Gazole\", \"@id\": \"1\", \"@maj\": \"2023-04-08 04:02:00\", \"@valeur\": \"1.863\"}, {\"@nom\": \"E85\", \"@id\": \"3\", \"@maj\": \"2023-04-08 04:03:00\", \"@valeur\": \"1.128\"}, {\"@nom\": \"E10\", \"@id\": \"5\", \"@maj\": \"2023-04-08 04:03:00\", \"@valeur\": \"1.971\"}, {\"@nom\": \"SP98\", \"@id\": \"6\", \"@maj\": \"2023-04-08 04:03:00\", \"@valeur\": \"2.051\"}]", "ville": "Plan-de-Cuques", "code_region": "93", "code_departement": "13", "latitude": "4334300", "departement": "Bouches-du-Rh\u00f4ne", "carburants_indisponibles": "SP95;GPLc", "geom": [43.343, 5.457], "sp98_maj": "2023-04-08 04:03:00", "sp98_prix": "2.051", "e85_prix": "1.128", "e10_prix": "1.971", "e85_maj": "2023-04-08 04:03:00", "carburants_disponibles": "Gazole;E85;E10;SP98", "adresse": "89 AVENUE DE LA LIBERATION", "longitude": 545700.0}, "geometry": {"type": "Point", "coordinates": [5.457, 43.343]}, "record_timestamp": "2023-04-10T19:04:00.926+02:00"},{"datasetid": "prix-des-carburants-en-france-flux-instantane-v2", "recordid": "756c933d7bbcae73051369ea81dfe4634982e807", "fields": {"services": "{\"service\": [\"Vente de p\u00e9trole lampant\", \"Station de gonflage\", \"Carburant additiv\u00e9\", \"Piste poids lourds\", \"Vente de gaz domestique (Butane, Propane)\", \"Automate CB 24/24\"]}", "region": "Provence-Alpes-C\u00f4te d'Azur", "gazole_maj": "2023-04-10 13:26:00", "id": 6201001, "horaires": "{\"@automate-24-24\": \"1\", \"jour\": [{\"@id\": \"1\", \"@nom\": \"Lundi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}, {\"@id\": \"2\", \"@nom\": \"Mardi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}, {\"@id\": \"3\", \"@nom\": \"Mercredi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}, {\"@id\": \"4\", \"@nom\": \"Jeudi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}, {\"@id\": \"5\", \"@nom\": \"Vendredi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}, {\"@id\": \"6\", \"@nom\": \"Samedi\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}, {\"@id\": \"7\", \"@nom\": \"Dimanche\", \"@ferme\": \"\", \"horaire\": {\"@ouverture\": \"00.00\", \"@fermeture\": \"00.00\"}}]}", "gazole_prix": "1.773", "pop": "R", "horaires_automate_24_24": "Oui", "sp95_maj": "2023-04-10 13:26:00", "e10_maj": "2023-04-10 13:26:00", "sp95_prix": "1.961", "cp": "06200", "prix": "[{\"@nom\": \"Gazole\", \"@id\": \"1\", \"@maj\": \"2023-04-10 13:26:00\", \"@valeur\": \"1.773\"}, {\"@nom\": \"SP95\", \"@id\": \"2\", \"@maj\": \"2023-04-10 13:26:00\", \"@valeur\": \"1.961\"}, {\"@nom\": \"GPLc\", \"@id\": \"4\", \"@maj\": \"2023-04-10 13:26:00\", \"@valeur\": \"0.958\"}, {\"@nom\": \"E10\", \"@id\": \"5\", \"@maj\": \"2023-04-10 13:26:00\", \"@valeur\": \"1.914\"}, {\"@nom\": \"SP98\", \"@id\": \"6\", \"@maj\": \"2023-04-10 13:26:00\", \"@valeur\": \"1.974\"}]", "ville": "Nice", "code_region": "93", "code_departement": "06", "latitude": "4373000", "departement": "Alpes-Maritimes", "carburants_indisponibles": "E85", "geom": [43.73, 7.186], "sp98_maj": "2023-04-10 13:26:00", "sp98_prix": "1.974", "services_service": "Vente de p\u00e9trole lampant//Station de gonflage//Carburant additiv\u00e9//Piste poids lourds//Vente de gaz domestique (Butane, Propane)//Automate CB 24/24", "gplc_prix": "0.958", "gplc_maj": "2023-04-10 13:26:00", "e10_prix": "1.914", "carburants_disponibles": "Gazole;SP95;GPLc;E10;SP98", "adresse": "606 boulevard du mercantour", "longitude": 718600.0}, "geometry": {"type": "Point", "coordinates": [7.186, 43.73]}, "record_timestamp": "2023-04-10T19:04:00.926+02:00"}]`

const WANT = `[{"distance":0.000051956875,"id":6201001,"coords":[43.73,7.186],"address_rd":"606 boulevard du mercantour","address_cp":"06200 Nice","automate_2424":true,"horaires":[[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]],[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]],[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]],[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]],[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]],[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]],[[{"Hour":0,"Minutes":0},{"Hour":0,"Minutes":0}]]],"services":["Vente de pétrole lampant","Station de gonflage","Carburant additivé","Piste poids lourds","Vente de gaz domestique (Butane, Propane)","Automate CB 24/24"],"gas":{"E10":{"date":"2023-04-10T13:26:00Z","amount":1914},"GPLc":{"date":"2023-04-10T13:26:00Z","amount":958},"Gazole":{"date":"2023-04-10T13:26:00Z","amount":1773},"SP95":{"date":"2023-04-10T13:26:00Z","amount":1961},"SP98":{"date":"2023-04-10T13:26:00Z","amount":1974}}}]
`

const WANTHISTORY = `[{"date":"2023-04-06T06:47:41Z","amount":1956}]
`

const WANTFAVORITES = `[{"distance":584.6775,"id":57430004,"coords":[48.987,7.028],"address_rd":"97 Route de Strasbourg","address_cp":"57430 Sarralbe","automate_2424":false,"horaires":[[[{"Hour":7,"Minutes":0},{"Hour":13,"Minutes":0}]],[[{"Hour":7,"Minutes":0},{"Hour":13,"Minutes":0}]],[[{"Hour":7,"Minutes":0},{"Hour":13,"Minutes":0}]],[[{"Hour":7,"Minutes":0},{"Hour":13,"Minutes":0}]],[[{"Hour":7,"Minutes":0},{"Hour":13,"Minutes":0}]],null,null],"services":["Boutique alimentaire","Boutique non alimentaire","Restauration à emporter","Carburant additivé","Vente de gaz domestique (Butane, Propane)","Vente d'additifs carburants","DAB (Distributeur automatique de billets)"],"gas":{"E10":{"date":"2023-04-06T06:47:41Z","amount":1956},"Gazole":{"date":"2023-04-06T06:47:40Z","amount":1883},"SP98":{"date":"2023-04-06T06:47:41Z","amount":2066}}},{"distance":145.86531,"id":13380002,"coords":[43.343,5.457],"address_rd":"89 AVENUE DE LA LIBERATION","address_cp":"13380 Plan-de-Cuques","automate_2424":false,"horaires":[[],[],[],[],[],[],[]],"services":[""],"gas":{"E10":{"date":"2023-04-08T04:03:00Z","amount":1971},"E85":{"date":"2023-04-08T04:03:00Z","amount":1128},"Gazole":{"date":"2023-04-08T04:02:00Z","amount":1863},"SP98":{"date":"2023-04-08T04:03:00Z","amount":2051}}}]
`
