package download

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "sync"

    "github.com/PuerkitoBio/goquery"
)

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

// GetSheetURLs fetches the list of .xlsx URLs from the RBI page
func GetSheetURLs(url string) ([]string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error fetching URL %s: %v", url, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
    }

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error parsing HTML: %v", err)
    }

    var urls []string
    doc.Find("a[href$='.xlsx']").Each(func(i int, s *goquery.Selection) {
        href, exists := s.Attr("href")
        if exists {
            urls = append(urls, href)
        }
    })

    if len(urls) < 1 {
        return nil, fmt.Errorf("couldn't find any .xlsx URLs")
    }

    return urls, nil
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
    urls, err := GetSheetURLs(scrapeURL)
    if err != nil {
        log.Fatalf("Failed to get sheet URLs: %v", err)
    }

    etags, err := LoadEtags(etagsFile)
    if err != nil {
        log.Printf("Could not load etags: %v", err)
        etags = make(map[string]string)
    }

    for _, url := range urls {
        wg.Add(1)
        go func(url string) {
            defer wg.Done()

            fname := filepath.Base(url)
            xlsxPath := filepath.Join(xlsxDir, fname)

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
        }(url)
    }

    wg.Wait()
    log.Println("All downloads completed.")
}
