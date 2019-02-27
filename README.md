# RIMEL Project: Impact of minor contributors on code quality of open-source projects

This goal of this project is to find out whether or not minor contributors brings more software quality problems in their contributions than major contributors. Details of the experience can be found [here](https://github.com/RIMEL-UCA/Book/blob/master/git-make-merging-great-again/contents-2.md)

## Getting Started

This repository has 4 folders, each one containing parts of the experience :
-  "kaggle", which contains the dataset we used for the experience. The source of the dataset is [here](https://www.kaggle.com/github/github-repos).
-  "clone_repositories", which provides scripts that allows user to clone multiple repository from the output of the dataset filters.
-  "select_repositories", which filters each repository
-  "sonarcube_analysis", which peforms the sonarqube analysis
-  "sonarqube_delta_analysis", which computes the delta between metrics of each commits
