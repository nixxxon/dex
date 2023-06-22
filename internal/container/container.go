package container

type Container struct {
	ID    string
	Names []string
	Image string
}

//go:generate mockery --name Service
type Service interface {
	GetAll() ([]Container, error)
	RunCommand(containerID, command string) error
	Close()
}
