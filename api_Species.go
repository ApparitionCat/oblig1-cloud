package oblig1

import (
	"encoding/json"
	"net/http"
	"strings"
	"bytes"
	"strconv"
)

var DBs = SpeciesDB{}

var structure = new(ResultList)

type YearHolder struct{
	Year									string
}





func GetSpeciesJSON(url string) ([]Species, error) {	// Fetches specific json reauested
	resp, err := http.Get(url)										            // Gets information
	if err != nil {																// error handler
		return nil, err
	}
	defer resp.Body.Close()												        // Opens the html body
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &structure); err != nil {	            // stores json into structure
		return nil, err															// error handler
	}
	for idx, row := range structure.AllSpecies {	// For each species
		if idx == 0 {															// Skip first slot
			continue
		}
		url := "http://api.gbif.org/v1/species/"		                        // Open the url for each species to get date
		url = url + strconv.Itoa(int(row.Key)) + "/name"	                    // Gets and uses species key to find the date
		resp, err := http.Get(url)									            // Gets information from this url
		if err != nil {
			return nil, err														// Error handling
		}
		var yearFetch YearHolder										        // Temp for result containing date
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respByte := buf.Bytes()
		json.Unmarshal(respByte, &yearFetch)				                    // Copies json into variable
		row.Year = yearFetch.Year										        // Assigns the fetched year to species
		structure.AllSpecies[idx] = row							                // Places species in map
	}

	return structure.AllSpecies, nil							                // Returns map with structures
}

func replyWithAllSpecies(w http.ResponseWriter, DB SpeciesStorage) {
		url := "http://api.gbif.org/v1/species?offset=20000&limit=25"
		speciesList, err := GetSpeciesJSON(url)
		if err != nil {							                               //error handling, gives the user a message according to the error
			http.Error(w, "Service could not be accessed", http.StatusServiceUnavailable)
		}

		for idx, row := range speciesList {					                   // For each species, if the idx is 0, keep looping
		  if idx == 0 {
		  	continue
	  	}
      DBs.Add(row)															    // Put in main Species struct
    }


		if DB.Count() == 0 {												    // If there are no species, write an empty json
			json.NewEncoder(w).Encode([]Species{})
		} else {																// If there are species stored then they are fetched and printed oot
			a := make([]Species, 0, DB.Count())
			for _, s := range DB.GetAll() {
				a = append(a, s)
			}
			json.NewEncoder(w).Encode(a)							            // And print as JSON to website
	}
}


func replyWithSpecies(w http.ResponseWriter,DB SpeciesStorage, id string) {    //replies with a specific species

	    var singleFetch Species														// Temporary variable used for storing species
		url := "http://api.gbif.org/v1/species/" + id			                // Uses the url combined with id to find species
		resp, err := http.Get(url)												// gets url
		if err != nil {
			http.Error(w, "The species service is currently unavailable", http.StatusServiceUnavailable)
			DN.TestApi("species")											    // Error handler, returns the 503 code: service unavaliable
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respByte := buf.Bytes()
		json.Unmarshal(respByte, &singleFetch)						            //This block of code gets the species requested

		DBs.Add(singleFetch)													// and this code stores the species
		idINT, err := strconv.ParseUint(id, 10, 64)				                // Casts it to unit64
		if err != nil {
			http.Error(w, "The species service is currently unavailable", http.StatusServiceUnavailable)
			DN.TestApi("species")                                               // Error handler, returns the 503 code: service unavaliable
		}
		s, ok := DBs.Get(idINT)												    // Gets teh species with the key that is like id.
		if !ok {
			http.Error(w, "Species not found", http.StatusNotFound)
			return																// Error handler, returns code 404 not found.
		}
		json.NewEncoder(w).Encode(s)											// Prints relevant info as JSON to website
        }


func HandlerSpecies(w http.ResponseWriter, r *http.Request) {
http.Header.Add(w.Header(), "content-type", "application/json")	                // Sets web application type
		parts := strings.Split(r.URL.Path, "/")						            // Splits teh url into sections or parts using /
		if len(parts) == 6 || parts[1] != "conservation"{                       // If the url has too few parts or too many, an error handler is used
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return                                                              // Error handler, returns the bad error request
		}
		if parts[4] == "" {													    // If part 4 is empty reply with all species stored using the function
			replyWithAllSpecies(w, &DBs)
		} else {																// Otherwise print one specific species
			replyWithSpecies(w, &DBs, parts[4])
		}
}
