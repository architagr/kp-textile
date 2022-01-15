package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"client-service/common"
	"client-service/persistance"

	uuid "github.com/iris-contrib/go.uuid"
)

var ClientServiceObj *ClientServiceService

type ClientServiceService struct {
	clientServiceRepo *persistance.ClientServicePersistance
}

func InitClientServiceService() (*ClientServiceService, *commonModels.ErrorDetail) {
	if ClientServiceObj == nil {
		repo, err := persistance.InitClientServicePersistance()
		if err != nil {
			return nil, err
		}
		ClientServiceObj = &ClientServiceService{
			clientServiceRepo: repo,
		}
	}
	return ClientServiceObj, nil
}

func (service *ClientServiceService) Add(client commonModels.AddClientRequest) commonModels.AddClientResponse {
	clientid, _ := uuid.NewV1()
	client.ClientId = clientid.String()
	client.SortKey = fmt.Sprintf("%s|%s", common.ClientSortKey, client.ClientId)

	_, err := service.clientServiceRepo.AddClient(client.ClientDto)
	if err != nil {
		return commonModels.AddClientResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add client - %s", client.CompanyName),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	errors := make([]commonModels.ErrorDetail, 0)
	clientContacts := make([]commonModels.ContactPersonDto, len(client.ContactPersons))
	for i, contact := range client.ContactPersons {
		contactId, _ := uuid.NewV1()
		contact.ContactId = contactId.String()
		contact.ClientId = client.ClientId
		contact.SortKey = fmt.Sprintf("%s|%s|%s", common.ContactSortKey, client.ClientId, contact.ContactId)
		_, err := service.clientServiceRepo.AddClientContact(contact)

		if err != nil {
			errors = append(errors, *err)
		}
		clientContacts[i] = contact
	}
	client.ContactPersons = clientContacts

	var status int = http.StatusCreated
	if len(errors) > 0 {
		status = http.StatusPartialContent
	}
	return commonModels.AddClientResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: status,
			Errors:     errors,
		},
		Data: client,
	}
}

func (service *ClientServiceService) GetClient(request commonModels.GetClientRequestDto) commonModels.AddClientResponse {
	client, err := service.clientServiceRepo.GetClient(request)
	if err != nil {
		return commonModels.AddClientResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get client with client id  - %s", request.ClientId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	clientPersons, err := service.clientServiceRepo.GetPersonByClientId(request)
	if err != nil {
		return commonModels.AddClientResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get client with client id  - %s", request.ClientId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	return commonModels.AddClientResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: commonModels.AddClientRequest{
			ClientDto:      client,
			ContactPersons: clientPersons,
		},
	}
}

func (service *ClientServiceService) GetAll(request commonModels.ClientListRequest) commonModels.ClientListResponse {

	client, lastEvaluationKey, err := service.clientServiceRepo.GetClientByFilter(request)
	if err != nil {
		return commonModels.ClientListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get clients",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	request.LastEvalutionKey = nil
	//request.PageSize = 100
	count, _ := service.clientServiceRepo.GetClientTotalByFilter(request)

	return commonModels.ClientListResponse{
		CommonListResponse: commonModels.CommonListResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			LastEvalutionKey: lastEvaluationKey,
			PageSize:         int64(len(client)),
			Total:            count,
		},
		Data: client,
	}
}
