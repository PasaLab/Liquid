package main

type State int

const (
	// submitted
	Created State = iota
	// scheduling
	Starting
	// running
	Running
	// stopped
	Stopped
	// finished successfully
	Finished

	Failed
)
