package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gens)
}

// CreatePersonneEndpoint endpoint pour ajouter une nouvelle personne à la liste
func CreatePersonneEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var personne Personne
	json.NewDecoder(req.Body).Decode(&personne)
	personne.ID = params["id"]
	gens = append(gens, personne)
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gens)
}

// QueryData argument du template SOAP
type QueryData struct {
	Country string
}

// MyRespEnvelope enveloppe de la réponse
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    Body
	Header  Header
}

// Body body de la réponse
type Body struct {
	XMLName     xml.Name
	GetResponse completeResponse `xml:"getCountryResponse" json:"getCountryResponse"`
}

// Header header de la réponse
type Header struct {
	XMLName xml.Name
}

// completeResponse object pays
type completeResponse struct {
	XMLName xml.Name `xml:"getCountryResponse"`
	Country Pays     `xml:"country" json:"country"`
}

// Pays objet décrivant un pays
type Pays struct {
	Name       string `xml:"name" json:"name"`
	Population string `xml:"population" json:"population"`
	Capital    string `xml:"capital" json:"capital"`
	Currency   string `xml:"currency" json:"currency"`
}

// GetCountryEndpoint endpoint pour retourner un pays
func GetCountryEndpoint(w http.ResponseWriter, req *http.Request) {
	// récupération des paramèters de l'URL
	params := mux.Vars(req)

	// appel et récupération depuis le WS SOAP
	queryCountry(params["id"], w)
}

func queryCountry(country string, w http.ResponseWriter) {

	// création de la requête du WS SOAP
	url := "http://localhost:8080/ws"
	client := &http.Client{}
	stringRequestContent := generateRequestContent(country)
	bytesRequestContent := []byte(stringRequestContent)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(bytesRequestContent))
	if err != nil {
		println(&err)
	}
	request.Header.Add("Content-Type", "text/xml")

	// récupération de la réponse
	response, err := client.Do(request)
	if err != nil {
		println(err)
	}
	if response.StatusCode != 200 {
		println("Error Response " + response.Status)
	}

	// http.response -> string
	bufferResponse := new(bytes.Buffer)
	bufferResponse.ReadFrom(response.Body)
	stringResponse := bufferResponse.String()
	fmt.Println(stringResponse)

	// string -> struct
	enveloppe := &MyRespEnvelope{}
	var bytesResponse = []byte(stringResponse)
	err = xml.Unmarshal(bytesResponse, enveloppe)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(enveloppe)
	}

	// écriture de la réponse en json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(enveloppe.Body.GetResponse.Country)

	defer response.Body.Close()
}

// generateRequestContent générer une requête SOAP au format String
func generateRequestContent(country string) string {

	// template de l'enveloppe SOAP
	const getTemplate = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:gs="http://spring.io/guides/gs-producing-web-service">
	<soapenv:Header/><soapenv:Body><gs:getCountryRequest><gs:name>{{.Country}}</gs:name></gs:getCountryRequest></soapenv:Body></soapenv:Envelope>`

	// changement de paramètre en fonction de ce qui a été envoyé dans l'URL
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
	gens = append(gens, Personne{ID: "2", Prenom: "Sébastien", Nom: "B"})
	gens = append(gens, Personne{ID: "3", Prenom: "Lucas", Nom: "B"})
	gens = append(gens, Personne{ID: "4", Prenom: "Cyril", Nom: "C"})
	gens = append(gens, Personne{ID: "5", Prenom: "Alex", Nom: "B"})
	gens = append(gens, Personne{ID: "6", Prenom: "Axel", Nom: "P", Addresse: &Addresse{Ville: "Paris", Pays: "France"}})
	gens = append(gens, Personne{ID: "7", Prenom: "Pascal", Nom: "S", Addresse: &Addresse{Ville: "Nice", Pays: "France"}})

	// routage des appels vers les endpoints correspondants
	router.HandleFunc("/personnes", GetGensEndpoint).Methods("GET")
	router.HandleFunc("/personnes/{id}", GetPersonneEndpoint).Methods("GET")
	router.HandleFunc("/personnes/{id}", CreatePersonneEndpoint).Methods("POST")
	router.HandleFunc("/personnes/{id}", DeletePersonneEndpoint).Methods("DELETE")
	router.HandleFunc("/country/{id}", GetCountryEndpoint).Methods("GET")

	// démarrage du server
	log.Fatal(http.ListenAndServe(":12345", router))
}
