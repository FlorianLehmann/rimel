package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	Hash            string
	Bugs            int
	Vulnerabilities int
	CodeSmells      int
	DuplicatedLines int
	Date            string
}

type Contributor struct {
	Hashs   []string
	name    string
	percent int
	change  int
}

type ApiSonarSearchAnalyses struct {
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	Analyses []struct {
		Key    string `json:"key"`
		Date   string `json:"date"`
		Events []struct {
			Key         string `json:"key"`
			Category    string `json:"category"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
		} `json:"events"`
	} `json:"analyses"`
}

type ApiSonarMeasures struct {
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	Measures []struct {
		Metric  string `json:"metric"`
		History []struct {
			Date  string `json:"date"`
			Value string `json:"value"`
		} `json:"history"`
	} `json:"measures"`
}

func main() {
	var pathToRepositories string = os.Args[1]
	var pathRepositoriesToAnalyse []string = iterateOverRepositories(pathToRepositories)

	for _, pathRepository := range pathRepositoriesToAnalyse {
		var contributors []Contributor = getContributors(pathRepository)
		var results []Result = retrieveSonarResult(pathRepository)
		an := computeDelta(contributors, results, pathRepository)

		log.Panicln("***************************")
		log.Println(pathRepository)
		log.Printf("%+v\n", an)

		var numberMinoritaryContributors, numberMajoritaryContributors int
		var numberMinoritaryContributorModifications, numberMajoritaryContributorModifications int
		for _, c := range contributors {
			if c.percent <= 5 {
				numberMinoritaryContributors++
				numberMinoritaryContributorModifications += c.change
			} else {
				numberMajoritaryContributors++
				numberMajoritaryContributorModifications += c.change
			}
		}

		fmt.Print("numberMinoritaryContributors")
		fmt.Println(numberMinoritaryContributors)
		fmt.Print("numberMajoritaryContributors")
		fmt.Println(numberMajoritaryContributors)
		fmt.Print("numberMinoritaryContributorModifications")
		fmt.Println(numberMinoritaryContributorModifications)
		fmt.Print("numberMajoritaryContributorModifications")
		fmt.Println(numberMajoritaryContributorModifications)

		log.Panicln("***************************")

	}
}

func iterateOverRepositories(pathToRepositories string) []string {
	var repositoriesToAnalyse []string
	files, err := ioutil.ReadDir(pathToRepositories)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			repositoriesToAnalyse = append(repositoriesToAnalyse, path.Join(pathToRepositories, file.Name()))
		}
	}

	return repositoriesToAnalyse
}

func getTotalNumberContribution(pathRepository string) int {
	out, _ := exec.Command("sh", "-c", "cd "+pathRepository+" && git log  master --pretty=format:'%H' | wc -l ").Output()
	value, _ := strconv.Atoi(strings.Replace(strings.TrimSpace(string(out)), "\n", "", -1))
	return value
}

func getContributorNames(pathRepository string) []string {
	out, _ := exec.Command("sh", "-c", "cd "+pathRepository+" && git log  master --pretty=format:'%an' | sort -u ").Output()
	return strings.Split(string(out), "\n")
}

func getNumberContribution(name string, pathRepository string) int {
	out, _ := exec.Command("sh", "-c", "cd "+pathRepository+" && git log  master --pretty=format:'%an' | grep '"+name+"' | wc -l ").Output()
	numberOfContributions, _ := strconv.Atoi(strings.Replace(string(out), "\n", "", -1))
	return numberOfContributions
}

func getHashContributions(name string, pathRepository string) []string {
	out, _ := exec.Command("sh", "-c", "cd "+pathRepository+" && git log  master --author='"+name+"' --pretty=format:'%H'").Output()
	return strings.Split(string(out), "\n")
}

func getContributors(pathRepository string) []Contributor {
	var contributors []Contributor

	for _, name := range getContributorNames(pathRepository) {
		percentage := float64(getNumberContribution(name, pathRepository)) / float64((getTotalNumberContribution(pathRepository)))
		hashs := getHashContributions(name, pathRepository)
		contributors = append(contributors, Contributor{hashs, name, int(percentage * 100), getNumberModifications(name, pathRepository)})
	}

	return contributors
}

func getNumberModifications(name string, pathRepository string) int {
	out, _ := exec.Command("sh", "-c", "cd "+pathRepository+" && git log --author='"+name+"' --oneline --shortstat").Output()

	var numberModifications int

	regex, err := regexp.Compile("(\\d+) insertion")
	if err != nil {
		log.Fatal(err)
	}
	groups := regex.FindAllStringSubmatch(string(out), -1)
	for _, group := range groups {
		contributions, _ := strconv.Atoi(group[1])
		numberModifications += contributions
	}

	regex, err = regexp.Compile("(\\d+) deletions")
	if err != nil {
		log.Fatal(err)
	}
	groups = regex.FindAllStringSubmatch(string(out), -1)
	for _, group := range groups {
		contributions, _ := strconv.Atoi(group[1])
		numberModifications += contributions
	}

	return numberModifications
}

func searchAnalysis(repoName string, page int) ApiSonarSearchAnalyses {
	resp, err := http.Get("http://192.168.1.41:9000/api/project_analyses/search?project=" + repoName + "&category=VERSION&p=" + strconv.Itoa(page))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var apiSonarVersions ApiSonarSearchAnalyses
	json.Unmarshal(body, &apiSonarVersions)
	return apiSonarVersions
}

func gatherResultAnalysis(repoName string, date string) ApiSonarMeasures {
	resp, err := http.Get("http://192.168.1.41:9000/api/measures/search_history?component=" + repoName + "&metrics=bugs,vulnerabilities,code_smells,duplicated_lines&from=" + date + "&to=" + date)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var apiSonarVersionMeasures ApiSonarMeasures
	json.Unmarshal(body, &apiSonarVersionMeasures)
	return apiSonarVersionMeasures
}

func retrieveSonarResult(pathRepository string) []Result {

	repoName := path.Base(pathRepository)

	var result []Result

	var TotalSize int = -1
	var page int = 1

	for {

		var apiSonarVersions ApiSonarSearchAnalyses = searchAnalysis(repoName, page)

		for _, analyse := range apiSonarVersions.Analyses {
			date := strings.Replace(analyse.Date, "+", "%2B", -1)

			var apiSonarVersionMeasures ApiSonarMeasures = gatherResultAnalysis(repoName, date)

			Bugs, err := strconv.Atoi(apiSonarVersionMeasures.Measures[2].History[0].Value)
			if err != nil {
				log.Fatal(err)
			}
			Vulnerabilities, err := strconv.Atoi(apiSonarVersionMeasures.Measures[3].History[0].Value)
			if err != nil {
				log.Fatal(err)
			}
			CodeSmells, err := strconv.Atoi(apiSonarVersionMeasures.Measures[1].History[0].Value)
			if err != nil {
				log.Fatal(err)
			}
			DuplicatedLines, err := strconv.Atoi(apiSonarVersionMeasures.Measures[0].History[0].Value)
			if err != nil {
				log.Fatal(err)
			}

			result = append(result, Result{analyse.Events[0].Name,
				Bugs,
				Vulnerabilities,
				CodeSmells,
				DuplicatedLines,
				date})
		}

		if TotalSize == -1 {
			TotalSize = apiSonarVersions.Paging.Total
		}

		if apiSonarVersions.Paging.PageSize < TotalSize {
			TotalSize -= apiSonarVersions.Paging.PageSize
			page++
		} else {
			break
		}
	}

	return result

}

func getAllHashsRepository(pathRepository string) []string {
	out, _ := exec.Command("sh", "-c", "cd "+pathRepository+" && git log  master --pretty=format:'%H'").Output()
	return strings.Split(string(out), "\n")
}

type AnalysisMeasures struct {
	nb_bugs_added_by_Minority            int
	nb_bugs_added_by_Majority            int
	nb_vulnerabilities_added_by_Minority int
	nb_vulnerabilities_added_by_Majority int
	nb_smell_added_by_Minority           int
	nb_smell_added_by_Majority           int
	nb_duplicated_added_by_Minority      int
	nb_duplicated_added_by_Majority      int
}

func computeDelta(contributors []Contributor, results []Result, pathRepository string) AnalysisMeasures {
	hashs := getAllHashsRepository(pathRepository)

	var an AnalysisMeasures

	for index := len(hashs) - 1; index > 1; index-- {

		var oldResult Result = findResult(hashs[index-1], results)
		var newResult Result = findResult(hashs[index], results)
		var contributeur Contributor = findContributor(hashs[index], contributors)

		if contributeur.percent <= 5 {
			if oldResult.Bugs < newResult.Bugs {
				an.nb_bugs_added_by_Minority += newResult.Bugs - oldResult.Bugs
			}
			if oldResult.Vulnerabilities < newResult.Vulnerabilities {
				an.nb_vulnerabilities_added_by_Minority += newResult.Vulnerabilities - oldResult.Vulnerabilities
			}
			if oldResult.CodeSmells < newResult.CodeSmells {
				an.nb_smell_added_by_Minority += newResult.CodeSmells - oldResult.CodeSmells
			}
			if oldResult.DuplicatedLines < newResult.DuplicatedLines {
				an.nb_duplicated_added_by_Minority += newResult.DuplicatedLines - oldResult.DuplicatedLines
			}
		} else {
			if oldResult.Bugs < newResult.Bugs {
				an.nb_bugs_added_by_Majority += newResult.Bugs - oldResult.Bugs
			}
			if oldResult.Vulnerabilities < newResult.Vulnerabilities {
				an.nb_vulnerabilities_added_by_Majority += newResult.Vulnerabilities - oldResult.Vulnerabilities
			}
			if oldResult.CodeSmells < newResult.CodeSmells {
				an.nb_smell_added_by_Majority += newResult.CodeSmells - oldResult.CodeSmells
			}
			if oldResult.DuplicatedLines < newResult.DuplicatedLines {
				an.nb_duplicated_added_by_Majority += newResult.DuplicatedLines - oldResult.DuplicatedLines
			}
		}
	}

	return an
}

func findResult(hash string, results []Result) Result {
	for _, res := range results {
		if res.Hash == hash {
			return res
		}
	}
	return Result{}
}

func findContributor(hash string, contributors []Contributor) Contributor {
	for _, contributor := range contributors {
		for _, hash1 := range contributor.Hashs {
			if hash1 == hash {
				return contributor
			}
		}
	}
	return Contributor{}
}
