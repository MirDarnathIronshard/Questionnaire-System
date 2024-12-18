package services

import (
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/internal/repositories"
	"github.com/QBG-P2/Voting-System/pkg/service_errors"
)

type OptionService struct {
	OptionRepo repositories.OptionRepository
}

func NewOptionService(optionRepo repositories.OptionRepository) *OptionService {
	return &OptionService{
		OptionRepo: optionRepo,
	}
}

func (s *OptionService) CreateOption(option *models.Option) (*models.Option, error) {
	if err := option.Validate(); err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage:   service_errors.UnExpectedError,
			TechnicalMessage: "Invalid option data validation",
			Err:              err,
		}
	}

	err := s.OptionRepo.Create(option)

	if err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage:   service_errors.UnExpectedError,
			TechnicalMessage: "Error saving option to database",
			Err:              err,
		}
	}

	return option, nil
}

func (s *OptionService) GetOptionByID(id uint) (*models.Option, error) {
	option, err := s.OptionRepo.GetByID(id)
	if err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage:   service_errors.RecordNotFound,
			TechnicalMessage: "Option not found with provided ID",
			Err:              err,
		}
	}

	return option, nil
}

func (s *OptionService) UpdateOption(id uint, updatedOption *models.Option) (*models.Option, error) {

	if err := updatedOption.Validate(); err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage:   service_errors.UnExpectedError,
			TechnicalMessage: "Invalid option data validation",
			Err:              err,
		}
	}

	_, err := s.OptionRepo.GetByID(id)
	if err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage:   service_errors.RecordNotFound,
			TechnicalMessage: "Option not found for update",
			Err:              err,
		}
	}

	err = s.OptionRepo.Update(updatedOption)

	if err != nil {
		return nil, &service_errors.ServiceError{
			EndUserMessage:   service_errors.UnExpectedError,
			TechnicalMessage: "Error updating option in database",
			Err:              err,
		}
	}

	return updatedOption, nil
}

func (s *OptionService) DeleteOption(id uint) error {
	err := s.OptionRepo.Delete(id)
	if err != nil {
		return &service_errors.ServiceError{
			EndUserMessage:   service_errors.UnExpectedError,
			TechnicalMessage: "Error deleting option from database",
			Err:              err,
		}
	}
	return nil
}
