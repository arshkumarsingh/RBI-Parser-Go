package rbiparser

import (
    "encoding/csv"
    "fmt"
    "os"
    "regexp"
    "strings"
)

var (
    alphanumeric          = regexp.MustCompile(`[^a-zA-Z0-9]`)
    spaces                = regexp.MustCompile(`\s+`)
    excludeWords          = map[string]struct{}{"to": {}, "the": {}, "at": {}, "of": {}, "by": {}, "as": {}, "for": {}, "via": {}}
)

// CleanRow cleans a single row from the CSV.
func CleanRow(row []string) []string {
    row[0] = cleanName(row[0])
    row[1] = strings.ToUpper(row[1])
    row[2] = cleanMICR(row[2])
    row[3] = cleanLine(row[3], true)
    row[4] = cleanLine(row[4], true)
    row[6] = cleanLine(row[6], true)
    row[7] = cleanLine(row[7], true)
    row[8] = cleanLine(row[8], false)
    return row
}

func cleanName(name string) string {
    return strings.ToUpper(name)
}

func cleanMICR(micr string) string {
    if len(micr) > 5 {
        micr = alphanumeric.ReplaceAllString(micr, "")
    }
    return micr
}

func cleanLine(line string, complicated bool) string {
    line = spaces.ReplaceAllString(line, " ")
    return strings.TrimSpace(line)
}

// CleanCSV cleans and processes the CSV file.
func CleanCSV(src, dest string) error {
    file, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("can't open CSV file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return fmt.Errorf("can't read CSV file: %w", err)
    }

    outFile, err := os.Create(dest)
    if err != nil {
        return fmt.Errorf("can't create CSV file: %w", err)
    }
    defer outFile.Close()

    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    for i, record := range records {
        if i == 0 {
            writer.Write(record) // Write header
        } else {
            writer.Write(CleanRow(record))
        }
    }
    return nil
}
