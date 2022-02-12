package service

import (
	commonModels "commonpkg/models"
	"organization-service/persistance"
)

var godownServiceObj *GodownService

type GodownService struct {
	repo *persistance.GodownPersistance
}

func InitGodownService() (*GodownService, *commonModels.ErrorDetail) {
	if godownServiceObj == nil {
		repo, err := persistance.InitGodownPersistance()
		if err != nil {
			return nil, err
		}
		godownServiceObj = &GodownService{
			repo: repo,
		}
	}
	return godownServiceObj, nil
}
