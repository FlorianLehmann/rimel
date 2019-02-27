# Sonarqube delta analysis

The provided script gather all results of sonarqube analysis and compute the delta of four metrics (bugs, vulnerabilities, duplications, code smell) for each commit in the master branch.

## Prerequisites

In order to run this script, you need to install the following tools:

```
Go (1.11.4)
Git
Sonarqube (running)
```

## Executing

In order to run this script, you have to execute the following commands:

```
cd src
go get gopkg.in/cheggaaa/pb.v1
go build -o sonarqube_delta_analysis
./sonarqube_delta_analysis <path_to_repositories_folders>
```

Where <path_to_repositories_folders> is the path of the folder where all the cloned repositories are.
