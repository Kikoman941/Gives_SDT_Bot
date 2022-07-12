package member

import (
	"Gives_SDT_Bot/internal/data"
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
		MemberTgId: memberTgId,
	}

	if err := s.repository.Create(context.TODO(), member); err != nil {
		if err != data.ERROR_MEMBER_ALREADY_EXIST {
			s.logger.Errorf("cannot save give member giveId=%d MemberTgId=%s: %s", giveId, memberTgId, err)
		}
		return err
	}

	return nil
}

func (s *Service) GetRandomMembersByGiveId(giveId int, count int) ([]string, error) {
	var members []string
	mbs, err := s.repository.FindRandomLimitWithConditions(context.TODO(), count, `"giveId"=?`, giveId)
	if err != nil {
		s.logger.Errorf("cannot get give members giveId=%d: %s", giveId, err)
		return nil, err
	} else if len(mbs) == 0 {
		s.logger.Infof("not found give members giveId=%d", giveId)
		return nil, nil
	}

	for _, member := range mbs {
		members = append(members, member.MemberTgId)
	}

	return members, nil
}
