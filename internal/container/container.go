package container

// Container holds container information
type Container struct {
	ID    string
	Names []string
	Image string
}

// Service interface for getting containers and running commands against them
//
//go:generate mockery --name Service
type Service interface {
	GetAll() ([]Container, error)
	RunCommand(containerID, command string) error
	Close()
}
