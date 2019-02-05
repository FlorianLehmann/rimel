package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func main() {
	var repositoriesNames []string
	repositoriesNames = readFile(os.Args[1])
	cloneRepositories(repositoriesNames)
}

func readFile(path string) []string {
	var repositoriesNames []string
	file, err := os.Open(path)
	check_error(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		repositoriesNames = append(repositoriesNames, string(line))
	}

	return repositoriesNames
}

func cloneRepositories(repositoriesNames []string) {
	err := os.Mkdir("repositories", os.ModePerm)
	check_error(err)
	bar := pb.StartNew(len(repositoriesNames))
	for _, repositoryName := range repositoriesNames {
		directoryName := strings.Replace(repositoryName, "/", "_", -1)
		os.Mkdir("repositories/"+directoryName, os.ModePerm)
		cmd := "git clone https://github.com/" + repositoryName + " repositories/" + directoryName
		log.Println(cmd)
		err := exec.Command("sh", "-c", cmd).Run()
		if err != nil {
			log.Println(err)
		}
		bar.Increment()
	}
	bar.Finish()
}

func check_error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
