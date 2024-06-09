package main

import (
    "log"
    "os"
    "download"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatalf("Usage: %s <command> [<args>]\nCommands:\n  download\n  convert\n  combine", os.Args[0])
    }

    switch os.Args[1] {
    case "download":
        download.DownloadAll("https://www.rbi.org.in/scripts/bs_viewcontent.aspx?Id=2009", "downloads", "etags.json")
    case "convert":
        // Implement conversion logic here if needed
    case "combine":
        // Implement combine logic here if needed
    default:
        log.Fatalf("Unknown command: %s", os.Args[1])
    }
}
