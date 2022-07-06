package member

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewMemberService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) SaveGiveMember(giveId int, memberTgId string) error {
	member := &Member{
		GiveId:     giveId,
		memberTgId: memberTgId,
	}

	if err := s.repository.Create(context.TODO(), member); err != nil {
		s.logger.Errorf("cannot save give member giveId=%d memberTgId=%s: %s", giveId, memberTgId, err)
		return err
	}

	return nil
}
