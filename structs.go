package main

import (
	"fmt"
	"strings"
)

type goal struct {
	product string
	time    bool
}

type resource struct {
	name     string
	quantity int
}

type process struct {
	name           string
	ingredients    []resource
	products       []resource
	time           int
	successor      *process
	predecessors   []*process
	initial        bool
	final          bool
	minCount       rational
	count          int
	activityNumber int
}

type rational struct {
	Numerator   int
	Denominator int
}

func (r rational) Plus(other rational) rational {
	return rational{
		Numerator:   r.Numerator*other.Denominator + other.Numerator*r.Denominator,
		Denominator: r.Denominator * other.Denominator,
	}
}

func (r rational) Times(other rational) rational {
	return rational{
		Numerator:   r.Numerator * other.Numerator,
		Denominator: r.Denominator * other.Denominator,
	}
}

func (r rational) Simplify() rational {
	gcd := gcd(r.Numerator, r.Denominator)
	return rational{
		Numerator:   r.Numerator / gcd,
		Denominator: r.Denominator / gcd,
	}
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

// func lcm(a []int) int {
// 	if len(a) == 0 {
// 		return 0
// 	}
// 	m := a[0]
// 	for _, n := range a[1:] {
// 		m = m * n / gcd(m, n)
// 	}
// 	return m
// }

func (g goal) string() string {
	return fmt.Sprintf("PRODUCT: %s, TIME: %v", g.product, g.time)
}

func (r resource) string() string {
	return fmt.Sprintf("NAME: %s, QUANTITY: %d", r.name, r.quantity)
}

func (p process) string() string {
	ingredients := make([]string, len(p.ingredients))
	for i, ingredient := range p.ingredients {
		ingredients[i] = fmt.Sprintf("%s: %d", ingredient.name, ingredient.quantity)
	}

	products := make([]string, len(p.products))
	for i, product := range p.products {
		products[i] = fmt.Sprintf("%s: %d", product.name, product.quantity)
	}

	var sucessor string
	if p.successor != nil {
		sucessor = p.successor.name
	} else if p.final {
		sucessor = "none"
	}

	var predecessors string
	if len(p.predecessors) == 0 {
		predecessors = "none"
	}
	for _, predecessor := range p.predecessors {
		predecessors += predecessor.name + ", "
	}
	predecessors = strings.TrimSuffix(predecessors, ", ")

	result := fmt.Sprintf("NAME: %s, INGREDIENTS: %s, PRODUCTS: %s, TIME: %d,\nSUCCESOR: %s, PREDECESSORS: %s,",
		p.name, strings.Join(ingredients, ", "), strings.Join(products, ", "), p.time, sucessor, predecessors)
	result += fmt.Sprintf("\nMINCOUNT: %d, COUNT: %d, ACTIVITY NUMBER: %d\n", p.minCount, p.count, p.activityNumber)

	return result
}

// func (p process) strength() int {
// 	n := p.count
// 	d := 0
// 	for _, ingredient := range p.ingredients {
// 		d += ingredient.quantity
// 	}
// 	strength := n / d
// 	return strength
// }
