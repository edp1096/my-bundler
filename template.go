package main

import (
	"errors"
	"os"

	"io/fs"
)

func exportOrPass(fname string, fdata []byte) {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		os.WriteFile(fname, fdata, 0644)
	}
}

func exportTemplate() error {
	var err error

	exportOrPass("package.json", pkgJSON)
	exportOrPass("tsconfig.json", tsconfJSON)
	exportOrPass("types.d.ts", typesDts)

	if _, err = os.Stat("src"); !os.IsNotExist(err) {
		return errors.New("directory src already exist")
	}

	root := embedSRC
	prefix := "embed/"

	err = fs.WalkDir(root, "embed/src", func(path string, d fs.DirEntry, err error) error {
		outname := path[len(prefix):]
		if d.IsDir() {
			os.MkdirAll(outname, 0755)
			return nil
		}

		fdata, er := fs.ReadFile(root, path)
		if err != nil {
			return er
		}

		out, er := os.Create(outname)
		if err != nil {
			return er
		}
		defer out.Close()

		_, er = out.Write(fdata)
		if err != nil {
			return er
		}

		return nil
	})

	return err
}
