package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/gobbler-inc/gobblerd/parser"
)

const file = "testdata/test.bbrz"

func main() {
	r, err := zip.OpenReader(file)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	//b := bytes.NewBuffer(nil)
	//ob := bytes.NewBuffer(nil)
	//fp, err := os.OpenFile("out.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//defer fp.Close()

	f := r.File[0]

	fmt.Printf("Contents of %s:\n", f.Name)
	rc, err := f.Open()
	if err != nil {
		log.Fatal(err)
	}
	//_, err = io.Copy(b, rc)
	//if err != nil {
	//	log.Fatal(err)
	//}
	defer rc.Close()

	var rr parser.Replay
	decoder := xml.NewDecoder(rc)
	err = decoder.Decode(&rr)
	if err != nil {
		panic(err)
	}

	record := parser.NewRecordFromReplay(rr)

	log.Printf("%s (%s) MVP: %s", record.Home.Name, record.Home.Race, record.Home.MVP)
	log.Printf("%s (%s) MVP: %s", record.Away.Name, record.Away.Race, record.Away.MVP)

	// spew.Dump(record)

	// if err := json.WriteJSON(rr); err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("Successfully written output")
}
