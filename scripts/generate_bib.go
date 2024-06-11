package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const style = "harvard-bournemouth-university"
const getCollectionsUrl = "https://api.zotero.org/groups/4536134/collections?limit=100&q=%s"
const getBibUrl = "https://api.zotero.org/groups/4536134/collections/%s/items?format=bib&%s"

func GenerateBib() {
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

			var siteName = site.Site
			resp, err := getCollections(siteName)
			if err != nil {
				panic(err)
			}

			if len(resp) <= 0 {
				break
			}

			for i, entry := range resp {
				entryNameSplit := strings.Split(entry.Data.Name, ",")
				if entryNameSplit[0] != siteName {
					break
				}

				siteBib, err := getBib(resp[i].Data.Key)
				if err != nil {
					panic(err)
				}

				site.Bibliography = siteBib
				writeBibliography(file.Name(), &site)
			}
		}
	}
}

func writeBibliography(fileName string, site *Site) error {
	fmt.Println(fileName)
	marshalledData, err := json.MarshalIndent(site, "", "   ")
	if err != nil {
		return err
	}

	err = os.WriteFile("../data/sites/"+fileName, marshalledData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getBib(id string) ([]string, error) {
	resp, err := http.Get(fmt.Sprintf(getBibUrl, id, style))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Failed to parse the HTML document", err)
	}

	foundCitations := doc.Find(".csl-entry")
	var references = make([]string, foundCitations.Length())
	foundCitations.Each(func(i int, s *goquery.Selection) {
		ref := s.Text()
		references[i] = ref
	})

	return references, nil
}

func getCollections(q string) ([]Entry, error) {
	requestUrl := fmt.Sprintf(getCollectionsUrl, url.QueryEscape(q))
	resp, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	err = json.Unmarshal(body, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

type Entry struct {
	Data Data
}

type Data struct {
	Key  string
	Name string
}
