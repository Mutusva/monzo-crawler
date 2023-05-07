package monzo_interview

import "sync"

type Crawler interface {
	Start(bool) error
}

type Worker interface {
	Run(filters []string, visited map[string]bool)
	GetResultChan() chan map[string][]string
	GenerateJobs(filters []string, visited map[string]bool, wg *sync.WaitGroup)
}

type Job struct {
	// Links []string
	StartUrl   string
	Filters    []string
	ProcessFun func(curUrl string, filters []string) ([]string, error)
}

func (j *Job) Run() ([]string, error) {
	return j.ProcessFun(j.StartUrl, j.Filters)
}
