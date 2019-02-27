# Sonarcube Analysis

## Prerequisites

In order to run this script, you need to install the following tools:

```
Go (1.11.4)
Sonarcube Server
sonar-scanner (exported in PATH)
```

You will need to create the following commands:

```
mkdir /tmp/empty
```

## Executing

In order to run this script, you have to execute the following commands:

```
cd src
go get gopkg.in/cheggaaa/pb.v1
go build -o sonarqube_analysis
./sonarqube_analysis <path_to_repositories_folders>
```
Where <path_to_repositories_folders> is the path of the folder where all the cloned repositories are.
