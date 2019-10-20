package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"oblig1"
)

// This is simply the default message. Its shown when nothing is typed after /

func handlerNil(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}


func main() {
// struct map initialization, so that functions from structs in other files can be called

      oblig1.DBc.Init()
      oblig1.DBs.Init()
      oblig1.DN.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	} //these run functions based in the url/whats written after /
	http.HandleFunc("/", handlerNil)
	http.HandleFunc("/conservation/v1/country", oblig1.HandlerCountry)
	http.HandleFunc("/conservation/v1/species", oblig1.HandlerSpecies)
	http.HandleFunc("/conservation/v1/diag", oblig1.HandlerDiag)
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
