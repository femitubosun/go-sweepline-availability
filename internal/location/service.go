package location

import "context"

type Service interface {
	GetInfo(ctx context.Context) (Location, error)
}

type svc struct {
}

func NewService() Service {
	return &svc{}
}

func (s *svc) GetInfo(ctx context.Context) (Location, error) {
	return NewDefaultLocation(), nil
}
