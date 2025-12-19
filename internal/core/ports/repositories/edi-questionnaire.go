package ports

import "backend/internal/core/domains"

type QuestionnaireRepository interface {
	CreateQuestionnaireRepository(User *domains.Questionnaire) error
	GetQuestionnaireByID(questionnaireID string) (domains.Questionnaire, error)
	GetAllQuestionnaire() ([]domains.Questionnaire, error)
	GetQuestionnaire(limit, offset int) ([]domains.Questionnaire, error)
	UpdateQuestionnaireWithMap(questionnaireID string, updates map[string]interface{}) error
	GetQuestionnaireCount() (int64, error)
}
