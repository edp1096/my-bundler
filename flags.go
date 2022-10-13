package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func procFlags() {
	FlagV := flag.Bool("v", false, " Print version")
	FlagVersion := flag.CommandLine.Bool("version", false, " Print version")
	FlagB := flag.Bool("b", false, " Build files")
	FlagBuild := flag.CommandLine.Bool("build", false, " Build files")
	FlagW := flag.Bool("w", false, " Watch files")
	FlagWatch := flag.CommandLine.Bool("watch", false, " Watch files")
	FlagS := flag.String("s", "[src]", " Set source root")
	FlagSrc := flag.CommandLine.String("source", "[src]", " Set source root")
	FlagD := flag.String("d", "[dist]", " Set distribution root")
	FlagDist := flag.CommandLine.String("dist", "[dist]", " Set distribution root")
	FlagT := flag.Bool("t", false, " Get template src")
	FlagTemplate := flag.CommandLine.Bool("template", false, " Get template src")

	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("Usage: %s [-b|-w] [-s path]\n", "builder")
		fmt.Printf("  %-24s Build from ./src\n", " no option")

		order := []string{"v", "b", "w", "t", "s", "d"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			switch name {
			case "v":
				fmt.Printf("  -%-23s%s\n", "v, --version", flag.Usage)
			case "b":
				fmt.Printf("  -%-23s%s\n", "b, --build", flag.Usage)
			case "w":
				fmt.Printf("  -%-23s%s\n", "w, --watch", flag.Usage)
			case "t":
				fmt.Printf("  -%-23s%s\n", "t, --template", flag.Usage)
			case "s":
				fmt.Printf("  -%-20s%s\n", "s"+" "+flag.Value.String()+", --source"+" "+flag.Value.String(), flag.Usage)
			case "d":
				fmt.Printf("  -%-20s%s\n", "d"+" "+flag.Value.String()+", --dist"+" "+flag.Value.String(), flag.Usage)
			}
		}
	}

	flag.Parse()

	if *FlagV || *FlagVersion {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}

	if *FlagBuild && *FlagWatch {
		log.Fatal("Cannot use --build and --watch together")
	}

	switch {
	case *FlagB || *FlagBuild:
		Mode = "build"
	case *FlagW || *FlagWatch:
		Mode = "watch"
	case *FlagT || *FlagTemplate:
		Mode = "template"
	}

	if *FlagS != "[src]" {
		sourceRoot = *FlagS
	}
	if *FlagSrc != "[src]" {
		sourceRoot = *FlagSrc
	}
	if *FlagD != "[dist]" {
		distRoot = *FlagD
	}
	if *FlagDist != "[dist]" {
		distRoot = *FlagDist
	}
}
