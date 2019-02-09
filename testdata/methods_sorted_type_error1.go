package testdata

type (
	MethodSortedTypeA1 struct{}
	MethodSortedTypeB1 struct{}
)

func (MethodSortedTypeB1) F() {}
func (MethodSortedTypeB1) G() {}

func (MethodSortedTypeA1) A() {}
