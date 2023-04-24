package errorx

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	e := C()
	fmt.Println(e)
	fmt.Println(e)
}

func A() error {
	return B()
}

func B() error {
	return C()
}

func C() error {
	return InternalServer("internal server error").
		WithError(fmt.Errorf("db connection error")).
		WithStack()
}
