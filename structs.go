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
	name        string
	ingredients []resource
	products    []resource
	time        int
	successor   *process
	initial     bool
	final       bool
}

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
		sucessor = "nil"
	}

	return fmt.Sprintf("NAME: %s, INGREDIENTS: %s, PRODUCTS: %s, TIME: %d, SUCCESOR: %s",
		p.name, strings.Join(ingredients, ", "), strings.Join(products, ", "), p.time, sucessor)
}
