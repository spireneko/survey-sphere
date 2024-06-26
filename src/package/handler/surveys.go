package handler

import (
	"net/http"
	"src/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SurveyInput struct {
	Token string `json:"token" binding:"required"`
	Topic string `json:"topic"`
}

// func (handler *Handler) createSurvey(context *gin.Context) {
// 	var input SurveyInput
// 	if err := context.BindJSON(&input); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	userID, err := handler.service.Authorization.ParseToken(input.Token)
// 	if err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	survey, err := handler.service.Surveys.CreateSurvey(models.Survey{
// 		CreatorID: uint32(userID),
// 		Topic:     input.Topic,
// 	})
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	context.JSON(http.StatusOK, gin.H{"survey": survey})
// }

func (handler *Handler) getSurveys(context *gin.Context) {
	var input SurveyInput
	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := handler.service.Authorization.ParseToken(input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	surveys, err := handler.service.Surveys.GetSurveys(models.Survey{CreatorID: uint32(userID)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"surveys": surveys})
}

func (handler *Handler) getSurveyAnswers(context *gin.Context) {
	var input SurveyInput
	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := handler.service.Authorization.ParseToken(input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	surveyID, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid survey id param"})
		return
	}

	answers, err := handler.service.Surveys.GetSurveyAnswers(models.Survey{ID: uint32(surveyID)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"answers": answers})
}

func (handler *Handler) createSurveyWithQuestions(context *gin.Context) {
	var input struct {
		Token     string   `json:"token" binding:"required"`
		Topic     string   `json:"topic" binding:"required"`
		Questions []string `json:"questions" binding:"required"`
	}

	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := handler.service.Authorization.ParseToken(input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	survey, err := handler.service.Surveys.CreateSurvey(models.Survey{
		CreatorID: uint32(userID),
		Topic:     input.Topic,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var questions []models.Question
	for index, question := range input.Questions {
		questions = append(questions, models.Question{
			SurveyID:     uint32(survey.ID),
			Number:       int32(index) + 1,
			QuestionText: question,
			Type:         models.WithText,
		})
	}

	_, err = handler.service.CreateQuestions(questions)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"user_id": userID})
}

// params:
// + Token(json),
// + survey_id(context)
// returns:
// - user.ID,
// - survey.Topic
// - []questions
func (handler *Handler) getSurveyWithQuestions(context *gin.Context) {
	var input SurveyInput
	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := handler.service.Authorization.ParseToken(input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	surveyID, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid survey id param"})
		return
	}

	survey, err := handler.service.Surveys.GetSurvey(models.Survey{ID: uint32(surveyID)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	questions, err := handler.service.Questions.GetQuestions(models.Question{SurveyID: survey.ID})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"topic":     survey.Topic,
		"questions": questions,
	})
}

func (handler *Handler) deleteSurvey(context *gin.Context) {
	var input SurveyInput
	if err := context.BindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := handler.service.Authorization.ParseToken(input.Token)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	surveyID, err := strconv.ParseUint(context.Param("id"), 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid survey id param"})
		return
	}

	survey, err := handler.service.Surveys.DeleteSurvey(models.Survey{ID: uint32(surveyID)})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"survey": survey})
}
