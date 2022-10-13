package arc

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func DecompressZIP(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	destFull := dest
	if dest == "" || dest == "." || dest == "."+string(os.PathSeparator) || dest == "./" {
		dest = ""
		destFull, _ = os.Getwd()
	}

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if dest != "" && !strings.HasPrefix(path, filepath.Clean(destFull)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: " + path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Fatalln(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func DecompressTGZ(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		path := filepath.Join(dest, header.Name)
		info := header.FileInfo()

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(path, info.Mode())
			continue
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(path), info.Mode())
			f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, tr)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown type: " + string(header.Typeflag))
		}
	}

	return nil
}
