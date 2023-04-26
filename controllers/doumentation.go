package controllers

import(
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"golangmovietask/dtos"
	"html/template"

)

var templ *template.Template

func APIDocumentation(w http.ResponseWriter, r *http.Request){
	urlPath := r.URL.Path
	resp := strings.Split(urlPath, "/")
	title := resp[1]
	content,err := ioutil.ReadFile("APIDocumentation.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}	
	endspointsdata := dtos.EndpointsHead{}
	json.Unmarshal(content,&endspointsdata)
	s:=dtos.DocumentationParseData{Title: title,Endpointsdata: endspointsdata}
	templ = template.Must(template.ParseFiles("index.html"))
	err = templ.Execute(w,s)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}