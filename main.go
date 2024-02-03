package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var resources []resource
	var processes []process
	var goal goal

	resources, processes, goal, err := parseFile("./" + os.Args[1])
	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	fmt.Println("\nResources:")
	for _, r := range resources {
		fmt.Println(r.string())
	}

	fmt.Println("\nProcesses:")
	for _, p := range processes {
		fmt.Println(p.string())
	}

	fmt.Println("\nOptimize:")
	fmt.Println(goal.string())
	fmt.Println()
}
