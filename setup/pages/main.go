package main

import "flag"

func main() {
	var project string
	var workers int
	flag.StringVar(&project, "p", "simplewiki", "Database name for the project")
	flag.IntVar(&workers, "w", 10, "Number of workers")
	flag.Parse()

}
