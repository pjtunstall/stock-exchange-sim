package main

func buildNetwork(resources []resource, processes []process, g goal) {
	var curr []*process

	// Maximum number of predecessors for a process. Adjust as needed.
	n := 64

	for i := range processes {
		processes[i].predecessors = make([]*process, 0, n)
	}

	// Find the processes that produce the goal product.
	// Can't use range here because we're modifying elements of the processes
	// slice while iterating over it. Range would create a copy of the slice.
	for i := range processes {
		for _, product := range processes[i].products {
			if product.name == g.product {
				processes[i].final = true
				curr = append(curr, &processes[i])
			}
		}
		for _, ingredient := range processes[i].ingredients {
			for _, resource := range resources {
				if ingredient.name == resource.name {
					processes[i].initial = true
					continue
				}
			}
		}
	}

	// Work back through the processes to build precedence relationships.
	for len(curr) > 0 {
		var next []*process
		for k := range curr {
			for _, ingredient := range curr[k].ingredients {
				// Can't use range here because we're modifying elements of the processes
				// slice while iterating over it. Range would create a copy of the slice.
				for i := range processes {
					for _, product := range processes[i].products {
						if ingredient.name == product.name {
							processes[i].successor = curr[k]
							curr[k].predecessors = append(curr[k].predecessors, &processes[i])
							next = append(next, &processes[i])
						}
					}
				}
			}
		}
		curr = next
	}
}
