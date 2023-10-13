package main

import "solid/lsp"

func main() {
	r := lsp.NewRectangle(2, 3)
	s := lsp.NewSquare(2)
	lsp.PrintArea(r)
	lsp.PrintArea(s)
}
