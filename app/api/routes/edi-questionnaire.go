package api

import (
	"backend/internal/core/services"
	"backend/internal/handlers"
	"backend/internal/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RoutesQuestionnaire(db *gorm.DB) *fiber.App {
	if db == nil {
		panic("Database connection is nil")
	}

	app := fiber.New()

	QuestionnaireRepository := repositories.NewQuestionnaireRepositoryDB(db)
	QuestionnaireService := services.NewQuestionnaireService(QuestionnaireRepository)
	QuestionnaireHandler := handlers.NewQuestionnaireHandler(QuestionnaireService)

	app.Post("/create-questionnaire", QuestionnaireHandler.CreateQuestionnaireHandler)
	app.Get("/get-questionnaires", QuestionnaireHandler.GetQuestionnairesHandler)
	return app
}
