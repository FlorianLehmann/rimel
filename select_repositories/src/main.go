package main

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"./model"
	"./process"
)

var repositories []model.Repository

func main() {
	readCsv(os.Args[1])

	var pipeline process.Pipeline
	pipeline.Done = make(chan bool)

	var f2 process.FilterByLanguageProportion = process.FilterByLanguageProportion{
		Names:      []string{"Python", "Java", "JavaScript"},
		Percentage: 51}
	var f4 process.FilterMinimumContributor = process.FilterMinimumContributor{
		Minimum: 2,
	}
	var f5 process.FilterLimitNumberProjects = process.FilterLimitNumberProjects{
		Limit: 4980,
	}
	var f1 process.FilterDuplicateProject = process.FilterDuplicateProject{
		RepositoriesFiltered: make(map[string]bool),
	}

	pipeline.AddFilter(&f1)
	pipeline.AddFilter(&f5)
	pipeline.AddFilter(&f2)
	pipeline.AddFilter(&f4)

	pipeline.Execute(repositories)

	<-pipeline.Done

	repositories = pipeline.GetFilteredRepositories()
	writeNameCompliantRepositories()
}

func readCsv(file string) {
	f, err := os.Open(file)
	checkError(err)
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	checkError(err)

	for _, line := range lines[1:] {
		var r model.Repository

		r.Organization = strings.Split(line[0], "/")[0]
		r.Name = strings.Split(line[0], "/")[1]

		regex, err := regexp.Compile("{'name': '([^']*)', 'bytes': (\\d+)}")
		checkError(err)

		groups := regex.FindAllStringSubmatch(line[1], -1)
		for _, group := range groups {
			var language model.Language
			language.Name = group[1]
			language.Bytes, _ = strconv.Atoi(group[2])
			r.Languages = append(r.Languages, language)
		}

		repositories = append(repositories, r)
	}
}

func writeNameCompliantRepositories() {
	f, err := os.Create("output")
	checkError(err)
	defer f.Close()

	for k, repository := range repositories {
		if k > 0 {
			f.WriteString("\n")
		}
		f.WriteString(repository.Organization + "/" + repository.Name)
	}

}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
