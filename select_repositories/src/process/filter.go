package process

import (
	"../api"
	"../model"
)

type Filter interface {
	execute(inRepositories chan model.Repository, outRepositories chan model.Repository)
}

type FilterByLanguageProportion struct {
	Names      []string
	Percentage int
	count      int
}

func (f *FilterByLanguageProportion) execute(inRepositories chan model.Repository, outRepositories chan model.Repository) {
	var repository model.Repository
	var open bool
	for {
		repository, open = <-inRepositories
		if !open {
			close(outRepositories)
			return
		}
		for _, name := range f.Names {
			if repository.PercentageLanguage(name) > f.Percentage {
				outRepositories <- repository
				break
			}
		}
		f.count++
	}

}

type FilterByAmountBytes struct {
	Amount int
}

func (f *FilterByAmountBytes) execute(inRepositories chan model.Repository, outRepositories chan model.Repository) {
	var repository model.Repository
	var open bool
	for {
		repository, open = <-inRepositories
		if !open {
			close(outRepositories)
			return
		} else if repository.GetTotalLanguagesBytes() > f.Amount {
			outRepositories <- repository
		}
	}
}

type FilterLimitNumberProjects struct {
	Limit           int
	projectAnalyzed int
}

func (f *FilterLimitNumberProjects) execute(inRepositories chan model.Repository, outRepositories chan model.Repository) {
	var repository model.Repository
	var open bool
	for {
		repository, open = <-inRepositories
		if !open {
			close(outRepositories)
			return
		} else if f.projectAnalyzed < f.Limit {
			outRepositories <- repository
			f.projectAnalyzed++
		}
	}
}

type FilterMinimumContributor struct {
	Minimum int
}

func (f *FilterMinimumContributor) execute(inRepositories chan model.Repository, outRepositories chan model.Repository) {
	var repository model.Repository
	var open bool
	for {
		repository, open = <-inRepositories
		if !open {
			close(outRepositories)
			return
		} else if api.GetNumberOfContributors(repository.Organization, repository.Name) >= f.Minimum {
			outRepositories <- repository
		}
	}
}

type FilterMinimumContributionLastYear struct {
	Minimum int
}

func (f *FilterMinimumContributionLastYear) execute(inRepositories chan model.Repository, outRepositories chan model.Repository) {
	var repository model.Repository
	var open bool
	for {
		repository, open = <-inRepositories
		if !open {
			close(outRepositories)
			return
		} else if api.GetNumberOfContributionsLastYear(repository.Organization, repository.Name) >= f.Minimum {
			outRepositories <- repository
		}
	}
}

type FilterDuplicateProject struct {
	RepositoriesFiltered map[string]bool
}

func (f *FilterDuplicateProject) execute(inRepositories chan model.Repository, outRepositories chan model.Repository) {
	var repository model.Repository
	var open bool
	for {
		repository, open = <-inRepositories
		repositoryName := repository.Organization + "/" + repository.Name
		if !open {
			close(outRepositories)
			return
		} else if _, isPresent := f.RepositoriesFiltered[repositoryName]; !isPresent {
			f.RepositoriesFiltered[repositoryName] = true
			outRepositories <- repository
		}
	}
}
