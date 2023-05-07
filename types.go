package monzo_interview

type Crawler interface {
	Start(bool) error
}

type Worker interface {
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
