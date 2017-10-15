package myetcd

type Registry interface {
	Start() error
	Stop()
	CreateService(service *Service) error
	ListServices() ([]*Service, error)
	GetService(domain, service, version string) (*Service, error)
	UpdateService(oldService, newService *Service) error
	DeleteService(domain, service, version string) error

	FindService(domain, service, version string) *Service

	RegisterEndpoint(domain, service, version, addr string, delegate bool) (*Service, error)
	UnregisterEndpoint(domain, service, version, addr string, delegate bool) error
	RefreshEndpoint(domain, service, version, addr string, tiemout uint64) error
	DeleteEndpoint(domain, service, version, addr string) error

	GetEndpoint(domain, serviceName, version string) (string, error)
	GetEndpoints(domain, serviceName, version string) ([]string, error)
}
