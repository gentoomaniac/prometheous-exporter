package types

// https://golang.org/pkg/plugin/
type Metric struct {
	Name    string
	Value   int
	Help    string
	Labels  map[string]string
	Comment string
}
