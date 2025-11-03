package attachment

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase interface {
	Create(input *AttachmentCreateDTO) (*AttachmentFoundDTO, error)
	FindByID(id uuid.UUID) (*AttachmentFoundDTO, error)
	Update(id uuid.UUID, input *AttachmentUpdateDTO) (*AttachmentFoundDTO, error)
	Delete(id uuid.UUID) error
}

type service struct {
	log        *zap.SugaredLogger
	repository Repository
}

func NewService(
	log *zap.SugaredLogger,
) UseCase {
	return &service{
		log:        log,
		repository: NewRepository(log),
	}
}

func (s *service) Create(input *AttachmentCreateDTO) (*AttachmentFoundDTO, error) {
	attachment := NewAttachment(
		input.URL,
		input.Type,
		input.Description,
	)

	createdAttachment, err := s.repository.Create(attachment)
	if err != nil {
		s.log.Errorw("error creating attachment", "error", err)
		return nil, err
	}

	return &AttachmentFoundDTO{
		ID:          createdAttachment.ID,
		URL:         createdAttachment.URL,
		Type:        createdAttachment.Type,
		Description: createdAttachment.Description,
		CreatedAt:   createdAttachment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   createdAttachment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *service) FindByID(id uuid.UUID) (*AttachmentFoundDTO, error) {
	attachment, err := s.repository.FindByID(id)
	if err != nil {
		s.log.Errorw("error finding attachment by ID", "id", id, "error", err)
		return nil, err
	}
	if attachment == nil {
		s.log.Warnw("attachment not found", "id", id)
		return nil, nil
	}
	return &AttachmentFoundDTO{
		ID:          attachment.ID,
		URL:         attachment.URL,
		Type:        attachment.Type,
		Description: attachment.Description,
		CreatedAt:   attachment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   attachment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *service) Update(id uuid.UUID, input *AttachmentUpdateDTO) (*AttachmentFoundDTO, error) {
	// First check if attachment exists
	existingAttachment, err := s.repository.FindByID(id)
	if err != nil {
		s.log.Errorw("error finding attachment for update", "id", id, "error", err)
		return nil, err
	}
	if existingAttachment == nil {
		s.log.Warnw("attachment not found for update", "id", id)
		return nil, nil
	}

	// Update only provided fields
	updateAttachment := &Attachment{
		URL:         existingAttachment.URL,
		Type:        existingAttachment.Type,
		Description: existingAttachment.Description,
	}

	if input.URL != nil {
		updateAttachment.URL = *input.URL
	}
	if input.Type != nil {
		updateAttachment.Type = input.Type
	}
	if input.Description != nil {
		updateAttachment.Description = input.Description
	}

	updatedAttachment, err := s.repository.Update(id, updateAttachment)
	if err != nil {
		s.log.Errorw("error updating attachment", "id", id, "error", err)
		return nil, err
	}

	return &AttachmentFoundDTO{
		ID:          updatedAttachment.ID,
		URL:         updatedAttachment.URL,
		Type:        updatedAttachment.Type,
		Description: updatedAttachment.Description,
		CreatedAt:   updatedAttachment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedAttachment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *service) Delete(id uuid.UUID) error {
	// Check if attachment exists
	existingAttachment, err := s.repository.FindByID(id)
	if err != nil {
		s.log.Errorw("error finding attachment for deletion", "id", id, "error", err)
		return err
	}
	if existingAttachment == nil {
		s.log.Warnw("attachment not found for deletion", "id", id)
		return nil
	}

	err = s.repository.Delete(id)
	if err != nil {
		s.log.Errorw("error deleting attachment", "id", id, "error", err)
		return err
	}

	s.log.Infow("attachment deleted successfully", "id", id)
	return nil
}
