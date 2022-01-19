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
	client.SortKey = common.GetClientSortKey(client.ClientId)

	_, err := service.clientServiceRepo.UpsertClient(client.ClientDto, true)
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
		contact.BranchId = client.BranchId
		contact.SortKey = common.GetClientContactSortKey(client.ClientId, contact.ContactId)
		_, err := service.clientServiceRepo.UpsertClientContact(contact)

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

func (service *ClientServiceService) Put(client commonModels.AddClientRequest) commonModels.AddClientResponse {
	client.SortKey = common.GetClientSortKey(client.ClientId)

	_, err := service.clientServiceRepo.UpsertClient(client.ClientDto, false)
	if err != nil {
		return commonModels.AddClientResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could no update client - %s", client.CompanyName),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	errors := make([]commonModels.ErrorDetail, 0)
	clientContacts := make([]commonModels.ContactPersonDto, len(client.ContactPersons))

	errDelete := deleteClientContact(client.BranchId, client.ClientId, client.ContactPersons)
	if errDelete != nil {
		return commonModels.AddClientResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could no delete contact person data client id - %s that were rermoved", client.ClientId),
				Errors: []commonModels.ErrorDetail{
					*errDelete,
				},
			},
		}
	}
	for i, contact := range client.ContactPersons {
		if contact.ContactId == "" {
			contactId, _ := uuid.NewV1()
			contact.ContactId = contactId.String()
		}
		contact.ClientId = client.ClientId
		contact.BranchId = client.BranchId
		contact.SortKey = common.GetClientContactSortKey(client.ClientId, contact.ContactId)
		_, err := service.clientServiceRepo.UpsertClientContact(contact)

		if err != nil {
			errors = append(errors, *err)
		}
		clientContacts[i] = contact
	}
	client.ContactPersons = clientContacts

	var status int = http.StatusOK
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
func deleteClientContact(branchId, clientId string, contactPersons []commonModels.ContactPersonDto) *commonModels.ErrorDetail {
	existingContact, err := ClientServiceObj.clientServiceRepo.GetPersonByClientId(commonModels.GetClientRequestDto{
		BranchId: branchId,
		ClientId: clientId,
	})

	if err != nil {
		return err
	}
	for _, exClientPerson := range existingContact {
		found := false
		for _, person := range contactPersons {
			if person.ContactId == exClientPerson.ContactId {
				found = true
			}
		}
		if !found {
			deleteErr := ClientServiceObj.clientServiceRepo.DeleteClientContact(branchId, clientId, exClientPerson.ContactId)
			if deleteErr != nil {
				return deleteErr
			}
		}
	}
	return nil
}
func (service *ClientServiceService) DeleteClient(request commonModels.GetClientRequestDto) commonModels.CommonResponse {
	existingContact, err := ClientServiceObj.clientServiceRepo.GetPersonByClientId(commonModels.GetClientRequestDto{
		BranchId: request.BranchId,
		ClientId: request.ClientId,
	})
	if err != nil {
		return commonModels.CommonResponse{
			ErrorMessage: fmt.Sprintf("error in getting contacts for client id %s", request.ClientId),
			StatusCode:   http.StatusBadRequest,
			Errors: []commonModels.ErrorDetail{
				*err,
			},
		}
	}
	for _, exClientPerson := range existingContact {
		deleteErr := ClientServiceObj.clientServiceRepo.DeleteClientContact(request.BranchId, request.ClientId, exClientPerson.ContactId)
		if deleteErr != nil {
			return commonModels.CommonResponse{
				ErrorMessage: fmt.Sprintf("error in deleting client contact id - %s for client id %s", exClientPerson.ContactId, request.ClientId),
				StatusCode:   http.StatusBadRequest,
				Errors: []commonModels.ErrorDetail{
					*deleteErr,
				},
			}
		}
	}
	clientDeleteErr := ClientServiceObj.clientServiceRepo.DeleteClient(request.BranchId, request.ClientId)
	if clientDeleteErr != nil {
		return commonModels.CommonResponse{
			ErrorMessage: fmt.Sprintf("error in deleting client id %s", request.ClientId),
			StatusCode:   http.StatusBadRequest,
			Errors: []commonModels.ErrorDetail{
				*clientDeleteErr,
			},
		}
	}
	return commonModels.CommonResponse{
		StatusCode: http.StatusOK,
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
			PageSize:         request.PageSize,
			Total:            count,
		},
		Data: client,
	}
}
