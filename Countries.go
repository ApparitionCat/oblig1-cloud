package oblig1

// Unified way of accessing country data, essentially local functions
type countryStorage interface {
	Init()
	Add(c Country) error
	Count() int
	Get(key string) (Country, bool)
	GetAll() []Country
	AssignSpecies(occ Occurence)
}


type Country struct {
	Code              string `json:"alpha2Code"`
	CountryName       string `json:"name"`
	CountryFlag       string `json:"flag"`
	Species[]					string
	SpeciesKey[]			uint64
}

// Establish countries in memory
type CountriesDB struct {
	countries map[string]Country
}

// Initialise the country slice
func (db *CountriesDB) Init() {
	db.countries = make(map[string]Country)
}

// Add new country
func (db *CountriesDB) Add(c Country) error {
	db.countries[c.Code] = c
	return nil
}

// Return amount of countries stored
func (db *CountriesDB) Count() int {
	return len(db.countries)
}

// Return country with given code
func (db *CountriesDB) Get(codeID  string) (Country, bool) {
	s, ok := db.countries[codeID]
	return s, ok
}

// Return all countries as a slice
func (db *CountriesDB) GetAll() []Country {
	all := make([]Country, 0, db.Count())
	for _, s := range db.countries {
		all = append(all, s)
	}
	return all
}

// See if string already exists in array of strings
func stringExists(a string, list []string) bool{
	for _, b := range list {
		if b == a{
			return true
		}
	}
	return false
}

// Assigns generic name and species key to country struct from occurence
func (db *CountriesDB) AssignSpecies(occ Occurence) {
	var tempStorage = db.countries[occ.CountryCode];
	if !stringExists(occ.GenericName, tempStorage.Species) {
		tempStorage.Species = append(tempStorage.Species, occ.GenericName)
		tempStorage.SpeciesKey = append(tempStorage.SpeciesKey, occ.SpeciesKey)
	}
	db.countries[occ.CountryCode] = tempStorage
}
