package exception

import "fmt"

func ExampleNullReference() {
	var e *NullReference = NewNullReference("a null reference exception", nil)
	fmt.Println(e.Message())
	// Output: FIXME
}
