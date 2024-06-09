package rbiparser

import (
    "encoding/csv"
    "fmt"
    "os"

    "github.com/tealeg/xlsx"
)

// ConvertXLSXToCSV converts an .xlsx file to a .csv file.
func ConvertXLSXToCSV(src, target string) error {
    xlFile, err := xlsx.OpenFile(src)
    if err != nil {
        return fmt.Errorf("can't open sheet: %w", err)
    }

    outFile, err := os.Create(target)
    if err != nil {
        return fmt.Errorf("can't create CSV file: %w", err)
    }
    defer outFile.Close()

    writer := csv.NewWriter(outFile)
    defer writer.Flush()

    for _, sheet := range xlFile.Sheets {
        for _, row := range sheet.Rows {
            var record []string
            for _, cell := range row.Cells {
                record = append(record, cell.String())
            }
            if err := writer.Write(record); err != nil {
                return fmt.Errorf("can't write to CSV: %w", err)
            }
        }
    }
    return nil
}
