package repositories

import (
	"backend/internal/core/domains"
	ports "backend/internal/core/ports/repositories"
	"fmt"

	"gorm.io/gorm"
)

type QuestionnaireRepositoryDB struct {
	db *gorm.DB
}

func NewQuestionnaireRepositoryDB(db *gorm.DB) ports.QuestionnaireRepository {
	// if err := db.AutoMigrate(&domains.Questionnaire{}); err != nil {
	// 	fmt.Printf("failed to auto migrate: %v", err)
	// }
	return &QuestionnaireRepositoryDB{db: db}
}

func (r *QuestionnaireRepositoryDB) CreateQuestionnaireRepository(User *domains.Questionnaire) error {
	if err := r.db.Create(User).Error; err != nil {
		fmt.Printf("CreateUserRepository error: %v\n", err)
		return err
	}
	return nil
}

func (r *QuestionnaireRepositoryDB) GetQuestionnaireByID(questionnaireID string) (domains.Questionnaire, error) {
	var user domains.Questionnaire
	if err := r.db.Where("questionnaire_id = ?", questionnaireID).First(&user).Error; err != nil {
		return domains.Questionnaire{}, err
	}
	return user, nil
}

func (r *QuestionnaireRepositoryDB) GetAllQuestionnaire() ([]domains.Questionnaire, error) {
	var reviews []domains.Questionnaire
	return reviews, r.db.Find(&reviews).Error
}

func (r *QuestionnaireRepositoryDB) GetQuestionnaire(limit, offset int) ([]domains.Questionnaire, error) {
	var jobs []domains.Questionnaire

	query := r.db.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC")

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *QuestionnaireRepositoryDB) UpdateQuestionnaireWithMap(questionnaireID string, updates map[string]interface{}) error {
	return r.db.Model(&domains.Questionnaire{}).
		Where("questionnaire_id = ?", questionnaireID).
		Updates(updates).
		Error
}

func (r *QuestionnaireRepositoryDB) GetQuestionnaireCount() (int64, error) {
	var count int64
	return count, r.db.Model(&domains.Questionnaire{}).Count(&count).Error
}
