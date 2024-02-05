package main

import "fmt"

func buildNetwork(resources []resource, processes []process, g goal) []process {
	var curr []process
	for i := range processes {
		for _, product := range processes[i].products {
			if product.name == g.product {
				processes[i].final = true
				curr = append(curr, processes[i])
				// fmt.Println("Final process: ", processes[i].string())
			}
		}
	}

	// Work back through the processes to build precedence relationships.
	for len(curr) > 0 {
		var next []process
		for _, c := range curr {
			for _, ingredient := range c.ingredients {
				for _, resource := range resources {
					if ingredient.name == resource.name {
						c.initial = true
					}
				}
				for i := range processes {
					for _, product := range processes[i].products {
						if ingredient.name == product.name {
							processes[i].successor = &c
							fmt.Println("Process: ", processes[i].string())
							next = append(next, processes[i])
						}
					}
				}
			}
		}
		curr = next
	}

	return processes
}
