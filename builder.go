package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	gosass "github.com/bep/godartsass"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"

	xnhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var htmPlugin = api.Plugin{
	Name: "htm",
	Setup: func(build api.PluginBuild) {
		build.OnResolve(
			api.OnResolveOptions{Filter: `^*.htm$`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Namespace: "htm-ns",
					Path:      strings.ReplaceAll(args.ResolveDir, "\\", "/") + "/" + args.Path,
				}, nil
			})

		build.OnLoad(api.OnLoadOptions{Filter: `.*`, Namespace: "htm-ns"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				htmPath := args.Path

				data, err := os.ReadFile(htmPath)
				if err != nil {
					return api.OnLoadResult{}, err
				}

				options := html.Minifier{
					KeepWhitespace:          false,
					KeepComments:            false,
					KeepConditionalComments: false,
					KeepEndTags:             false,
					KeepQuotes:              false,
				}

				m := minify.New()
				m.AddFunc("text/html", options.Minify)

				minified, err := m.Bytes("text/html", data)
				if err != nil {
					log.Println("minify error: " + err.Error())
					return api.OnLoadResult{}, err
				}

				re := regexp.MustCompile(`\s{2,}|\n`)
				minified = re.ReplaceAll(minified, []byte(""))

				contents := string(minified)
				return api.OnLoadResult{
					Contents: &contents,
					Loader:   api.LoaderText,
				}, nil
			})
	},
}

func buildJS() {
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{sourceRoot + "/" + "ts/main.ts"},
		Outfile:           distRoot + "/" + "app.js",
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Target:            api.ES2015,
		Write:             true,
		Plugins:           []api.Plugin{htmPlugin},
	})

	if len(result.Errors) > 0 {
		fmt.Println("BuildJS:", result.Errors[0].Text)
	}
}

func buildSCSS() {
	mainSCSS := "main.scss"
	includes := []string{sourceRoot + "/" + "scss", sourceRoot + "/" + "scss/consts"}

	output := distRoot + "/" + "style.css"
	createSourceMAP := true

	options := gosass.Options{DartSassEmbeddedFilename: "sass_embedded\\dart-sass-embedded.bat"}
	args := gosass.Args{
		IncludePaths:    includes,
		Source:          `@import "` + mainSCSS + `";`,
		OutputStyle:     gosass.OutputStyleCompressed,
		SourceSyntax:    gosass.SourceSyntaxSCSS,
		EnableSourceMap: createSourceMAP,
	}

	tp, err := gosass.Start(options)
	if err != nil {
		if Mode == "watch" {
			log.Println("SCSS Start: " + err.Error())
		} else {
			log.Fatalln("SCSS Start: " + err.Error())
		}
	}
	result, err := tp.Execute(args)
	if err != nil {
		if Mode == "watch" {
			log.Println("SCSS Execute: " + err.Error())
		} else {
			log.Fatalln("SCSS Execute: " + err.Error())
		}
	}

	err = os.WriteFile(output, []byte(result.CSS), os.ModeAppend|os.ModePerm)
	if err != nil {
		if Mode == "watch" {
			log.Println("SCSS WriteFile: " + err.Error())
		} else {
			log.Fatalln("SCSS WriteFile: " + err.Error())
		}
	}

	if createSourceMAP {
		err = os.WriteFile(output+".map", []byte(result.SourceMap), os.ModeAppend|os.ModePerm)
		if err != nil {
			if Mode == "watch" {
				log.Println("SCSS WriteFile: " + err.Error())
			} else {
				log.Fatalln("SCSS WriteFile: " + err.Error())
			}
		}
	}
}

func buildHTML() {
	source := sourceRoot + "/" + "html/index.html"
	output := distRoot + "/" + "index.html"

	// use later
	// f, _ := os.Stat(output)
	// if f.IsDir() {
	// 	output = strings.TrimSuffix(output, "/") + "/index.html"
	// }

	data, err := os.ReadFile(source)
	if err != nil {
		log.Fatalln("HTML ReadFile: " + err.Error())
	}

	node, err := xnhtml.Parse(bytes.NewReader(data))
	if err != nil {
		log.Fatalln("HTML Parse: " + err.Error())
	}

	childCSS := xnhtml.Node{
		Type:     xnhtml.ElementNode,
		Data:     "link",
		DataAtom: atom.Link,
		Attr: []xnhtml.Attribute{
			{Key: "rel", Val: "stylesheet"},
			{Key: "href", Val: "style.css"},
		},
	}
	node.LastChild.LastChild.AppendChild(&childCSS)

	childJS := xnhtml.Node{
		Type:     xnhtml.ElementNode,
		Data:     "script",
		DataAtom: atom.Script,
		Attr:     []xnhtml.Attribute{{Key: "src", Val: "app.js"}},
	}
	node.LastChild.LastChild.AppendChild(&childJS)

	sse, err := xnhtml.Parse(strings.NewReader(`<script>(() => new EventSource("/hello-esbuild-event").onmessage = () => location.reload())()</script>`))
	if err != nil {
		log.Fatalln("HTML Parse: " + err.Error())
	}
	if Mode == "watch" {
		node.LastChild.LastChild.AppendChild(sse)
	}

	var buf bytes.Buffer
	err = xnhtml.Render(&buf, node)
	if err != nil {
		log.Fatalln("HTML Render: " + err.Error())
	}

	options := html.Minifier{KeepDocumentTags: true}

	m := minify.New()
	m.AddFunc("text/html", options.Minify)

	minified, err := m.Bytes("text/html", buf.Bytes())
	if err != nil {
		log.Fatalln("HTML Minify: " + err.Error())
	}

	os.WriteFile(output, minified, os.ModeAppend|os.ModePerm)
}

func buildAll() {
	buildJS()
	buildSCSS()
	buildHTML()
}
