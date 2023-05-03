package controllers

import(
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"golangmovietask/models"
	"html/template"

)

var templ *template.Template

func (m *Controllers) APIDocumentation(w http.ResponseWriter, r *http.Request){
	urlPath := r.URL.Path
	resp := strings.Split(urlPath, "/")
	title := resp[1]
	content,err := ioutil.ReadFile("APIDocumentation.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}	
	endspointsdata := models.EndpointsHead{}
	json.Unmarshal(content,&endspointsdata)
	s:=models.DocumentationParseData{Title: title,Endpointsdata: endspointsdata}
	templ = template.Must(template.ParseFiles("views/index.html"))
	err = templ.Execute(w,s)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}