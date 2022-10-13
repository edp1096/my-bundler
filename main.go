package main // import "builder"

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
)

var Version = "development"

var Mode = "build"

var sourceRoot = "src"
var distRoot = "dist"

var addr = "localhost:8864"

//go:embed embed/package.json
var pkgJSON []byte

//go:embed embed/tsconfig.json
var tsconfJSON []byte

//go:embed embed/types.d.ts
var typesDts []byte

//go:embed embed/src
var embedSRC embed.FS

var chReload chan bool

func handlerSSE(w http.ResponseWriter, r *http.Request) {
	if <-chReload {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("data: reload\n\n"))
	}
}

func checkAndGetSass() {
	err := checkSassExists()
	if err != nil {
		fmt.Println("checkSassExists:", err)
		getSass()
		os.RemoveAll("sass_embedded.zip")
	}
}

func main() {
	var err error

	if len(os.Args) > 1 {
		procFlags()
	}

	switch Mode {
	case "build":
		fmt.Println("Build..")

		checkAndGetSass()
		buildAll()

	case "watch":
		checkAndGetSass()
		buildAll()

		chFinish = make(chan bool)
		fmt.Println("Watching on " + addr)
		go watch()

		chReload = make(chan bool)

		http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, distRoot+"/index.html")
		})
		http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, distRoot+"/index.html")
		})
		http.Handle("/", http.FileServer(http.Dir(distRoot)))
		http.HandleFunc("/hello-esbuild-event", handlerSSE)
		http.ListenAndServe(addr, nil)

		<-chFinish
	case "template":
		fmt.Println("Generating template..")
		err = exportTemplate()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Done")

	default:
		log.Fatal("Invalid mode")
	}
}
