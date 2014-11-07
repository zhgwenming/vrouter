package service

// Service Instance
type Instance struct {
	Addr string
	Port string
}

// A TCP service will publish to outer network
type Service struct {
	Name    string
	Addr    string
	Port    string
	Targets []Instance
}

func NewService() *Service {
	srv := new(Service)
	tgt := make([]Instance, 4)
	srv.Targets = tgt
	return srv
}

type LBProxy struct {
}
