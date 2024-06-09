package main

import (
    "fmt"
    "github.com/arshkumarsingh/RBI-Parser-Go/internal/parser"
)

func main() {
    url := "https://www.rbi.org.in/Scripts/BS_PressReleaseDisplay.aspx?prid=49515"
    data, err := parser.FetchData(url)
    if err != nil {
        fmt.Println("Error fetching data:", err)
        return
    }

    fmt.Println("Data fetched successfully:")
    fmt.Println(data)
}
