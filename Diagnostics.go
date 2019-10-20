package oblig1

import (
  "time"
  "net/http"
  "encoding/json"
  "strings"
)

type diagStorage interface {
	Init()
  Get() (Diag, error)
  CheckIfWorks(api string)
  GetAll() []Diag
}

type Diag struct {
  Gbif             			int
  Restcountries         int
  Version             	string
	Uptime             		time.Duration
}

var startTime time.Time                  // GLobal time variable for storing when service started
var DN = diagDB{}                        // Map vaiable for fiagnostics

type diagDB struct {                     // Diagnostics map stored in memory
	diag map[int]Diag
}

func (db *diagDB) Init() {               // Initialised for use
	db.diag = make(map[int]Diag)
  startTime = time.Now()                 // Stores application start Time
  var tempDiag Diag                      // Temp to hold to be modified diagnostics value
  tempDiag = db.diag[0]                  // Copies object
  tempDiag.Gbif = http.StatusOK          // Assigns default start up values
  tempDiag.Restcountries = http.StatusOK
  tempDiag.Version = "v1"
  db.diag[0] = tempDiag                  // Places udated information into diag
}

func (db *diagDB) Get() (Diag, bool){    // Get specific diagnostics
  s, ok := db.diag[0]
	return s, ok
}

func (db *diagDB) TestApi(api string){   // Error handling
  var tempDiag Diag
  tempDiag = db.diag[0]
  if api == "country"{                   // For restcountries
    tempDiag.Gbif = http.StatusServiceUnavailable
  }else if api == "species" || api == "occurrence"{ //for Gbif
    tempDiag.Restcountries = http.StatusServiceUnavailable
  }

  db.diag[0] = tempDiag
}

func (db *diagDB) GetAll() []Diag {     // gets all the diagnostics
  var tempDiag Diag
  tempDiag = db.diag[0]
  tempDiag.Uptime = time.Since(startTime) / time.Second
  db.diag[0] = tempDiag
	all := make([]Diag, 0, 1)
	for _, s := range db.diag {
		all = append(all, s)
	}
	return all
}

                                         // returns teh webservice
func printDiagnostics(w http.ResponseWriter) {
  a := make([]Diag, 0, 1)
  for _, s := range DN.GetAll() {
    a = append(a, s)
  }
  json.NewEncoder(w).Encode(a)
}

func handlerDiag(w http.ResponseWriter, r *http.Request) {
		http.Header.Add(w.Header(), "content-type", "application/json")

		parts := strings.Split(r.URL.Path, "/")

		if len(parts) == 6 || parts[1] != "conservation" {
			http.Error(w, "Bad request:", http.StatusBadRequest)
			return
		}

    printDiagnostics(w)
}
