package parser

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "golang.org/x/net/html"
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
func parseHTML(htmlContent string) (string, error) {
    doc, err := html.Parse(strings.NewReader(htmlContent))
    if err != nil {
        return "", fmt.Errorf("failed to parse HTML: %w", err)
    }

    var textContent strings.Builder
    extractText(doc, &textContent)
    return textContent.String(), nil
}

// extractText recursively extracts text from the HTML node.
func extractText(n *html.Node, textContent *strings.Builder) {
    if n.Type == html.TextNode {
        textContent.WriteString(n.Data)
        textContent.WriteString("\n")
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        extractText(c, textContent)
    }
}
