package module

type Module struct {
	Name        string
	Initialize  func()
	IsInstalled func() bool
}

func (m *Module) Run() {
	for {
	}
}
