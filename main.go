package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/stamblerre/sheets"
	gsheets "google.golang.org/api/sheets/v4"
)

var (
	urlFlag         = flag.String("url", "", "URL of existing Google sheets")
	credentialsFile = flag.String("credentials", "", "path to credentials file for Google Sheets")
	tokenFile       = flag.String("token", "", "path to token file for authentication in Google sheets")

	playlistFlag = flag.String("playlist", "", "Name of Spotify playlist to sync")
)

func main() {
	flag.Parse()

	if *playlistFlag == "" {
		log.Fatal("Please provide the --playlist flag.")
	}

	// Determine if the user has provided a valid Google Sheets URL.
	var spreadsheetID string
	if *urlFlag != "" {
		var err error
		spreadsheetID, err = sheets.GetSpreadsheetID(*urlFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	ctx := context.Background()

	srv, err := sheets.GoogleSheetsService(ctx, *credentialsFile, *tokenFile)
	if err != nil {
		log.Fatal(err)
	}
	var spreadsheet *gsheets.Spreadsheet
	rowData := make(map[string][]*gsheets.RowData)
	if *urlFlag == "new" {
		title := fmt.Sprintf("Spotify Sync %q as of %s", *playlistFlag, time.Now().Format("2006-01-02"))
		spreadsheet, err = sheets.CreateSheet(ctx, srv, title, rowData)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		spreadsheet, err = sheets.AppendToSheet(ctx, srv, spreadsheetID, rowData)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := sheets.ResizeColumns(ctx, srv, spreadsheet); err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote data to Google Sheet: %s\n", spreadsheet.SpreadsheetUrl)
}
