package grpc

type ServiceProtocol interface {
}

type ServiceGRPC struct {
}

func NewServiceGRPC() *ServiceGRPC {
	return &ServiceGRPC{}
}

func (s *ServiceGRPC) RRR() {

}
