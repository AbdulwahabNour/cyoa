package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Option struct{
    Text string `json:"text"`
    Arc string `json:"arc"`
}

type Chapter struct {
    Title   string   `json:"title"`
    Story   []string `json:"story"`
    Options []Option  `json:"options"`
}  

type Story map[string]Chapter

var temp *template.Template
func init(){
     temp = template.Must(template.ParseFiles("cyoa.html"))   
}
func main(){
   flagStoryFilename := flag.String("story", "gopher.json", "The path to Json file for the story to render default `gopher.json` ")
   flag.Parse()
   f, err:= os.Open(*flagStoryFilename)
   if err != nil{
        log.Fatalf("Failed to open %s: %v", *flagStoryFilename, err)
   }
   defer f.Close()
   var story Story
   err  = json.NewDecoder(f).Decode(&story)
    if err != nil{
         log.Fatalf("Failed to decode %v",  err)
    }
    s := StoryHandler{story:story, filename: *flagStoryFilename}
    http.ListenAndServe(":8080", s)

}

type StoryHandler struct{
     story Story
     filename string
}
func (s StoryHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
     path := strings.TrimSpace(r.URL.Path)

     if path == "" || path =="/"{
          path ="/intro"
     }
     path = path[1:]
     temp.Execute(w, s.story[path])
}