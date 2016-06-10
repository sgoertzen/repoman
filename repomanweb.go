package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"html/template"
	"log"
	"net/http"
)

type config struct {
	githubkey    *string
	organization *string
}


var _config config

// Program to read in poms and determine
func main() {

	_config = getConfiguration()
	runServer()

	log.Println("Done")
}

func showRepos(w http.ResponseWriter, r *http.Request) {
	// TODO: Include contents of log file on the main page
	repos := getAllRepos(*_config.organization)
	data := map[string]interface{}{
		"Repos": repos,
	}
	showTemplatedFile(w, "html/default.html", data)
}

func runServer() {
	http.HandleFunc("/", handler)

	// TODO: Allow port to be passed in
	log.Println("Webserver ready")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]

	switch path {
	case "about":
		http.NotFound(w, r)
	case "":
		showRepos(w, r)
	default:
		http.NotFound(w, r)
	}
}

// TODO: Read these parameters from either command line or config file
// Can look at calendar pi on how to read configuration from file
func getConfiguration() config {
	config := config{}
	config.organization = kingpin.Arg("org", "Github organization to analyze for upgrades").Required().String()
	config.githubkey = kingpin.Arg("githubkey", "Api key for github.").Required().String()
	kingpin.Version("1.0.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('?')
	kingpin.Parse()
	return config
}

func showTemplatedFile(w http.ResponseWriter, filename string, data map[string]interface{}) {

	t, err := template.ParseFiles(filename)
	if err != nil {
		log.Println("Error while parsing template file", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Println("Error while showing list ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
