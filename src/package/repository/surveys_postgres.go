package repository

import (
	"src/models"

	"gorm.io/gorm"
)

type SurveyPostgres struct {
	DB *gorm.DB
}

func NewSurveyPostgres(DB *gorm.DB) *SurveyPostgres {
	return &SurveyPostgres{DB: DB}
}

func (postgres *SurveyPostgres) CreateSurvey(survey models.Survey) (*models.Survey, error) {
	err := postgres.DB.Create(&survey).Error
	if err != nil {
		return nil, err
	}
	return postgres.GetSurvey(survey)
}

func (postgres *SurveyPostgres) GetSurvey(survey models.Survey) (*models.Survey, error) {
	err := postgres.DB.Where(&survey).First(&survey).Error
	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (postgres *SurveyPostgres) GetSurveys(survey models.Survey) (*[]models.Survey, error) {
	var surveys []models.Survey
	err := postgres.DB.Where("creator_id = ?", survey.CreatorID).Find(&surveys).Error
	if err != nil {
		return nil, err
	}
	return &surveys, nil
}

func (postgres *SurveyPostgres) DeleteSurvey(_survey models.Survey) (*models.Survey, error) {
	survey, err := postgres.GetSurvey(_survey)
	if err != nil {
		return nil, err
	}

	err = postgres.DB.Where("id = ?", survey.ID).Delete(&models.Survey{}).Error
	if err != nil {
		return nil, err
	}

	return survey, nil
}

func (postgres *SurveyPostgres) GetSurveyAnswers(survey models.Survey) (*[]models.Answer, error) {
	var answers *[]models.Answer
	err := postgres.DB.Where("survey_id = ?", survey.ID).Find(&answers).Error
	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (postgres *SurveyPostgres) GetAllSurveys() (*[]models.Survey, error) {
	var surveys *[]models.Survey
	err := postgres.DB.Find(&surveys).Error
	if err != nil {
		return nil, err
	}
	return surveys, nil
}

func (postgres *SurveyPostgres) GetSurveysNumber(user models.User) (int, error) {
	var count int64
	err := postgres.DB.Where("creator_id = ?", user.ID).Find(&models.Survey{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), err
}
