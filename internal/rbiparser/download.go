package rbiparser

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"

    "golang.org/x/net/html"
    "golang.org/x/net/html/atom"
)

// GetSheetURLs scrapes the RBI page and returns the list of .xlsx URLs.
func GetSheetURLs(url string) ([]string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch URL: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
    }

    doc, err := html.Parse(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse HTML: %w", err)
    }

    var urls []string
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.DataAtom == atom.A {
            for _, a := range n.Attr {
                if a.Key == "href" && strings.HasSuffix(a.Val, ".xlsx") {
                    urls = append(urls, a.Val)
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
    return urls, nil
}

// DownloadFile downloads a file from the given URL and saves it to the specified path.
func DownloadFile(url, filepath string) error {
    resp, err := http.Get(url)
    if err != nil {
        return fmt.Errorf("failed to download file: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
    }

    out, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("failed to save file: %w", err)
    }
    return nil
}
