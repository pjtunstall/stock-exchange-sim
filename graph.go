package main

func buildGraph(resources []resource, processes []process, g goal) bool {
	var curr []*process

	// Maximum number of predecessors for a process. Adjust as needed.
	n := 64

	// This will be used to assign activity numbers to the processes
	// in such a way that successors always have greater numbers than
	// their predecessors. This will allow us to have a default order
	// of precedence for the processes to fall back on when there's
	// nothing else to choose between them.
	activityNumber := len(processes)

	for i := range processes {
		processes[i].predecessors = make([]*process, 0, n)
		processes[i].start = -1
	}

	// Find the processes that produce the goal product.
	// Can't use _, range here because we want to modify elements of the processes
	// slice while iterating over it. Range would create a copy of the slice.
	for i := range processes {
		for _, product := range processes[i].products {
			if product.name == g.product {
				processes[i].added = true
				processes[i].final = true
				processes[i].minCount = rational{numerator: 1, denominator: 1}
				processes[i].activityNumber = activityNumber
				activityNumber--
				curr = append(curr, &processes[i])
			}
		}
		for _, ingredient := range processes[i].ingredients {
			for _, resource := range resources {
				if ingredient.name == resource.name && resource.quantity > 0 && resource.quantity >= ingredient.quantity {
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
				for i := range processes {
					for _, product := range processes[i].products {
						if ingredient.name == product.name {
							if processes[i].added {
								return false
							}
							processes[i].added = true
							processes[i].successor = curr[k]
							r := rational{ingredient.quantity, 1}.Times(curr[k].minCount)
							processes[i].minCount = r.Times(rational{1, product.quantity})
							processes[i].activityNumber = activityNumber
							activityNumber--
							curr[k].predecessors = append(curr[k].predecessors, &processes[i])
							next = append(next, &processes[i])
						}
					}
				}
			}
		}
		curr = next
	}
	return true
}
