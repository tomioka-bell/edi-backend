package ports

import "backend/internal/core/models"

type QuestionnaireService interface {
	CreateQuestionnaireService(req models.QuestionnaireResp) error
	GetQuestionnaires(limit, offset int) ([]models.QuestionnaireReq, error)
}
