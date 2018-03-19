package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Connect to database...")
	//connect to the database
	db, err := sql.Open("sqlite3", "mediaDB.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createtable(db)
	log.Println("creating table...")
	if err != nil {
		log.Fatal(err)
	}
}
func createtable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE 
						Filelists(
								fileid               INTEGER  NOT NULL PRIMARY KEY,
								filename             TEXT  NOT NULL  ,
								filesuffix           TEXT    ,
								filelocation         TEXT    ,
								filesize             NUMERIC    ,
								filehash             TEXT    ,
								filedate             TIMESTAMP,
								rowaction            TEXT  NOT NULL  ,
								rowactiondatetime    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP 
								);
							`)
	log.Printf("error: %s", err)
	return err
}
