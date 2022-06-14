package main

import "fmt"

func main() {
	p := post{
		[]tag{
			{
				name: "go",
			},
			{
				name: "programming",
			},
			{
				name: "range",
			},
		},
	}

	// the object returned by range is a copy of
	// the item of the slice. Therefore any modification is not kept
	for _, tag := range p.tags {
		tag.name = "changed"
	}

	fmt.Println(p)
	// prints: {[{go} {programming} {loop}]}

	// To modify the item, we need to access its reference via the index of the array/slice
	for i := range p.tags {
		p.tags[i].name = "changed"
	}

	fmt.Println(p)
	// prints: {[{changed} {changed} {changed}]}
}

type post struct {
	tags []tag
}

type tag struct {
	name string
}
