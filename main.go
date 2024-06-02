package main

import (
	"flag"
	"log"
	"strconv"
)

func main() {
	checkerFlag := flag.Bool("checker", false, "run the checker")
	flag.Parse()

	c := make(chan struct{})
	setTimer(checkArgs(*checkerFlag), c)

	resources, processes, goal, err := parseFile("./" + flag.Arg(0))
	if err != nil {
		log.Fatalf("error parsing file: %v", err)
	}

	// `ubik` is the the name of a resource that is produced by every process.
	finite, ubik := buildGraph(resources, processes, goal)
	currentResources := make(map[string]int)
	for i := range resources {
		currentResources[resources[i].name] = resources[i].quantity
	}

	if *checkerFlag {
		checker(currentResources, processes, goal)
	}

	end, t := schedule(currentResources, processes, goal, finite, c, ubik)

	s := buildOutput(currentResources, processes, end, finite, t)

	// Uncomment to print the output to the console, but be warned that very
	// large outputs will take a long time to print. This time is not limited
	// by the time argument. With cyclic examples (renewable resources), and
	// large outputs, you're better off writing the output to a file.
	// fmt.Println(s)

	writeOutput(s)
}

func checkArgs(ch bool) float64 {
	// Default waiting time: 1 second
	usage := "usage: ./stock <file> <waiting_time>\n./stock -checker <file> <log_file>"
	f := 1.0
	if flag.NArg() > 2 {
		log.Println("too many arguments")
		log.Fatal(usage)
	}
	if flag.NArg() < 1 {
		log.Println("not enough arguments")
		log.Fatal(usage)
	}
	if flag.NArg() == 2 {
		if ch {
			return f
		}
		g, err := strconv.ParseFloat(flag.Arg(1), 64)
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
