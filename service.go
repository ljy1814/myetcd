package myetcd

func NewService(domain, name, version string) *Service {
	return &Service{
		Domain:          domain,
		Name:            name,
		Version:         version,
		Type:            HttpService,
		OnlyLeaderServe: false,
		LBPolicy:        DefaultLBPolicy,
		RetryTimes:      DefaultRetryTimes,
		DialTimeout:     DefaultDialTimeout,
		EndpointTimeout: DefaultEndpointTimeout,

		Endpoints: make(map[string]*Endpoint),
		Msgs:      make(map[string]string),
	}
}

func (s *Service) NewEndpoint(addr string) (ep *Endpoint) {
	ep = newEndpoint(addr)
	s.Endpoints[addr] = ep
	return ep
}

func (s *Service) GetEndpoint(addr string) *Endpoint {
	return s.Endpoints[addr]
}

func newEndpoint(addr string) *Endpoint {
	return &Endpoint{
		Addr:           addr,
		Status:         EndpointStatusNormal,
		FreezeDuration: 0,
	}
}
