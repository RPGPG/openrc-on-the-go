package structures

type Service struct {
	Name    string
	Started bool
}

type JsonOutput struct {
	Started  int
	Stopped  int
	Services []Service
}
