package container

type Container struct {
	ID    string
	Names []string
	Image string
}

//go:generate mockery --name ContainerService
type ContainerService interface {
	GetAll() ([]Container, error)
	RunCommand(containerID, command string) error
	Close()
}
