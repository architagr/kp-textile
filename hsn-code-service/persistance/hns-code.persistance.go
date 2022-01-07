package persistance

import (
	"hsn-code-service/model"
)

var hnsCodes = []model.HnsCodeDto{
	{
		Id:      1,
		HnsCode: "1-123",
	},
	{
		Id:      2,
		HnsCode: "2-123",
	},
}
var hnsCodepersistanceObj *HnsCodePersistance

type HnsCodePersistance struct {
}

func InitHnsCodePersistance() *HnsCodePersistance {
	// TODO: get connection to DB
	if hnsCodepersistanceObj == nil {
		hnsCodepersistanceObj = &HnsCodePersistance{}
	}

	return hnsCodepersistanceObj
}

func (repo *HnsCodePersistance) GetAll() []model.HnsCodeDto {
	// TODO: get data from DB
	return hnsCodes
}

func (repo *HnsCodePersistance) Get(id int) model.HnsCodeDto {
	var hnsCode *model.HnsCodeDto
	// TODO: get data from DB
	for _, value := range hnsCodes {
		if value.Id == id {
			hnsCode = &value
			break
		}
	}
	// TODO: throw custom error if hnsCode == nil

	return *hnsCode
}

func (repo *HnsCodePersistance) Add(code string) model.HnsCodeDto {
	length := len(hnsCodes)

	id := hnsCodes[length-1].Id + 1

	newHnsCode := model.HnsCodeDto{Id: id, HnsCode: code}
	hnsCodes = append(hnsCodes, newHnsCode)
	return newHnsCode
}

func (repo *HnsCodePersistance) AddMultiple(codes []string) []model.HnsCodeDto {
	var newHnsCodes []model.HnsCodeDto
	for _, val := range codes {
		newHnsCode := repo.Add(val)

		newHnsCodes = append(newHnsCodes, newHnsCode)
	}

	// TODO: return error of the codes which have not been added
	return newHnsCodes
}
