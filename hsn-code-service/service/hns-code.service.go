package service

import (
	"hsn-code-service/model"
	"hsn-code-service/persistance"
)

// TODO: return standard response model and error
var hnsCodeServiceObj *HnsCodeService

type HnsCodeService struct {
	hnsCodeRepo *persistance.HnsCodePersistance
}

func InitHnsCodeService(hnsCodeRepoObj *persistance.HnsCodePersistance) *HnsCodeService {
	if hnsCodeServiceObj == nil {
		hnsCodeServiceObj = &HnsCodeService{
			hnsCodeRepo: hnsCodeRepoObj,
		}
	}
	return hnsCodeServiceObj
}

func (service *HnsCodeService) GetAll() []model.HnsCodeDto {
	return service.hnsCodeRepo.GetAll()
}

func (service *HnsCodeService) Get(id int) model.HnsCodeDto {
	return service.hnsCodeRepo.Get(id)
}

func (service *HnsCodeService) Add(code string) model.HnsCodeDto {
	return service.hnsCodeRepo.Add(code)
}

func (service *HnsCodeService) AddMultiple(codes []string) []model.HnsCodeDto {
	return service.hnsCodeRepo.AddMultiple(codes)
}
