package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	c := make(chan struct{})
	setTimer(checkArgs(), c)

	resources, processes, goal, err := parseFile("./" + os.Args[1])
	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	// `ubik` is the the name of a resource that is produced by every process.
	finite, ubik := buildGraph(resources, processes, goal)
	currentResources := make(map[string]int)
	for i := range resources {
		currentResources[resources[i].name] = resources[i].quantity
	}

	end, t := schedule(currentResources, processes, goal, finite, c, ubik)
	// printConfig(resources, processes, goal)
	s := buildOutput(currentResources, processes, end, finite, t)
	fmt.Println()
	fmt.Println(s)
	// writeOutput(s)
}

func checkArgs() float64 {
	// Default waiting time: 1 second
	f := 1.0
	if len(os.Args) > 3 {
		log.Println("too many arguments")
		log.Fatal("usage: ./stock <file> <waiting_time>")
	}
	if len(os.Args) < 2 {
		log.Println("not enough arguments")
		log.Fatal("usage: ./stock <file> <waiting_time>")
	}
	if len(os.Args) == 3 {
		g, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			log.Fatalf("error parsing wait time: %v", err)
		}
		f = g
	}
	return f
}

// func printConfig(resources []resource, processes []process, goal goal) {
// 	fmt.Println("\nResources:")
// 	for _, r := range resources {
// 		fmt.Println(r.string())
// 	}
// 	fmt.Println()

// 	fmt.Println("\nProcesses:")
// 	for _, p := range processes {
// 		fmt.Println(p.string())
// 	}

// 	fmt.Println("\nOptimize:")
// 	fmt.Println(goal.string())
// 	fmt.Println()
// }
