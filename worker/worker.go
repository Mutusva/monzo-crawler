package worker

import (
	"fmt"
	monzo_interview "github.com/Mutusva/monzo-webcrawler"
	"sync"
)

type ProcessFn func(curUrl string, filters []string) ([]string, error)

type worker struct {
	PoolCount int
	Jobs      chan monzo_interview.Job
	Results   chan map[string][]string
	links     chan string
	Mx        *sync.Mutex
	ProcessFn ProcessFn
	Queue     []string
}

func (w *worker) GenerateJobs(filters []string, visited map[string]bool) {

	for len(w.Queue) > 0 {

		w.Mx.Lock()
		curUrl := w.Queue[0]
		w.Queue = w.Queue[1:]
		w.Mx.Unlock()

		if visited[curUrl] {
			continue
		}

		job := monzo_interview.Job{
			StartUrl:   curUrl,
			Filters:    filters,
			ProcessFun: w.ProcessFn,
		}

		w.Jobs <- job

		for link := range w.links {
			w.Mx.Lock()
			w.Queue = append(w.Queue, link)
			w.Mx.Unlock()
		}
		visited[curUrl] = true
	}
}

func (w *worker) DoWork(wg *sync.WaitGroup, jobs <-chan monzo_interview.Job, result chan<- map[string][]string) {
	defer wg.Done()
	for job := range jobs {
		wg.Add(1)
		links, _ := job.Run()

		// add error link to another queue
		fmt.Printf("url %s, links \n %v", job.StartUrl, links)
		w.Mx.Lock()
		w.Queue = append(w.Queue, links...)
		w.Mx.Unlock()

		result <- map[string][]string{
			job.StartUrl: links,
		}

		for _, link := range links {
			w.links <- link
		}
		wg.Done()
	}
}

func (w *worker) Run(filters []string, visited map[string]bool) {
	wg := sync.WaitGroup{}
	go w.GenerateJobs(filters, visited)

	for i := 0; i < w.PoolCount; i++ {
		wg.Add(1)
		go w.DoWork(&wg, w.Jobs, w.Results)
	}

	wg.Wait()
	close(w.Jobs)
}

func (w *worker) GetResultChan() chan map[string][]string {
	return w.Results
}

func NewWorker(poolCount int, queue []string, pf ProcessFn) monzo_interview.Worker {
	return &worker{
		PoolCount: poolCount,
		Jobs:      make(chan monzo_interview.Job),
		Results:   make(chan map[string][]string),
		Mx:        &sync.Mutex{},
		Queue:     queue,
		ProcessFn: pf,
	}
}
