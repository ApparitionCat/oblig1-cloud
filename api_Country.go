package oblig1

import (
	"encoding/json"
	"net/http"
	"strings"
	"bytes"
)


type OccurenceList struct {									
	Oresult[] 					Occurence `json:"results"`
}

type Occurence struct {											
	CountryCode 					string `json:"countryCode"`
	GenericName						string `json:"genericName"`
	SpeciesKey						uint64 `json:"speciesKey"`
}



func handlerGetCountry(w http.ResponseWriter) {                 // Gets countries information
	url := "https://restcountries.eu/rest/v2/all"		        // url used
	countryList, err := getCountryJSON(url)				    
	if err != nil {												// If an error occurs, give the 503 error code
		http.Error(w, "Country Service Unavailable", http.StatusServiceUnavailable)
		DN.TestApi("country")										
	}

	for idx, row := range countryList {				            // For each country in country list
		if idx == 0 {											//this loop runs if idx is 0		
			continue
		}

		DBc.Add(row)														
	}
}

var OCCstructure = new(OccurenceList)				            // Structure for Occurence
var DBc = CountriesDB{}											// Stores countries
// Fetches crountries from url
func getCountryJSON(url string) ([]Country, error) {
	resp, err := http.Get(url)								  
	if err != nil {												// Gets url, returns error if something goes wrong
		return nil, err
	}

	defer resp.Body.Close()
	var countryList []Country
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()										
	if err := json.Unmarshal(respByte, &countryList); err != nil {
		return nil, err											// return error if something goes wrong
	}

	return countryList, nil										
}


//-------------------------------------------------------------




func replyWithAllCountries(w http.ResponseWriter, DB countryStorage) {

		if DB.Count() == 0 {									 //if there are no countries, print out empty json
			json.NewEncoder(w).Encode([]Country{})
		} else {																
			a := make([]Country, 0, DB.Count())		             // make map variable for printing
			for _, s := range DB.GetAll() {				         
				a = append(a, s)										
			}
			json.NewEncoder(w).Encode(a)					     // Display as JSON on browser
		}
	}
                                                                 // Reply with specific country with specified amount of occurrences
func handlerCountry(w http.ResponseWriter, r *http.Request) {
		handlerGetCountry(w)											// gets countries using called function
		http.Header.Add(w.Header(), "content-type", "application/json")	// Assigns content type
		parts := strings.Split(r.URL.Path, "/")		                    // Splits url into string variables
		var limit string = r.URL.Query().Get("limit")	                // Fetches the query limit
		if limit == ""{														
			limit = "10"											    // The limit is set to 10 by default if no limit was set
		}

		if len(parts) == 6 || parts[1] != "conservation" {	// Errorhandling for malformed url
			http.Error(w, "Bad request:", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {											// If 4th section does not contain counry code:part is empty it will reply with all countries
		} else {												    //otherwise it replies with one			
			replyWithCountry(w, &DBc, parts[4], limit)	        
		}
}

func replyWithCountry(w http.ResponseWriter, DB countryStorage, id string, limit string) {
	url := "http://api.gbif.org/v1/occurrence/search?country=" + id + "&limit=" + limit
	resp, err := http.Get(url)								     // get url above
	if err != nil {										         // Error handling, returns error code
		http.Error(w, "Occurrence Service Unavailable", http.StatusServiceUnavailable)
		DN.TestApi("occurrence")								 
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	json.Unmarshal(respByte, &OCCstructure)		                  // Copy JSON Into map

	for idx, x := range OCCstructure.Oresult{                     // For each occurrence loop is run
		if idx == 0 {														
			continue
		}
		DBc.AssignSpecies(x)									   // Assign the species (occurrence) to the specific country using function
	}

	s, ok := DB.Get(id)											    // Get country id
	if !ok {														// CHeck if country with id exists
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return                                                      //if its not found, error code 404
	}
		json.NewEncoder(w).Encode(s)							    // Print as JSON to display
}