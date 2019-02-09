package testdata

type (
	methodValid struct{}
)

func (methodValid) MethodValid1()    {}
func (m *methodValid) MethodValid2() {}

func (methodValid) methodValid1()    {}
func (m *methodValid) methodValid2() {}
