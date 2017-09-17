package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// Liste de personnes
var gens []Personne

// Personne objet représentant une personne
type Personne struct {
	ID       string    `json:"id,omitempty"`
	Prenom   string    `json:"prenom,omitempty"`
	Nom      string    `json:"nom,omitempty"`
	Addresse *Addresse `json:"addresse,omitempty"`
}

// Addresse objet représentant l'Adresse d'un personne
type Addresse struct {
	Ville string `json:"ville,omitempty"`
	Pays  string `json:"pays,omitempty"`
}

// GetPersonneEndpoint endpoint pour retourner une personne identifiée par un Id
func GetPersonneEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range gens {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Personne{})
}

// GetGensEndpoint endpoint pour retourner la liste des personnes
func GetGensEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(gens)
}

// CreatePersonneEndpoint endpoint pour ajouter une nouvelle personne à la liste
func CreatePersonneEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var personne Personne
	_ = json.NewDecoder(req.Body).Decode(&personne)
	personne.ID = params["id"]
	gens = append(gens, personne)
	json.NewEncoder(w).Encode(gens)
}

// DeletePersonneEndpoint endpoint pour supprimer une personne de la liste
func DeletePersonneEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range gens {
		if item.ID == params["id"] {
			gens = append(gens[:index], gens[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(gens)
}

// Pays object décrivant un pays
type Pays struct {
	Name       string `xml:"name" json:"name"`
	Population string `xml:"population" json:"population"`
	Capital    string `xml:"capital" json:"capital"`
	Currency   string `xml:"currency" json:"currency"`
}

// GetCountryEndpoint endpoint pour retourner un pays
func GetCountryEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	soapResponse := queryCountry(params["id"])

	jsonData, err := json.Marshal(soapResponse.Body)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	// json.NewEncoder(w).Encode(soapResponse)
}

func queryCountry(country string) *http.Response {
	url := "http://localhost:8080/ws"
	client := &http.Client{}
	sRequestContent := generateRequestContent(country)
	// fmt.Printf(sRequestContent)
	requestContent := []byte(sRequestContent)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestContent))
	if err != nil {
		println(&err)
	}

	// req.Header.Add("SOAPAction", `" golang"`)
	req.Header.Add("Content-Type", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		println(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		println("Error Response " + resp.Status)
	}

	// print de la réponse
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	fmt.Printf(newStr)

	return resp
}

// generateRequestContent générer une requête SOAP au format String
func generateRequestContent(country string) string {
	type QueryData struct {
		Country string
	}
	const getTemplate = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
				  xmlns:gs="http://spring.io/guides/gs-producing-web-service">
   <soapenv:Header/>	
   <soapenv:Body>
      <gs:getCountryRequest>
         <gs:name>{{.Country}}</gs:name>
      </gs:getCountryRequest>
   </soapenv:Body>
</soapenv:Envelope>`
	querydata := QueryData{Country: country}
	tmpl, err := template.New("getCoutry").Parse(getTemplate)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, querydata)
	if err != nil {
		panic(err)
	}
	return doc.String()
}

func main() {
	router := mux.NewRouter()

	// initialisation des personnes de la liste
	gens = append(gens, Personne{ID: "1", Prenom: "Ryan", Nom: "Gosling", Addresse: &Addresse{Ville: "Los Angeles", Pays: "Etats-Unis"}})
	gens = append(gens, Personne{ID: "2", Prenom: "Sébastien", Nom: "Barbara"})
	gens = append(gens, Personne{ID: "3", Prenom: "Lucas", Nom: "Bowler"})
	gens = append(gens, Personne{ID: "4", Prenom: "Cyril", Nom: "Constant"})
	gens = append(gens, Personne{ID: "5", Prenom: "Alex", Nom: "Binguy"})
	gens = append(gens, Personne{ID: "6", Prenom: "Axel", Nom: "Prieur", Addresse: &Addresse{Ville: "Paris", Pays: "France"}})
	gens = append(gens, Personne{ID: "7", Prenom: "Pascal", Nom: "Spadone", Addresse: &Addresse{Ville: "Nice", Pays: "France"}})

	// routage des appels vers les endpoints correspondants
	router.HandleFunc("/personnes", GetGensEndpoint).Methods("GET")
	router.HandleFunc("/personnes/{id}", GetPersonneEndpoint).Methods("GET")
	router.HandleFunc("/personnes/{id}", CreatePersonneEndpoint).Methods("POST")
	router.HandleFunc("/personnes/{id}", DeletePersonneEndpoint).Methods("DELETE")

	router.HandleFunc("/country/{id}", GetCountryEndpoint).Methods("GET")

	// démarrage du server
	log.Fatal(http.ListenAndServe(":12345", router))

}
