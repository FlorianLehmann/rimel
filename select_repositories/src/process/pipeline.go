package process

import (
	"../model"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type Pipeline struct {
	channels        []chan model.Repository
	filters         []Filter
	Done            chan bool
	outRepositories []model.Repository
}

func (p *Pipeline) AddFilter(filter Filter) {
	p.filters = append(p.filters, filter)
}

func (p *Pipeline) Execute(inRepositories []model.Repository) {

	p.createChannels()

	go p.sendDataToPipeline(inRepositories)
	go p.collectDataProcessed()

	p.startFilters()

}

func (p *Pipeline) GetFilteredRepositories() []model.Repository {
	return p.outRepositories
}

func (p *Pipeline) createChannels() {
	for index := 0; index < len(p.filters)+1; index++ {
		p.channels = append(p.channels, make(chan model.Repository))
	}
}

func (p *Pipeline) sendDataToPipeline(inRepositories []model.Repository) {
	count := len(inRepositories)
	bar := pb.StartNew(count)
	for _, repo := range inRepositories {
		bar.Increment()
		p.channels[0] <- repo
	}
	bar.Finish()
	close(p.channels[0])
}

func (p *Pipeline) collectDataProcessed() {
	for {
		repository, ok := <-p.channels[len(p.channels)-1]
		if !ok {
			p.Done <- true
			return
		}
		p.outRepositories = append(p.outRepositories, repository)
	}
}

func (p *Pipeline) startFilters() {
	for i, filter := range p.filters {
		go filter.execute(p.channels[i], p.channels[i+1])
	}
}
