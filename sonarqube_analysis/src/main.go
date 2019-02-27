package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func main() {
	var pathToRepositories string = os.Args[1]
	var pathRepositoriesToAnalyse []string = iterateOverRepositories(pathToRepositories)

	bar := pb.StartNew(len(pathRepositoriesToAnalyse))
	for _, pathRepository := range pathRepositoriesToAnalyse {
		runAnalysis(pathRepository)
		bar.Increment()
	}
	bar.Finish()
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

func runAnalysis(pathRepository string) {
	log.Println("Start analysis of " + pathRepository)
	var hashs []string = retrieveCommitsMasterHash(pathRepository)
	for _, hash := range hashs {
		log.Println("Hash:" + hash)
		log.Println(len(hashs))
		_, err := exec.Command("sh", "-c", "cd "+pathRepository+" && git checkout -f "+hash).Output()
		if err != nil {
			log.Println(err)
		}
		sonarAnalysis(pathRepository, hash)
	}
	log.Println("End " + pathRepository)

}

func sonarAnalysis(pathRepository string, hash string) {
	var projectName string = path.Base(pathRepository)
	_, err := exec.Command("sh", "-c", "cd "+pathRepository+" && sonar-scanner -Dsonar.java.binaries=/tmp/empty -Dsonar.sources=. -Dsonar.projectKey="+projectName+" -Dsonar.projectVersion="+hash).Output()
	if err != nil {
		log.Println(err)
	}
}

func retrieveCommitsMasterHash(pathRepository string) []string {
	var hashs []string
	out, err := exec.Command("sh", "-c", "cd "+pathRepository+" && git log master --after='2018-02-15 00:00' --pretty=format:'%H'").Output()
	if err != nil {
		log.Println(err)
		return hashs
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()
		hashs = append(hashs, line)
	}

	return hashs
}
