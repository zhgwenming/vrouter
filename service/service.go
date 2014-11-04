package service

// A TCP service will publish to outer network
type Service struct {
	Name   string
	IpAddr string
	Port   string
}

type LBProxy struct {
}
