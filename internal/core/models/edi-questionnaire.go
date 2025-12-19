package models

type QuestionnaireReq struct {
	QuestionnaireID     string `json:"questionnaire_id"`
	QuestionnaireName   string `json:"questionnaire_name"`
	QuestionnairePhone  string `json:"questionnaire_phone"`
	QuestionnaireEmail  string `json:"questionnaire_email"`
	Question            string `json:"question"`
	QuestionnaireType   string `json:"questionnaire_type"`
	QuestionnaireStatus string `json:"questionnaire_status"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type QuestionnaireResp struct {
	QuestionnaireID     string `json:"questionnaire_id"`
	QuestionnaireName   string `json:"questionnaire_name"`
	QuestionnairePhone  string `json:"questionnaire_phone"`
	QuestionnaireEmail  string `json:"questionnaire_email"`
	Question            string `json:"question"`
	QuestionnaireType   string `json:"questionnaire_type"`
	QuestionnaireStatus string `json:"questionnaire_status"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}
