package handlers

import (
	"backend/internal/core/models"
	services "backend/internal/core/ports/services"
	"log"

	"github.com/gofiber/fiber/v2"
)

type QuestionnaireHandler struct {
	QuestionnaireSrv services.QuestionnaireService
}

func NewQuestionnaireHandler(insSrv services.QuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{QuestionnaireSrv: insSrv}
}

func (h *QuestionnaireHandler) CreateQuestionnaireHandler(c *fiber.Ctx) error {
	var req models.QuestionnaireResp

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if h.QuestionnaireSrv == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Service is not available",
		})
	}

	err := h.QuestionnaireSrv.CreateQuestionnaireService(req)
	if err != nil {
		log.Println("Error creating  Questionnaire:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create  Questionnaire",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Questionnaire  created successfully",
	})
}

func (h *QuestionnaireHandler) GetQuestionnairesHandler(c *fiber.Ctx) error {
	limit, offset := c.QueryInt("limit", 10), c.QueryInt("offset", 0)

	jobs, err := h.QuestionnaireSrv.GetQuestionnaires(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve Questionnaires",
		})
	}

	return c.JSON(jobs)
}
