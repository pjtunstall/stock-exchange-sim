package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func checker(resources map[string]int, processes []process, goal goal) {
	file, err := os.Open("./" + flag.Arg(1))
	if err != nil {
		log.Fatalf("error parsing log file: %v", err)
	}
	defer file.Close()

	p := make(map[string]*process)
	for i := range processes {
		p[processes[i].name] = &processes[i]
	}

	fmt.Println()

	scanner := bufio.NewScanner(file)
	isFirstLine := true
	for scanner.Scan() {
		if isFirstLine {
			isFirstLine = false
			continue
		}

		line := scanner.Text()
		ln := strings.Split(line, ":")
		if len(ln) != 2 {
			break
		}
		cycle := ln[0]
		name := ln[1]
		if _, ok := strconv.Atoi(cycle); ok != nil {
			break
		}

		process, ok := p[name]
		if !ok {
			log.Fatalf("process %s not found", name)
		}

		for _, ingredient := range process.ingredients {
			if resources[ingredient.name] < ingredient.quantity {
				fmt.Println("Error detected")
				fmt.Printf("at %s stock insufficient\n", line)
				fmt.Println("Exiting...")
				fmt.Println()
				os.Exit(0)
			}
			resources[ingredient.name] -= ingredient.quantity
		}
		for _, product := range process.products {
			resources[product.name] += product.quantity
		}
	}
	fmt.Println("Trace completed, no error detected.")
	fmt.Println()
	os.Exit(0)
}
