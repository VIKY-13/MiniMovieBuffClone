package main

import (
	// "database/sql"
	"encoding/json"
	// "fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	// "os"
	"strings"

	// "github.com/gorilla/mux"
	// "github.com/joho/godotenv"
	_ "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func APIDocumentation(w http.ResponseWriter, r *http.Request){
	urlPath := r.URL.Path
	resp := strings.Split(urlPath, "/")
	title := resp[1]
	content,err := ioutil.ReadFile("APIDocumentation.json")
	checkErr(err)
	endspointsdata := []EndpointDescriptions{}
	json.Unmarshal(content,&endspointsdata)
	s:=documentationparsedata{Title: title,Endpointsdata: endspointsdata}
	templ= template.Must(template.ParseFiles("index.html"))
	err = templ.Execute(w,s)
	checkErr(err)
}