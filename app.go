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
	Message 	string 				`json:"message"`
	Credentials map[string]string 	`json:"credentials"`
}

type RequestBody struct {
	Id 		string `json:"id"`
	Product string `json:"product"`
	Plan 	string `json:"plan"`
	Region 	string `jsson:"region"`
}

func main() {

	//Resources = make(map[string]int)
	//Credentials = make(map[string]int)

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
	/*body, _ := ioutil.ReadAll(r.Body)
	buf := bytes.NewBuffer(body)

	verifier, _ := signature.NewVerifier(signature.ManifoldKey)
	if err := verifier.Verify(r, buf); err != nil {
		// return an error...
		log.Print(err)
		SendResponse(w, http.StatusUnauthorized, "Invalid Signature")
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Signature verified")
	*/
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
	SendResponse(w, http.StatusNoContent, "Resource gone");
}

func PutCredentials(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("PutCredentials")
	//id := c.URLParams["id"]
	//w.Write(id)
}

func DeleteCredentials(c web.C, w http.ResponseWriter, r *http.Request) {
	log.Print("DeleteCredentials")
	//id := c.URLParams["id"]
	//w.Write(id)
}

func GetSso(w http.ResponseWriter, r *http.Request) {
	log.Print("GetSso")
}

// Helper functions
// TODO - move these to anyother file/package
func VerifySignature(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	buf := bytes.NewBuffer(body)

	verifier, _ := signature.NewVerifier(signature.ManifoldKey)
	if err := verifier.Verify(r, buf); err != nil {
		// return an error...
		http.Error(w, err.Error(), http.StatusUnauthorized)
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
  	fmt.Println(string(js))

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

func IsInArray(a string, list []string) bool {
    for _, b := range list {
    	fmt.Println ("a - ", a, "b - ", b)
        if b == a {
            return true
        }
    }
    return false
}
