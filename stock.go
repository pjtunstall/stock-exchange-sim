package main

import (
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
	currentResources := make(map[string]int)
	for i := range resources {
		currentResources[resources[i].name] = resources[i].quantity
	}

	// // Work in progress: count number of times each process must be run.
	end := count(currentResources, processes, goal)
	_ = buildOutput(currentResources, processes, end)
	// writeOutput(s)

	// fmt.Println("\nResources:")
	// for _, r := range resources {
	// 	fmt.Println(r.string())
	// }
	// fmt.Println()

	// fmt.Println("\nProcesses:")
	// for _, p := range processes {
	// 	fmt.Println(p.string())
	// }

	// fmt.Println("\nOptimize:")
	// fmt.Println(goal.string())
	// fmt.Println()

	// for resource, quantity := range currentResources {
	// 	fmt.Println(resource, quantity)
	// }
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
