package main

import (
	"log"

	"codeberg.org/polarhive/pasta/util"
	"codeberg.org/polarhive/pasta/web"
)

func main() {
	db, err := util.InitDB()
	defer db.Close()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	web.Serve(db)
}
