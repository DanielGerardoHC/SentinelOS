package model

type Route struct {
	ID          int
	Destination string
	Gateway     string
	Interface   string
	Metric      int
	Description string
}
