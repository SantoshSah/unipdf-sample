package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func init() {
	// To get your free API key for metered license, sign up on: https://cloud.unidoc.io
	// Make sure to be using UniPDF v3.19.1 or newer for Metered API key support.
	err := license.SetMeteredKey(`metered license key`)
	if err != nil {
		fmt.Printf("ERROR: Failed to set metered key: %v\n", err)
		fmt.Printf("Make sure to get a valid key from https://cloud.unidoc.io\n")
		panic(err)
	}
}

func main() {
	// Setting up meter licensed key
	lk := license.GetLicenseKey()
	if lk == nil {
		fmt.Printf("Failed retrieving license key")
		return
	}
	fmt.Printf("License: %s\n", lk.ToString())

	// GetMeteredState freshly checks the state, contacting the licensing server.
	state, err := license.GetMeteredState()
	if err != nil {
		fmt.Printf("ERROR getting metered state: %+v\n", err)
		panic(err)
	}
	fmt.Printf("State: %+v\n", state)
	if state.OK {
		fmt.Printf("State is OK\n")
	} else {
		fmt.Printf("State is not OK\n")
	}

	// Read pdf and extract policy terms
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == ".pdf" {
			inPath := f.Name()

			f, err := os.Open(f.Name())
			if err != nil {
				panic(err)
			}
			defer f.Close()

			pdfReader, err := pdf.NewPdfReaderLazy(f)
			if err != nil {
				fmt.Errorf("NewPdfReaderLazy failed. %q err=%v", inPath, err)
			}

			p, err := pdfReader.GetPage(1)
			if err != nil {
				fmt.Errorf("GetNumPages failed. %q  err=%v", inPath, err)
			}
			ex, err := extractor.New(p)
			if err != nil {
				fmt.Errorf("NewPdfReaderLazy failed. %q err=%v", inPath, err)
			}

			pageText, _, _, err := ex.ExtractPageText()
			if err != nil {
				fmt.Errorf("ExtractPageText failed. %q  err=%v", inPath, err)
			}

			text := pageText.Text()

			// Read start date
			fmt.Println("Start Date:", strings.ReplaceAll(between(text, "EFFECTIVE DATE", "EXPIRATION DATE"), "\n", ""))

			endDateText := after(text, "EXPIRATION DATE")

			// Read end date
			fmt.Println("End Date:", strings.ReplaceAll(endDateText[0:9], "\n", ""))
		}
	}
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
