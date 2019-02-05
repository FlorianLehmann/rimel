package model

type Repository struct {
	Organization   string
	Name           string
	NbCommit       int
	NbContributors int
	Languages      []Language
}

type Language struct {
	Name  string
	Bytes int
}

func (r *Repository) GetTotalLanguagesBytes() int {
	var sum int
	for _, language := range r.Languages {
		sum += language.Bytes
	}
	return sum
}

func (r *Repository) PercentageLanguage(languageName string) int {
	for _, language := range r.Languages {
		if language.Name == languageName {
			if totalBytes := r.GetTotalLanguagesBytes(); totalBytes != 0 {
				return int((float32(language.Bytes) / float32(totalBytes)) * 100.)
			}
		}
	}
	return 0
}
