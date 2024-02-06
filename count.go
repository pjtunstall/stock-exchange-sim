package main

// Work in progress: count number of times each process must be run.
func count(resources []resource, processes []process, goal goal) {
	curr := make([]*process, 0, 64)

	// This map will be used to keep track of the remaining resources
	// at each step, starting with initial resources, and at each step
	// containing the products of the previous step.
	rMap := make(map[string]int)
	for i := range resources {
		rMap[resources[i].name] = resources[i].quantity
	}

	// Seed the curr slice with the initial processes.
	for i := range processes {
		if processes[i].initial {
			curr = append(curr, &processes[i])
		}
	}

	// While there are processes to be run, decrement the resources
	// and add the products to the next slice of processes to be run.
	for len(curr) > 0 {
		next := make([]*process, 0, 64)
		for len(rMap) > 0 {
			for k := range curr {
				for j := range curr[k].ingredients {
					rMap[curr[k].ingredients[j].name] -= curr[k].ingredients[j].quantity
					if rMap[curr[k].ingredients[j].name] <= 0 {
						delete(rMap, curr[k].ingredients[j].name)
					}
				}
				d := curr[k].minCount.Denominator
				if d > 1 {
					curr[k].count *= d
				}
				next = append(next, curr[k].successor)
			}
		}
		curr = next
	}
}
