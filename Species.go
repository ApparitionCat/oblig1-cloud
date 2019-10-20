package oblig1

type Species struct {
	Key       uint64 `json:"speciesKey"`
	Kingdom   string `json:"kingdom"`
	Phylum    string `json:"phylum"`
	Order     string `json:"order"`
	Family    string `json:"family"`
	Genus     string `json:"genus"`
	SciName   string `json:"sciName"`
	CanName   string `json:"canName"`
	Year      int    `json:"year"`
}


type SpeciesStorage interface {
	Init()
	Add(s Species) error
	Count() int
	Get(key uint64) (Species, bool)
	GetAll() []Species
}
 ///////////////////////////////////////////
type ResultList struct {
	AllSpecies[] 					Species `json:"results"`
}
////////////////////////////////////////////////

//	Species struct created in a memory
type SpeciesDB struct {
	species map[uint64]Species
}


func (db *SpeciesDB) Init() {
	db.species = make(map[uint64]Species)
}

// Adds a new species with new id
func (db *SpeciesDB) Add(s Species) error {
	db.species[s.Key] = s
	return nil
}

// Returns the number of species stored
func (db *SpeciesDB) Count() int {
	return len(db.species)
}

// returns a species
func (db *SpeciesDB) Get(Key uint64) (Species, bool) {
	s, ok := db.species[Key]
	return s, ok
}

// Fget all species
func (db *SpeciesDB) GetAll() []Species {
	all := make([]Species, 0, db.Count())
	for _, s := range db.species {
		all = append(all, s)
	}
	return all
}
