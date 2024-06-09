package main

import (
    "fmt"
    "github.com/arshkumarsingh/RBI-Parser-Go/internal/rbiparser"
    "os"
    "path/filepath"
)

const (
    scrapeURL  = "https://www.rbi.org.in/Scripts/BS_PressReleaseDisplay.aspx?prid=49515"
    xlsxDir    = "downloads/xlsx"
    csvDir     = "downloads/csv"
    masterFile = "master.csv"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: rbiparser <download|convert|combine>")
        return
    }

    switch os.Args[1] {
    case "download":
        err := downloadFiles()
        if err != nil {
            fmt.Println("Error downloading files:", err)
        } else {
            fmt.Println("Files downloaded successfully.")
        }
    case "convert":
        err := convertFiles()
        if err != nil {
            fmt.Println("Error converting files:", err)
        } else {
            fmt.Println("Files converted successfully.")
        }
    case "combine":
        err := combineFiles()
        if err != nil {
            fmt.Println("Error combining files:", err)
        } else {
            fmt.Println("Files combined successfully.")
        }
    default:
        fmt.Println("Unknown command. Usage: rbiparser <download|convert|combine>")
    }
}

func downloadFiles() error {
    urls, err := rbiparser.GetSheetURLs(scrapeURL)
    if err != nil {
        return err
    }

    if err := os.MkdirAll(xlsxDir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    for _, url := range urls {
        fileName := filepath.Base(url)
        filePath := filepath.Join(xlsxDir, fileName)
        if err := rbiparser.DownloadFile(url, filePath); err != nil {
            return err
        }
    }
    return nil
}

func convertFiles() error {
    if err := os.MkdirAll(csvDir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    err := filepath.Walk(xlsxDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if filepath.Ext(path) == ".xlsx" {
            csvPath := filepath.Join(csvDir, filepath.Base(path)+".csv")
            if err := rbiparser.ConvertXLSXToCSV(path, csvPath); err != nil {
                return err
            }
        }
        return nil
    })
    return err
}

func combineFiles() error {
    return rbiparser.CombineCSVs(csvDir, masterFile)
}
