package dl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
)

// WriteCounter : Progress counter
type WriteCounter struct {
	Total uint64
}

// Write : Use as io.Writer for io.TeeReader
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress : Print progress to console
func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// DownloadFile : Download file
func DownloadFile(filepath string, url string) error {
	// Create filename + tmp
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	// Receive data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Create & pass progress
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// Below closes lock the file so, defer should not use
	resp.Body.Close()
	out.Close()

	fmt.Print("\n")

	// Remove tmp
	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}
