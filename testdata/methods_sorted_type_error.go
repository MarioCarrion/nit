package testdata

type (
	MethodSortedTypeA struct{}
	MethodSortedTypeB struct{}
)

func (MethodSortedTypeB) F() {}
func (MethodSortedTypeB) G() {}

//-

func (MethodSortedTypeA) A() {}
