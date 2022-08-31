package main

import (
	"bytes"
	"fmt"
	"github.com/yosssi/gohtml"
	"x-dry-go/src/internal/clone_detect"
	"x-dry-go/src/internal/compare"
	"x-dry-go/src/internal/service/aggregate"
	"x-dry-go/src/templates"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates

func main() {
	// names := []string{"Kate", "Go", "John", "Brad"}

	matches := []compare.Match{
		{
			Content: "blablabla",
			IndexA:  0,
			IndexB:  5,
		},
	}
	type1Clones := []clone_detect.Clone{
		{
			A:       "foo.php",
			B:       "bar.php",
			Matches: matches,
		},
		{
			A:       "foo.php",
			B:       "bar.php",
			Matches: matches,
		},
	}

	clones := map[string][]clone_detect.Clone{
		"1": type1Clones,
		"2": type1Clones,
		"3": type1Clones,
		"4": type1Clones,
	}

	cloneBundles := []aggregate.CloneBundle{
		{
			CloneType: 1,
			AggregatedClones: []aggregate.AggregatedClone{
				{
					Content: "blablabla",
					Instances: []aggregate.CloneInstance{
						{
							Path:  "foo.php",
							Index: 10,
						},
						{
							Path:  "bar.php",
							Index: 0,
						},
					},
				},
				{
					Content: "foooooo",
					Instances: []aggregate.CloneInstance{
						{
							Path:  "foo.php",
							Index: 0,
						},
						{
							Path:  "bum.php",
							Index: 0,
						},
					},
				},
			},
		},
		{
			CloneType: 2,
			AggregatedClones: []aggregate.AggregatedClone{
				{
					Content: "sdfgsdgsdffg",
					Instances: []aggregate.CloneInstance{
						{
							Path:  "foo.php",
							Index: 10,
						},
						{
							Path:  "bar.php",
							Index: 0,
						},
					},
				},
				{
					Content: "xxxx1x1x1x1",
					Instances: []aggregate.CloneInstance{
						{
							Path:  "foo.php",
							Index: 0,
						},
						{
							Path:  "bum.php",
							Index: 0,
						},
					},
				},
			},
		},
	}

	// qtc creates Write* function for each template function.
	// Such functions accept io.Writer as first parameter:
	var buf bytes.Buffer
	templates.WriteClones(&buf, clones)

	fmt.Println(gohtml.Format(buf.String()))
	//fmt.Printf("%s", buf.Bytes())
}
