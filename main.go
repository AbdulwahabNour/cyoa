package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Story map[string]Chapter

var temp *template.Template
 
func main() {
  
     flaghttp := flag.Bool("http", false, "Run as web server")
     flagStoryFilename := flag.String("story", "gopher.json", "The path to Json file for the story to render default `gopher.json` ")
     flag.Parse()
     filename := "arc.txt"

     story := storyDecoder(*flagStoryFilename)
     
     if *flaghttp{
          temp = template.Must(template.ParseFiles("cyoa.html"))
          webServerHandler(filename, &story)
          return
     }

     temp = template.Must(template.ParseFiles("cyoa.txt"))
   
     cmdHandler(story)
}

type StoryHandler struct {
	story    Story
	filename string
}

func (s StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]
	temp.Execute(w, s.story[path])
}
func webServerHandler(flagStoryFilename string, story *Story){
	
	s := StoryHandler{story: *story, filename: flagStoryFilename}
	http.ListenAndServe(":8080", s)
}
func cmdHandler(s Story){
     arc := s["intro"]
     var choice int
     for{
          arcOptionsLen := len(arc.Options) 
          temp.Execute(os.Stdout, arc)
          if arcOptionsLen == 0{
               break
          }
          fmt.Print("Choice: ")
          if _, err  := fmt.Scanf("%d\n", &choice); err !=nil{
               log.Fatalf("Failed to scan %v", err)
          }
          if  choice >= arcOptionsLen || choice < 0{
               log.Fatalf("Invalid choice (%d)", choice)
          }
          arc = s[arc.Options[choice].Arc]
    

     }   
}
func storyDecoder(filename string) Story{
     f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open %s: %v", filename, err)
	}
	defer f.Close()
     
	var story Story
	err = json.NewDecoder(f).Decode(&story)
	if err != nil {
		log.Fatalf("Failed to decode %v", err)
	}
     return story
}