package domains

import "time"

type Questionnaire struct {
	QuestionnaireID     string    `gorm:"type:uniqueidentifier;primaryKey;default:gen_random_uuid()" json:"questionnaire_id"`
	QuestionnaireName   string    `json:"questionnaire_name"`
	QuestionnairePhone  string    `json:"questionnaire_phone"`
	QuestionnaireEmail  string    `json:"questionnaire_email"`
	Question            string    `json:"question"`
	QuestionnaireType   string    `json:"questionnaire_type"`
	QuestionnaireStatus string    `json:"questionnaire_status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
