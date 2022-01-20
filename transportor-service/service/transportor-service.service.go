package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"transportor-service/common"
	"transportor-service/persistance"

	uuid "github.com/iris-contrib/go.uuid"
)

var transporterServiceObj *TransporterService

type TransporterService struct {
	transporterServiceRepo *persistance.TransporterPersistance
}

func InitTransporterService() (*TransporterService, *commonModels.ErrorDetail) {
	if transporterServiceObj == nil {
		repo, err := persistance.InitTransporterPersistance()
		if err != nil {
			return nil, err
		}
		transporterServiceObj = &TransporterService{
			transporterServiceRepo: repo,
		}
	}
	return transporterServiceObj, nil
}

func (service *TransporterService) Add(transporter commonModels.AddTransporterRequest) commonModels.AddTransporterResponse {
	transporterid, _ := uuid.NewV1()
	transporter.TransporterId = transporterid.String()
	transporter.SortKey = common.GetTransporterSortKey(transporter.TransporterId)

	_, err := service.transporterServiceRepo.UpsertTransporter(transporter.TransporterDto, true)
	if err != nil {
		return commonModels.AddTransporterResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add transporter - %s", transporter.CompanyName),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	errors := make([]commonModels.ErrorDetail, 0)
	transporterContacts := make([]commonModels.TransporterContactPersonDto, len(transporter.ContactPersons))
	for i, contact := range transporter.ContactPersons {
		contactId, _ := uuid.NewV1()
		contact.ContactId = contactId.String()
		contact.TransporterId = transporter.TransporterId
		contact.BranchId = transporter.BranchId
		contact.SortKey = common.GetTransporterContactSortKey(transporter.TransporterId, contact.ContactId)
		_, err := service.transporterServiceRepo.UpsertTransporterContact(contact)

		if err != nil {
			errors = append(errors, *err)
		}
		transporterContacts[i] = contact
	}
	transporter.ContactPersons = transporterContacts

	var status int = http.StatusCreated
	if len(errors) > 0 {
		status = http.StatusPartialContent
	}
	return commonModels.AddTransporterResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: status,
			Errors:     errors,
		},
		Data: transporter,
	}
}

func (service *TransporterService) Put(transporter commonModels.AddTransporterRequest) commonModels.AddTransporterResponse {
	transporter.SortKey = common.GetTransporterSortKey(transporter.TransporterId)

	_, err := service.transporterServiceRepo.UpsertTransporter(transporter.TransporterDto, false)
	if err != nil {
		return commonModels.AddTransporterResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could no update transporter - %s", transporter.CompanyName),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	errors := make([]commonModels.ErrorDetail, 0)
	transporterContacts := make([]commonModels.TransporterContactPersonDto, len(transporter.ContactPersons))

	errDelete := deleteTransporterContact(transporter.BranchId, transporter.TransporterId, transporter.ContactPersons)
	if errDelete != nil {
		return commonModels.AddTransporterResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could no delete contact person data transporter id - %s that were rermoved", transporter.TransporterId),
				Errors: []commonModels.ErrorDetail{
					*errDelete,
				},
			},
		}
	}
	for i, contact := range transporter.ContactPersons {
		if contact.ContactId == "" {
			contactId, _ := uuid.NewV1()
			contact.ContactId = contactId.String()
		}
		contact.TransporterId = transporter.TransporterId
		contact.BranchId = transporter.BranchId
		contact.SortKey = common.GetTransporterContactSortKey(transporter.TransporterId, contact.ContactId)
		_, err := service.transporterServiceRepo.UpsertTransporterContact(contact)

		if err != nil {
			errors = append(errors, *err)
		}
		transporterContacts[i] = contact
	}
	transporter.ContactPersons = transporterContacts

	var status int = http.StatusOK
	if len(errors) > 0 {
		status = http.StatusPartialContent
	}
	return commonModels.AddTransporterResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: status,
			Errors:     errors,
		},
		Data: transporter,
	}
}
func deleteTransporterContact(branchId, transporterId string, contactPersons []commonModels.TransporterContactPersonDto) *commonModels.ErrorDetail {
	existingContact, err := transporterServiceObj.transporterServiceRepo.GetPersonByTransporterId(commonModels.GetTransporterRequestDto{
		BranchId:      branchId,
		TransporterId: transporterId,
	})

	if err != nil {
		return err
	}
	for _, exTransporterPerson := range existingContact {
		found := false
		for _, person := range contactPersons {
			if person.ContactId == exTransporterPerson.ContactId {
				found = true
			}
		}
		if !found {
			deleteErr := transporterServiceObj.transporterServiceRepo.DeleteTransporterContact(branchId, transporterId, exTransporterPerson.ContactId)
			if deleteErr != nil {
				return deleteErr
			}
		}
	}
	return nil
}
func (service *TransporterService) DeleteTransporter(request commonModels.GetTransporterRequestDto) commonModels.CommonResponse {
	existingContact, err := transporterServiceObj.transporterServiceRepo.GetPersonByTransporterId(commonModels.GetTransporterRequestDto{
		BranchId:      request.BranchId,
		TransporterId: request.TransporterId,
	})
	if err != nil {
		return commonModels.CommonResponse{
			ErrorMessage: fmt.Sprintf("error in getting contacts for transporter id %s", request.TransporterId),
			StatusCode:   http.StatusBadRequest,
			Errors: []commonModels.ErrorDetail{
				*err,
			},
		}
	}
	for _, exTransporterPerson := range existingContact {
		deleteErr := transporterServiceObj.transporterServiceRepo.DeleteTransporterContact(request.BranchId, request.TransporterId, exTransporterPerson.ContactId)
		if deleteErr != nil {
			return commonModels.CommonResponse{
				ErrorMessage: fmt.Sprintf("error in deleting Transporter contact id - %s for Transporter id %s", exTransporterPerson.ContactId, request.TransporterId),
				StatusCode:   http.StatusBadRequest,
				Errors: []commonModels.ErrorDetail{
					*deleteErr,
				},
			}
		}
	}
	transporterDeleteErr := transporterServiceObj.transporterServiceRepo.DeleteTransporter(request.BranchId, request.TransporterId)
	if transporterDeleteErr != nil {
		return commonModels.CommonResponse{
			ErrorMessage: fmt.Sprintf("error in deleting transporter id %s", request.TransporterId),
			StatusCode:   http.StatusBadRequest,
			Errors: []commonModels.ErrorDetail{
				*transporterDeleteErr,
			},
		}
	}
	return commonModels.CommonResponse{
		StatusCode: http.StatusOK,
	}
}
func (service *TransporterService) GetTransporter(request commonModels.GetTransporterRequestDto) commonModels.AddTransporterResponse {
	transporter, err := service.transporterServiceRepo.GetTransporter(request)
	if err != nil {
		return commonModels.AddTransporterResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get transporter with transporter id  - %s", request.TransporterId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	transporterPersons, err := service.transporterServiceRepo.GetPersonByTransporterId(request)
	if err != nil {
		return commonModels.AddTransporterResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get transporter with transporter id  - %s", request.TransporterId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	return commonModels.AddTransporterResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: commonModels.AddTransporterRequest{
			TransporterDto: transporter,
			ContactPersons: transporterPersons,
		},
	}
}

func (service *TransporterService) GetAll(request commonModels.TransporterListRequest) commonModels.TransporterListResponse {

	transporter, lastEvaluationKey, err := service.transporterServiceRepo.GetTransporterByFilter(request)
	if err != nil {
		return commonModels.TransporterListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get transporter",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	request.LastEvalutionKey = nil
	//request.PageSize = 100
	count, _ := service.transporterServiceRepo.GetTransporterTotalByFilter(request)

	return commonModels.TransporterListResponse{
		CommonListResponse: commonModels.CommonListResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			LastEvalutionKey: lastEvaluationKey,
			PageSize:         request.PageSize,
			Total:            count,
		},
		Data: transporter,
	}
}
