# Selection of repositories

## Prerequisites

In order to run this script, you need to install the following tools:

```
Go (1.11.4)
```

You have to provide your username and password in `src/api/github.go` for making call to the Github API.

```go
var username string = "Username"
var password string = "Password"
```

## Executing

In order to run this script, you have to execute the following commands:

```
cd src
go get gopkg.in/cheggaaa/pb.v1
go build -o repositories_selection
./repositories_selection <kaggle_dataset>
```

<kaggle_dataset> is the path of the dataset containning a large amount of repositories. This file also contains the languages used in those projects.
