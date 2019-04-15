package testdata

type (
	MethodSortedTypeUnexported_A struct{}
	MethodSortedTypeUnexported_B struct{}
)

func (MethodSortedTypeUnexported_A) A()  {}
func (MethodSortedTypeUnexported_A) B()  {}
func (*MethodSortedTypeUnexported_A) c() {}

//-

func (MethodSortedTypeUnexported_B) A() {}
