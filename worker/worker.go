package worker

import (
	monzo_interview "github.com/Mutusva/monzo-webcrawler"
	"log"
	"sync"
)

type ProcessFn func(curUrl string, filters []string) ([]string, error)

type worker struct {
	PoolCount int
	Results   chan map[string][]string
	links     chan string
	Mx        *sync.Mutex
	ProcessFn ProcessFn
	Queue     []string
}

func (w *worker) GenerateJobs(filters []string, visited map[string]bool, wg *sync.WaitGroup) {

	defer wg.Done()
	w.Mx.Lock()
	count := len(w.Queue)
	w.Mx.Unlock()

	for count > 0 {

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

		links, err := job.Run()
		if err != nil {
			log.Printf("error on url %s %v", curUrl, err)
		}
		w.Mx.Lock()
		w.Queue = append(w.Queue, links...)
		w.Mx.Unlock()
		w.Results <- map[string][]string{
			job.StartUrl: links,
		}
		visited[curUrl] = true
	}
}

func (w *worker) Run(filters []string, visited map[string]bool) {
	wg := sync.WaitGroup{}
	for i := 0; i < w.PoolCount; i++ {
		wg.Add(1)
		go w.GenerateJobs(filters, visited, &wg)
	}

	wg.Wait()
	close(w.Results)
}

func (w *worker) GetResultChan() chan map[string][]string {
	return w.Results
}

func NewWorker(poolCount int, queue []string, pf ProcessFn) monzo_interview.Worker {
	return &worker{
		PoolCount: poolCount,
		Results:   make(chan map[string][]string),
		Mx:        &sync.Mutex{},
		Queue:     queue,
		ProcessFn: pf,
	}
}
