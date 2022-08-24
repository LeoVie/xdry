package main

import (
	"bytes"
	"fmt"
	"github.com/yosssi/gohtml"
	"x-dry-go/src/templates"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=templates

func main() {
	names := []string{"Kate", "Go", "John", "Brad"}

	// qtc creates Write* function for each template function.
	// Such functions accept io.Writer as first parameter:
	var buf bytes.Buffer
	templates.WriteGreetings(&buf, names)

	fmt.Println(gohtml.Format(buf.String()))
	//fmt.Printf("%s", buf.Bytes())
}
