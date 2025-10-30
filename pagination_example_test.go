package cursor

import "fmt"

func ExamplePaginate() {
	c := New[Int64](50, 0)
	c.Prev = new(Int64)
	p, _ := Paginate(c, []byte("qsddqdqsdqs"))
	fmt.Printf("%#v", p)
	// Output: rv
}

// c := cursor.New[int](10)
//
//	var c Cursor[int]
//	c.
//		fmt.Println("rv", c.Prev())
//
