package main

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"o365-attack-toolkit/model"
	"o365-attack-toolkit/server"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gcfg.v1"
)

func main() {

	model.GlbConfig = model.Config{}
	err := gcfg.ReadFileInto(&model.GlbConfig, "template.conf")

	if err != nil {
		log.Fatal(err.Error())
	}

	//initializeRules()
	go server.StartExtServer(model.GlbConfig)
	server.StartIntServer(model.GlbConfig)
	fmt.Println(model.GlbConfig)
}

func initializeRules() {

	var ruleFiles []string
	var tempRule model.Rule

	root := "rules"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		ruleFiles = append(ruleFiles, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range ruleFiles {

		ruleFile, err := os.Open(file)

		if err != nil {
			log.Println(err)
		}

		defer ruleFile.Close()

		byteValue, _ := ioutil.ReadAll(ruleFile)

		json.Unmarshal(byteValue, &tempRule)

		model.GlbRules = append(model.GlbRules, tempRule)

	}

	log.Printf("Loaded %d rules successfully.", len(model.GlbRules))
}
