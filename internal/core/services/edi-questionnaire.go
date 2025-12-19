package services

import (
	"backend/internal/core/domains"
	"backend/internal/core/models"
	ports "backend/internal/core/ports/repositories"
	servicesports "backend/internal/core/ports/services"
	"backend/internal/pkgs/logs"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type QuestionnaireService struct {
	questionnaireRepo ports.QuestionnaireRepository
}

func NewQuestionnaireService(questionnaireRepo ports.QuestionnaireRepository) servicesports.QuestionnaireService {
	return &QuestionnaireService{questionnaireRepo: questionnaireRepo}
}

func (s *QuestionnaireService) CreateQuestionnaireService(req models.QuestionnaireResp) error {
	newID := uuid.New()

	domainISR := domains.Questionnaire{
		QuestionnaireID:     newID.String(),
		QuestionnaireName:   req.QuestionnaireName,
		QuestionnairePhone:  req.QuestionnairePhone,
		QuestionnaireEmail:  req.QuestionnaireEmail,
		Question:            req.Question,
		QuestionnaireType:   req.QuestionnaireType,
		QuestionnaireStatus: req.QuestionnaireStatus,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if s.questionnaireRepo == nil {
		log.Println("UserRepo is nil")
		return fmt.Errorf("user repository is not initialized")
	}

	err := s.questionnaireRepo.CreateQuestionnaireRepository(&domainISR)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (s *QuestionnaireService) GetQuestionnaires(limit, offset int) ([]models.QuestionnaireReq, error) {
	var jobs []models.QuestionnaireReq

	query, err := s.questionnaireRepo.GetQuestionnaire(limit, offset)
	if err != nil {
		return nil, err
	}

	for _, job := range query {
		jobReq := models.QuestionnaireReq{
			QuestionnaireID:     job.QuestionnaireID,
			QuestionnaireName:   job.QuestionnaireName,
			QuestionnairePhone:  job.QuestionnairePhone,
			QuestionnaireEmail:  job.QuestionnaireEmail,
			Question:            job.Question,
			QuestionnaireType:   job.QuestionnaireType,
			QuestionnaireStatus: job.QuestionnaireStatus,
			CreatedAt:           job.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:           job.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		jobs = append(jobs, jobReq)
	}

	return jobs, nil
}
