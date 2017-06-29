package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"encoding/json"
	//"io"
	"io/ioutil"
	"bytes"

	"github.com/manifoldco/go-signature"

	//"github.com/goji/param"
	"github.com/zenazn/goji"
    "github.com/zenazn/goji/web"
    //"github.com/zenazn/goji/web/middleware"
)

// Grab env variables
var Master_key string = os.Getenv("MASTER_KEY")
var Client_id string = os.Getenv("CLIENT_ID")
var Client_secret = os.Getenv("CLIENT_SECRET")
var Connector_url = os.Getenv("CONNECTOR_URL")

var Plans = []string{"small", "large"}
var Products = []string{"numbers"}
var Regions = []string{"aws::us-east-1"}

var Resources = make(map[string]string)
var Credentials = make(map[string]string)


// Some helpful Structs
type ResponseBody struct {
	Message 	string 	`json:"message"`
	Credentials string 	`json:"credentials"`
}

type RequestBody struct {
	Id 		string `json:"id"`
	Product string `json:"product"`
	Plan 	string `json:"plan"`
	Region 	string `jsson:"region"`
}

type CredRequestBody struct {
	Id 			string `json:"id"`
	ResourceId	string `json:"resource_id"`
}

func main() {

	goji.Put("/v1/resources/:id", PutResources)

	goji.Patch("/v1/resources/:id", PatchResources)

	goji.Delete("/v1/resources/:id", DeleteResources)

	goji.Put("/v1/credentials/:id", PutCredentials)

	goji.Delete("/v1/credentials/:id", DeleteCredentials)

	goji.Get("/v1/sso", GetSso)

	goji.Serve()
}

func PutResources(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("PutResources")

	//VerifySignature(w, r)
	/*
	body, _ := ioutil.ReadAll(r.Body)
	buf := bytes.NewBuffer(body)

	verifier, _ := signature.NewVerifier(signature.ManifoldKey)
	if err := verifier.Verify(r, buf); err != nil {
		// return an error...
		log.Print(err)
		SendResponse(w, http.StatusUnauthorized, "Invalid Signature")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Signature verified")*/
	
	id := c.URLParams["id"]

    body := ParseBody(r)

    if !IsInArray(body.Plan, Plans) {
    	SendResponse(w, http.StatusBadRequest, "Bad Plan")
    }
    if !IsInArray(body.Region, Regions) {
    	SendResponse(w, http.StatusBadRequest, "Bad Region")
    }
    if !IsInArray(body.Product, Products) {
    	SendResponse(w, http.StatusBadRequest, "Bad Product")
    }

	existing := Resources[id]

	bodystring, err := json.Marshal(body)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
  	}

	if existing != "" && string(bodystring) != existing {
		SendResponse(w, http.StatusConflict, "Resource already exists")
		return
	} 
	
	Resources[id] = string(bodystring)

	SendResponse(w, http.StatusCreated, "Resource has been created")
}

func PatchResources(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("PatchResources")

	body := ParseBody(r)
	
	plan := body.Plan
	if !IsInArray(plan, Plans) {
		// bad plan
		SendResponse(w, http.StatusBadRequest, "Bad Plan")
		return
	}

	id := c.URLParams["id"]
	existing := Resources[id]
	if existing == ""{
		SendResponse(w, http.StatusNotFound, "Resource could not be found")
		return
	}

	SendResponse(w, http.StatusOK, "Plan has been changed")
}

func DeleteResources(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteResources")
	id := c.URLParams["id"]
	existing := Resources[id]
	if existing == ""{
		SendResponse(w, http.StatusNotFound, "Resource could not be found")
		return
	}

	delete(Resources, id)
	SendResponse(w, http.StatusNoContent, "Resource deleted")
}

func PutCredentials(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("PutCredentials")
	id := c.URLParams["id"]

	body := ParseCredBody(r)

	existing := Resources[body.ResourceId]
	if existing == "" {
		SendResponse(w, http.StatusNotFound, "Resource could not be found")
		return
	}

	bodystring, err := json.Marshal(body)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
  	}

  	fmt.Println("Creating Credentials")

  	Credentials[id] = string(bodystring)

  	SendResponse(w, http.StatusCreated, "Credentials created")
}

func DeleteCredentials(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteCredentials")
	id := c.URLParams["id"]

	existing := Credentials[id]
	if existing == ""{
		SendResponse(w, http.StatusNotFound, "Resource could not be found")
		return
	}

	delete(Credentials, id)
	SendResponse(w, http.StatusNoContent, "Resource deleted")
}

func GetSso(w http.ResponseWriter, r *http.Request) {
	log.Print("GetSso")
}

// Helper functions
// TODO - move these to anyother file/package
func VerifySignature(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	buf := bytes.NewBuffer(body)

	verifier, _ := signature.NewVerifier(Master_key)
	if err := verifier.Verify(r, buf); err != nil {
		// return an error...
		SendResponse(w, http.StatusUnauthorized, err.Error())
    	return
	}
}

func SendResponse(w http.ResponseWriter, code int, message string) {

	msg := &ResponseBody{Message: message}

	js, err := json.Marshal(msg)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
  	}
  	//Log for testing
  	//fmt.Println(string(js))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func SendCredResponse(w http.ResponseWriter, code int, message string, credentials string) {
	
	msg := &ResponseBody{Message: message, Credentials: credentials}

	js, err := json.Marshal(msg)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
  	}

  	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func ParseBody(r *http.Request) RequestBody {
	decoder := json.NewDecoder(r.Body)
    var body RequestBody   
    err := decoder.Decode(&body)
    if err != nil {
        panic(err)
    }

    return body
}

func ParseCredBody(r *http.Request) CredRequestBody {
	decoder := json.NewDecoder(r.Body)
    var body CredRequestBody   
    err := decoder.Decode(&body)
    if err != nil {
        panic(err)
    }

    return body
}

func IsInArray(a string, list []string) bool {
    for _, b := range list {
    	// Log for testing
    	//fmt.Println ("a - ", a, "b - ", b)
        if b == a {
            return true
        }
    }
    return false
}
// End of helper functions
