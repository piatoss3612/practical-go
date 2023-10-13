package lsp

type Sized interface {
	Width() int
	SetWidth(width int)
	Height() int
	SetHeight(height int)
}

// func PrintArea(s Sized) {
// 	s.SetWidth(3)
// 	s.SetHeight(4)
// 	println(s.Width() * s.Height())
// }

// func PrintArea(s Sized) {
// 	println(s.Width() * s.Height())
// }

type Shape interface {
	Area() int
}

func PrintArea(s Shape) {
	println(s.Area())
}

type Rectangle struct {
	width, height int
}

func NewRectangle(width, height int) *Rectangle {
	return &Rectangle{width, height}
}

func (r *Rectangle) Width() int {
	return r.width
}

func (r *Rectangle) SetWidth(width int) {
	r.width = width
}

func (r *Rectangle) Height() int {
	return r.height
}

func (r *Rectangle) SetHeight(height int) {
	r.height = height
}

func (r *Rectangle) Area() int {
	return r.width * r.height
}

type Square struct {
	size int
}

func NewSquare(size int) *Square {
	return &Square{size}
}

func (s *Square) Size() int {
	return s.size
}

func (s *Square) SetSize(size int) {
	s.size = size
}

func (s *Square) Area() int {
	return s.size * s.size
}

func (s *Square) Rectangle() *Rectangle {
	return NewRectangle(s.size, s.size)
}
