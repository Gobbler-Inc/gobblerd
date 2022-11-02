package main

import (
	"archive/zip"
	"bbrz/parser"
	"bbrz/writer/json"
	"encoding/xml"
	"fmt"
	"log"
)

const file = "Coach-495574-745a245967fb458af7274335d1a2a626_2022-10-05_17_35_01.bbrz"

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

	if err := json.WriteJSON(rr); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully written output")
}
