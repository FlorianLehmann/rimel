package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var username string = "Username"
var password string = "Password"

const githubAPIUri = "https://api.github.com/"

func GetNumberOfContributors(owner string, projectName string) int {
	waitForCreditRateLimit()
	var client *http.Client = &http.Client{}
	req, err := http.NewRequest("GET", githubAPIUri+"repos/"+owner+"/"+projectName+"/contributors", nil)
	checkError(err)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	checkError(err)
	bodyText, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	checkError(err)
	regex, err := regexp.Compile("contributions")
	checkError(err)

	groups := regex.FindAllStringSubmatch(string(bodyText), -1)
	return len(groups)
}

func GetNumberOfContributionsLastYear(owner string, projectName string) int {
	waitForCreditRateLimit()
	var client *http.Client = &http.Client{}
	req, err := http.NewRequest("GET", githubAPIUri+"repos/"+owner+"/"+projectName+"/stats/commit_activity", nil)
	checkError(err)
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	checkError(err)
	bodyText, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	checkError(err)

	regex, err := regexp.Compile("\"total\":(\\d+),")
	checkError(err)

	groups := regex.FindAllStringSubmatch(string(bodyText), -1)
	var numberContributions int
	for _, group := range groups {
		contributions, _ := strconv.Atoi(group[1])
		numberContributions += contributions
	}

	return numberContributions
}

func waitForCreditRateLimit() {
	for {
		runtime.Gosched()

		if !checkRateLimit() {
			return
		}
		time.Sleep(time.Minute)
	}
}

func checkRateLimit() bool {
	var client *http.Client = &http.Client{}
	req, err := http.NewRequest("GET", githubAPIUri+"rate_limit", nil)
	checkError(err)

	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	checkError(err)

	remaining, _ := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	resp.Body.Close()
	return remaining < 10
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
