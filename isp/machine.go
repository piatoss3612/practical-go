package isp

import "fmt"

type Document struct {
	Content string
}

type Machine interface {
	Print(d Document)
	Fax(d Document)
	Scan(d Document)
}

type OldFashionedPrinter struct{}

func (o OldFashionedPrinter) Print(d Document) {
	fmt.Printf("Print: %s\n", d.Content)
}

func (o OldFashionedPrinter) Fax(d Document) {
	panic("operation not supported")
}

func (o OldFashionedPrinter) Scan(d Document) {
	panic("operation not supported")
}

var _ Machine = (*OldFashionedPrinter)(nil)

type Printer interface {
	Print(d Document)
}

type Scanner interface {
	Scan(d Document)
}

type Faxer interface {
	Fax(d Document)
}

type CustomPrinter struct{}

func (c CustomPrinter) Print(d Document) {
	fmt.Printf("Print: %s\n", d.Content)
}

func (c CustomPrinter) Scan(d Document) {
	fmt.Printf("Scan: %s\n", d.Content)
}

var _ Printer = (*CustomPrinter)(nil)
var _ Scanner = (*CustomPrinter)(nil)

type CustomMachine interface {
	Printer
	Scanner
}

var _ CustomMachine = (*CustomPrinter)(nil)
