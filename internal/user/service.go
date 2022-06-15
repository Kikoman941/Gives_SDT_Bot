package user

type Service struct {
	repository Repository
}

func NewUserService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
