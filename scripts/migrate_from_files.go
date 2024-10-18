package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"codeberg.org/polarhive/pasta/util"
)

func main() {
	db, err := util.NewDB("pastebin.db")
	dataDir := "data" // Directory to update from

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	files, err := ioutil.ReadDir(dataDir)
	if err != nil {
		log.Fatalf("Failed to read data directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		id := file.Name()
		filePath := filepath.Join(dataDir, id)

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", id, err)
			continue
		}

		paste := &util.Paste{
			ID:      id,
			Content: string(content),
		}

		err = db.Create(paste)
		if err != nil {
			log.Printf("Failed to create paste for file %s: %v", id, err)
			continue
		}

		fmt.Printf("Migrated file: %s\n", id)
	}

	fmt.Println("Migration completed.")
}
