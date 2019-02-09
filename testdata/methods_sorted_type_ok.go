package testdata

type (
	MethodSortedTypeOK_A struct{}
	MethodSortedTypeOK_B struct{}
)

func (MethodSortedTypeOK_A) F() {}
func (MethodSortedTypeOK_A) G() {}

func (MethodSortedTypeOK_A) a() {}
func (MethodSortedTypeOK_A) b() {}

//-

func (MethodSortedTypeOK_B) a() {}
func (MethodSortedTypeOK_B) b() {}
