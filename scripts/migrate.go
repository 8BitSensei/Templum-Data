package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Migrate() {
	fmt.Println("Reading files..")
	file, err := os.ReadFile("../data/templum_sites.json")
	if err != nil {
		panic(err)
	}

	var site Sites
	err = json.Unmarshal(file, &site)
	if err != nil {
		panic(err)
	}

	for _, site := range site.Sites {
		siteData, err := json.MarshalIndent(site, "", "   ")
		if err != nil {
			panic(err)
		}

		newFileName := strings.Replace(site.Site, " ", "-", -1) + ".json"
		newFile, err := os.Create("../data/sites/" + strings.ToLower(newFileName))
		if err != nil {
			panic(err)
		}

		_, err = newFile.WriteString(string(siteData))
	}

}

type Sites struct {
	Sites []Site `json:"sites"`
}
