package main

type tasks struct {
	entries []task
	index   int
}

type task struct {
	done    bool
	details string
}
