package give

import (
	"Gives_SDT_Bot/pkg/logging"
	"context"
	"fmt"
)

type Service struct {
	repository Repository
	logger     *logging.Logger
}

func NewGiveService(repository Repository, logger *logging.Logger) *Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

func (s *Service) CreateGive(channel string, ownerId int) (int, error) {
	give := &Give{
		IsActive: false,
		Owner:    ownerId,
		Channel:  channel,
	}

	if err := s.repository.Create(context.TODO(), give); err != nil {
		s.logger.Errorf("cannot create give with channel %s: %s", channel, err)
	}

	return give.Id, nil
}

func (s *Service) GetAllUserGives(userId int) ([]Give, error) {
	gives, err := s.repository.FindAllWithConditions(context.TODO(), fmt.Sprintf("\"owner\"=%d", userId))
	if err != nil {
		s.logger.Errorf("cannot get userId=%d gives: %s", userId, err)
		return nil, err
	} else if len(gives) == 0 {
		s.logger.Errorf("not found gives for userId=%d", userId)
		return nil, err
	}

	return gives, nil
}

func (s *Service) UpdateGive(giveId int, update string, params ...interface{}) error {
	err := s.repository.UpdateWithConditions(
		context.TODO(),
		fmt.Sprintf("\"id\"=%d", giveId),
		update,
		params...,
	)
	if err != nil {
		s.logger.Errorf("cannot do update giveId=%d: %s", giveId, err)
		return err
	}

	return nil
}

func (s *Service) GetGiveById(giveId int) (Give, error) {
	gives, err := s.repository.FindAllWithConditions(context.TODO(), `"id"=?`, giveId)
	if err != nil {
		s.logger.Errorf("cannot get giveId=%d: %s", giveId, err)
		return Give{}, err
	} else if len(gives) == 0 {
		s.logger.Errorf("not found give giveId=%d", giveId)
	}

	return gives[0], nil
}

func (s *Service) GetGiveByTitle(giveTitle string) (Give, error) {
	gives, err := s.repository.FindAllWithConditions(context.TODO(), `"title"=?`, giveTitle)
	if err != nil {
		s.logger.Errorf("cannot get giveTitle=%s: %s", giveTitle, err)
		return Give{}, err
	} else if len(gives) == 0 {
		s.logger.Errorf("not found give giveTitle=%s", giveTitle)
	}

	return gives[0], nil
}

func (s *Service) CheckFilling(give *Give) []string {
	var unfilledFields []string

	if give.Title == "" || give.Title == " " {
		unfilledFields = append(unfilledFields, "Заголовок")
	} else if give.Description == "" || give.Description == " " {
		unfilledFields = append(unfilledFields, "Описание")
	} else if give.Image == "" {
		unfilledFields = append(unfilledFields, "Картинка")
	} else if give.StartAt.IsZero() {
		unfilledFields = append(unfilledFields, "Время старта")
	} else if give.FinishAt.IsZero() {
		unfilledFields = append(unfilledFields, "Время финиша")
	} else if give.WinnersCount == 0 {
		unfilledFields = append(unfilledFields, "Колличество побежителей")
	} else if give.Channel == "" || give.Channel == " " {
		unfilledFields = append(unfilledFields, "Канал проведения")
	} else if len(give.TargetChannels) == 0 {
		unfilledFields = append(unfilledFields, "Каналы для подписки")
	}

	return unfilledFields
}
