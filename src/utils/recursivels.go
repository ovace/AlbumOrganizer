package main

import (
	"flag"
	"fmt"
	"log"

	"./auth"

	"google.golang.org/api/drive/v2"
)

// https://developers.google.com/drive/v2/reference/files

var (
	queryFlag = flag.String("query", "", "Query for the root of the tree")
	idFlag    = flag.String("id", "", "File Id for the root of the tree")
)

// Lists files on google drive.
func main() {
	flag.Parse()
	if *queryFlag == "" && *idFlag == "" {
		log.Fatal("Either --id or --query is required.")
	}

	client, err := auth.DoAuth()
	if err != nil {
		log.Fatal(err)
	}

	query := *queryFlag
	if *idFlag != "" {
		query = fmt.Sprintf(`'%s' in parents`, *idFlag)
	}

	svc, err := drive.New(client)
	if err != nil {
		log.Fatalf("An error occurred opening driveservice: %v\n", err)
	}

	processQuery(svc, query)
}

func processQuery(svc *drive.Service, query string) {
	list, err := svc.Files.List().Q(query).Do()
	if err != nil {
		log.Fatalf("An error occurred listing files: %v\n", err)
	}
	var folders []string
	for _, item := range list.Items {
		if item.MimeType == "application/vnd.google-apps.folder" {
			folders = append(folders, item.Id)
			continue
		}
		printFile(item)
	}
	for _, item := range folders {
		query := fmt.Sprintf(`'%s' in parents`, item)
		processQuery(svc, query)
	}
}

func printFile(file *drive.File) {
	// TODO: Print more useful things
	fmt.Printf("%s %s\n", file.Id, file.Title)
}
