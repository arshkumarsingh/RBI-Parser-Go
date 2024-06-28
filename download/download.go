package download

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var bankNames map[string]string


// init initializes the bankNames map by loading bank names from a JSON file.
// The JSON file is expected to be named "banks.json" and located in the same directory as this file.
// The map is loaded using the encoding/json package.
// If the file cannot be opened or decoded, the program will log a fatal error and terminate.
func init() {
	// Open the JSON file
	file, err := os.Open("banks.json")
	if err != nil {
		// Log a fatal error and terminate the program if the file cannot be opened
		log.Fatalf("Could not open banks.json: %v", err)
	}
	defer file.Close()

	// Decode the JSON file into the bankNames map
	err = json.NewDecoder(file).Decode(&bankNames)
	if err != nil {
		// Log a fatal error and terminate the program if the file cannot be decoded
		log.Fatalf("Could not decode banks.json: %v", err)
	}
}


// LoadEtags loads the etags from a file
//
// Parameters:
// - filename: The name of the file to load the etags from.
//
// Returns:
// - map[string]string: The loaded etags.
// - error: An error if the file cannot be opened or decoded.
func LoadEtags(filename string) (map[string]string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a map to store the loaded etags
	etags := make(map[string]string)

	// Decode the JSON file into the etags map
	err = json.NewDecoder(file).Decode(&etags)
	if err != nil {
		return nil, err
	}

	return etags, nil
}

// SaveEtags saves the etags to a file
func SaveEtags(filename string, etags map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(etags)
}


// ExtractBankNameFromContext extracts the bank name from the surrounding context of the URL.
//
// Parameters:
// - s: The HTML selection containing the link to the bank's data.
//
// Returns:
// - string: The name of the bank.
// - error: An error if no matching bank name is found.
func ExtractBankNameFromContext(s *goquery.Selection) (string, error) {
	// Look for the closest preceding sibling that is a text node containing a bank name
	parent := s.Parent() // Get the parent element of the selection
	text := parent.Text() // Get the text of the parent element

	// Iterate over the bank names map
	for name := range bankNames {
		// Check if the name is present in the text
		matched, err := regexp.MatchString(`(?i)`+regexp.QuoteMeta(name), text)
		if err != nil {
			return "", err
		}
		if matched {
			return bankNames[name], nil // Return the corresponding bank name
		}
	}

	// If no match is found, try looking one level up in the DOM
	grandParent := parent.Parent() // Get the parent of the parent element
	text = grandParent.Text() // Get the text of the grandparent element

	// Iterate over the bank names map again
	for name := range bankNames {
		// Check if the name is present in the text
		matched, err := regexp.MatchString(`(?i)`+regexp.QuoteMeta(name), text)
		if err != nil {
			return "", err
		}
		if matched {
			return bankNames[name], nil // Return the corresponding bank name
		}
	}

	// If no matching bank name is found, return an error
	return "", fmt.Errorf("no matching bank name found in context: %s", text)
}

// DownloadFile downloads a file from a URL and saves it to the target path
func DownloadFile(url, target string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching file URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}

	return resp.Header.Get("etag"), nil
}

// DownloadAll downloads all .xlsx files from the RBI page and saves them to a directory
func DownloadAll(scrapeURL, xlsxDir, etagsFile string) {
	// Load the entity tags from the etags file or create a new map if the file cannot be opened
	etags, err := LoadEtags(etagsFile)
	if err != nil {
		etags = make(map[string]string)
	}

	// Fetch the HTML page from the RBI page
	resp, err := http.Get(scrapeURL)
	if err != nil {
		log.Fatalf("Error fetching URL %s: %v", scrapeURL, err)
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Find all the links to .xlsx files
	doc.Find("a[href$='.xlsx']").Each(func(i int, s *goquery.Selection) {
		// Add 1 to the WaitGroup counter
		wg.Add(1)

		// Start a new goroutine for each .xlsx file
		go func(s *goquery.Selection) {
			defer wg.Done()

			// Get the URL of the .xlsx file and check if the href attribute exists
			url, exists := s.Attr("href")
			if !exists {
				return
			}

			// Extract the bank name from the context of the .xlsx file
			bankName, err := ExtractBankNameFromContext(s)
			if err != nil {
				return
			}

			// Sanitize the bank name for use as a filename
			fileName := fmt.Sprintf("%s.xlsx", sanitizeFileName(bankName))
			xlsxPath := filepath.Join(xlsxDir, fileName)

			// Check if the .xlsx file has already been downloaded with the same entity tag
			if etag, ok := etags[url]; ok && etag != "" {
				return
			}

			// Download the .xlsx file and get the new entity tag
			newEtag, err := DownloadFile(url, xlsxPath)
			if err != nil {
				return
			}

			// Update the entity tags map with the new entity tag
			etags[url] = newEtag

			log.Printf("Downloaded %s", url)
		}(s)
	})

	// Wait for all the goroutines to finish
	wg.Wait()

	// Save the entity tags to the etags file
	if err := SaveEtags(etagsFile, etags); err != nil {
		log.Printf("Failed to save etags: %v", err)
	}

	log.Println("All downloads completed.")
}

// sanitizeFileName removes or replaces characters in the bank name that are not allowed in file names
func sanitizeFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}
