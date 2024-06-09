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

func init() {
	// Load bank names from JSON file
	file, err := os.Open("banks.json")
	if err != nil {
		log.Fatalf("Could not open banks.json: %v", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&bankNames)
	if err != nil {
		log.Fatalf("Could not decode banks.json: %v", err)
	}
}

// LoadEtags loads the etags from a file
func LoadEtags(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	etags := make(map[string]string)
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

// ExtractBankNameFromContext extracts the bank name from the surrounding context of the URL
func ExtractBankNameFromContext(s *goquery.Selection) (string, error) {
	// Look for the closest preceding sibling that is a text node containing a bank name
	parent := s.Parent()
	text := parent.Text()
	for name := range bankNames {
		matched, err := regexp.MatchString(`(?i)`+regexp.QuoteMeta(name), text)
		if err != nil {
			return "", err
		}
		if matched {
			return bankNames[name], nil
		}
	}

	// If no match is found, try looking one level up in the DOM
	grandParent := parent.Parent()
	text = grandParent.Text()
	for name := range bankNames {
		matched, err := regexp.MatchString(`(?i)`+regexp.QuoteMeta(name), text)
		if err != nil {
			return "", err
		}
		if matched {
			return bankNames[name], nil
		}
	}

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
	var wg sync.WaitGroup

	etags, err := LoadEtags(etagsFile)
	if err != nil {
		log.Printf("Could not load etags: %v", err)
		etags = make(map[string]string)
	}

	resp, err := http.Get(scrapeURL)
	if err != nil {
		log.Fatalf("Error fetching URL %s: %v", scrapeURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	doc.Find("a[href$='.xlsx']").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		go func(s *goquery.Selection) {
			defer wg.Done()
			url, exists := s.Attr("href")
			if !exists {
				log.Println("No href attribute found")
				return
			}

			bankName, err := ExtractBankNameFromContext(s)
			if err != nil {
				log.Printf("Failed to extract bank name for URL %s: %v", url, err)
				return
			}

			// Sanitize the bank name for use as a filename
			fileName := fmt.Sprintf("%s.xlsx", sanitizeFileName(bankName))
			xlsxPath := filepath.Join(xlsxDir, fileName)

			if etag, ok := etags[url]; ok && etag != "" {
				log.Printf("%s already downloaded with etag %s, skipping", url, etag)
				return
			}

			newEtag, err := DownloadFile(url, xlsxPath)
			if err != nil {
				log.Printf("Failed to download %s: %v", url, err)
				return
			}

			etags[url] = newEtag
			if err := SaveEtags(etagsFile, etags); err != nil {
				log.Printf("Failed to save etags: %v", err)
			}

			log.Printf("Downloaded %s", url)
		}(s)
	})

	wg.Wait()
	log.Println("All downloads completed.")
}

// sanitizeFileName removes or replaces characters in the bank name that are not allowed in file names
func sanitizeFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}
