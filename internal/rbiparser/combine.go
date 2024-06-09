package rbiparser

import (
    "encoding/csv"
    "fmt"
    "os"
    "path/filepath"
)

// CombineCSVs combines multiple CSV files into one.
func CombineCSVs(srcDir, destFile string) error {
    outFile, err := os.Create(destFile)
    if err != nil {
        return fmt.Errorf("can't create master CSV file: %w", err)
    }
    defer outFile.Close()

    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    headers := []string{"BANK", "IFSC", "MICR", "BRANCH", "ADDRESS", "CONTACT", "CITY", "DISTRICT", "STATE", "ABBREVIATION"}
    if err := writer.Write(headers); err != nil {
        return fmt.Errorf("can't write headers to CSV: %w", err)
    }

    err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".csv" {
            if err := appendCSV(path, writer); err != nil {
                return err
            }
        }
        return nil
    })
    if err != nil {
        return fmt.Errorf("failed to combine CSV files: %w", err)
    }

    return nil
}

func appendCSV(filePath string, writer *csv.Writer) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("can't open CSV file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return fmt.Errorf("can't read CSV file: %w", err)
    }

    for _, record := range records[1:] { // Skip header
        writer.Write(record)
    }
    return nil
}
