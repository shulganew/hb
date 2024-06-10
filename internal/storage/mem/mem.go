package mem

type Memory struct {
	logedIn map[string]struct{} // Autentificated telegram users.
}

func NewMemory() *Memory {
	ln := make(map[string]struct{}, 0)

	mem := &Memory{logedIn: ln}
	return mem
}

func (m *Memory) Add(user string) {
	m.logedIn[user] = struct{}{}
}

func (m *Memory) Get(user string) bool {
	_, ok := m.logedIn[user]
	return ok
}
