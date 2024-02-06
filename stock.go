package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	checkArgs()

	if len(os.Args) == 3 {
		setTimer()
	}

	resources, processes, goal, err := parseFile("./" + os.Args[1])
	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	buildNetwork(resources, processes, goal)

	// // Work in progress: count number of times each process must be run.
	// count(resources, processes, goal)

	fmt.Println("\nResources:")
	for _, r := range resources {
		fmt.Println(r.string())
	}
	fmt.Println()

	fmt.Println("\nProcesses:")
	for _, p := range processes {
		fmt.Println(p.string())
	}

	fmt.Println("\nOptimize:")
	fmt.Println(goal.string())
	fmt.Println()

	// // For testing timer.
	// fmt.Scanln()
}

func checkArgs() {
	if len(os.Args) > 3 {
		log.Println("too many arguments")
		log.Fatal("usage: ./stock <file> <waiting_time>")
	}
	if len(os.Args) < 2 {
		log.Println("not enough arguments")
		log.Fatal("usage: ./stock <file> <waiting_time>")
	}
}
