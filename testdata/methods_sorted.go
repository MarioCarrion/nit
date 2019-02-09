package testdata

type (
	MethodSorted struct{}
)

func (MethodSorted) B() {}
func (MethodSorted) C() {}

//-

func (MethodSorted) A() {}
