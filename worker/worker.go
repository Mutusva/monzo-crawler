package worker

import (
	monzo_interview "github.com/Mutusva/monzo-webcrawler"
	"sync"
)

type worker struct {
	PoolCount int
	Jobs      chan monzo_interview.Job
	Results   chan map[string][]string
	mx        *sync.Mutex
	ProcessFn func(curUrl string, filters []string) ([]string, error)
	Queue     []string
}

func (w *worker) GenerateJobs(filters []string) {
	for len(w.Queue) > 0 {
		curUrl := w.Queue[0]
		w.Queue = w.Queue[1:]

		job := monzo_interview.Job{
			StartUrl:   curUrl,
			Filters:    filters,
			ProcessFun: w.ProcessFn,
		}

		w.Jobs <- job
	}
}

func (w *worker) DoWork(wg *sync.WaitGroup, jobs <-chan monzo_interview.Job, result chan<- map[string][]string) {
	defer wg.Done()
	for job := range jobs {
		links, _ := job.Run()
		// add error link to another queue
		w.mx.Lock()
		w.Queue = append(w.Queue, links...)
		w.mx.Unlock()

		result <- map[string][]string{
			job.StartUrl: links,
		}
	}
}

func (w *worker) Run() {
	var wg *sync.WaitGroup
	for i := 0; i < w.PoolCount; i++ {
		wg.Add(1)
		go w.DoWork(wg, w.Jobs, w.Results)
	}
}
