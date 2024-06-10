package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

func GenerateMetadata() {
	fmt.Println("Reading files..")
	files, err := os.ReadDir("../data/sites/")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			fmt.Println("Reading " + file.Name())
			byteValue, err := os.ReadFile("../data/sites/" + file.Name())
			if err != nil {
				panic(err)
			}

			var site Site
			err = json.Unmarshal(byteValue, &site)
			if err != nil {
				panic(err)
			}

			var metadata = &SiteMetadata{
				Category:    "data",
				Filename:    file.Name(),
				Weight:      1,
				Title:       site.Site,
				Description: "Bibliographic data on " + site.Site,
				Type:        "application/json",
			}

			yamlData, err := yaml.Marshal(&metadata)
			if err != nil {
				panic(err)
			}

			var metadataFilename = strings.Replace(file.Name(), ".json", ".data", 1)
			metadataFile, err := os.Create("../data/sites/" + metadataFilename)
			if err != nil {
				panic(err)
			}

			defer func(metadataFile *os.File) {
				err := metadataFile.Close()
				if err != nil {
					panic(err)
				}
			}(metadataFile)

			_, err = metadataFile.WriteString("---\n" + string(yamlData) + "---")
		}
	}
}

type Site struct {
	Site         string   `json:"site"`
	Start        string   `json:"start"`
	End          string   `json:"end"`
	Latitude     string   `json:"latitude"`
	Longitude    string   `json:"longitude"`
	Status       string   `json:"status"`
	Location     string   `json:"location"`
	Tags         string   `json:"tags"`
	Description  string   `json:"description"`
	Bibliography []string `json:"bibliography"`
}

type SiteMetadata struct {
	Category    string
	Filename    string
	Weight      int
	Title       string
	Description string
	Type        string
}
