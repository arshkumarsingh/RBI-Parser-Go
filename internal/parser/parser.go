package parser

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
)

// FetchData fetches and parses data from the given URL.
func FetchData(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("failed to fetch URL: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("received non-200 response: %d", resp.StatusCode)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body: %w", err)
    }

    data, err := parseHTML(string(body))
    if err != nil {
        return "", fmt.Errorf("failed to parse HTML: %w", err)
    }

    return data, nil
}

// parseHTML parses the HTML content and extracts the relevant data.
func parseHTML(html string) (string, error) {
    // Regex to extract data, for example: <p>(.*?)</p>
    re := regexp.MustCompile(`<p>(.*?)</p>`)
    matches := re.FindAllStringSubmatch(html, -1)
    
    if matches == nil {
        return "", fmt.Errorf("no matches found")
    }

    // Extract and concatenate matches
    var data string
    for _, match := range matches {
        data += match[1] + "\n"
    }

    return data, nil
}
