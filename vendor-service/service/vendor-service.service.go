package service

import (
	commonModels "commonpkg/models"
	"fmt"
	"net/http"

	"vendor-service/common"
	"vendor-service/persistance"

	uuid "github.com/iris-contrib/go.uuid"
)

var VendorServiceObj *VendorService

type VendorService struct {
	vendorServiceRepo *persistance.VendorPersistance
}

func InitVendorService() (*VendorService, *commonModels.ErrorDetail) {
	if VendorServiceObj == nil {
		repo, err := persistance.InitVendorPersistance()
		if err != nil {
			return nil, err
		}
		VendorServiceObj = &VendorService{
			vendorServiceRepo: repo,
		}
	}
	return VendorServiceObj, nil
}

func (service *VendorService) Add(vendor commonModels.AddVendorRequest) commonModels.AddVendorResponse {
	vendorid, _ := uuid.NewV1()
	vendor.VendorId = vendorid.String()
	vendor.SortKey = common.GetVendorSortKey()

	_, err := service.vendorServiceRepo.UpsertVendor(vendor.VendorDto, true)
	if err != nil {
		return commonModels.AddVendorResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not add vendor - %s as, %s", vendor.CompanyName, err.ErrorMessage),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	errors := make([]commonModels.ErrorDetail, 0)
	vendorContacts := make([]commonModels.VendorContactPersonDto, len(vendor.ContactPersons))
	for i, contact := range vendor.ContactPersons {
		contactId, _ := uuid.NewV1()
		contact.ContactId = contactId.String()
		contact.VendorId = vendor.VendorId
		contact.SortKey = common.GetVendorContactSortKey(contact.ContactId)
		_, err := service.vendorServiceRepo.UpsertVendorContact(contact)

		if err != nil {
			errors = append(errors, *err)
		}
		vendorContacts[i] = contact
	}
	vendor.ContactPersons = vendorContacts

	var status int = http.StatusCreated
	if len(errors) > 0 {
		status = http.StatusPartialContent
	}
	return commonModels.AddVendorResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: status,
			Errors:     errors,
		},
		Data: vendor,
	}
}

func (service *VendorService) Put(vendor commonModels.AddVendorRequest) commonModels.AddVendorResponse {
	vendor.SortKey = common.GetVendorSortKey()

	_, err := service.vendorServiceRepo.UpsertVendor(vendor.VendorDto, false)
	if err != nil {
		return commonModels.AddVendorResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not update vendor - %s as, %s", vendor.CompanyName, err.ErrorMessage),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	errors := make([]commonModels.ErrorDetail, 0)
	vendorContacts := make([]commonModels.VendorContactPersonDto, len(vendor.ContactPersons))

	errDelete := deleteVendorContact(vendor.VendorId, vendor.ContactPersons)
	if errDelete != nil {
		return commonModels.AddVendorResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not delete contact person data vendor id - %s that were rermoved", vendor.VendorId),
				Errors: []commonModels.ErrorDetail{
					*errDelete,
				},
			},
		}
	}
	for i, contact := range vendor.ContactPersons {
		if contact.ContactId == "" {
			contactId, _ := uuid.NewV1()
			contact.ContactId = contactId.String()
		}
		contact.VendorId = vendor.VendorId
		contact.SortKey = common.GetVendorContactSortKey(contact.ContactId)
		_, err := service.vendorServiceRepo.UpsertVendorContact(contact)

		if err != nil {
			errors = append(errors, *err)
		}
		vendorContacts[i] = contact
	}
	vendor.ContactPersons = vendorContacts

	var status int = http.StatusOK
	if len(errors) > 0 {
		status = http.StatusPartialContent
	}
	return commonModels.AddVendorResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: status,
			Errors:     errors,
		},
		Data: vendor,
	}
}
func deleteVendorContact(vendorId string, contactPersons []commonModels.VendorContactPersonDto) *commonModels.ErrorDetail {
	existingContact, err := VendorServiceObj.vendorServiceRepo.GetPersonByVendorId(commonModels.GetVendorRequestDto{
		VendorId: vendorId,
	})

	if err != nil {
		return err
	}
	for _, exVendorPerson := range existingContact {
		found := false
		for _, person := range contactPersons {
			if person.ContactId == exVendorPerson.ContactId {
				found = true
			}
		}
		if !found {
			deleteErr := VendorServiceObj.vendorServiceRepo.DeleteVendorContact(vendorId, exVendorPerson.ContactId)
			if deleteErr != nil {
				return deleteErr
			}
		}
	}
	return nil
}
func (service *VendorService) DeleteVendor(request commonModels.GetVendorRequestDto) commonModels.CommonResponse {
	existingContact, err := VendorServiceObj.vendorServiceRepo.GetPersonByVendorId(commonModels.GetVendorRequestDto{
		VendorId: request.VendorId,
	})
	if err != nil {
		return commonModels.CommonResponse{
			ErrorMessage: fmt.Sprintf("error in getting contacts for vendor id %s", request.VendorId),
			StatusCode:   http.StatusBadRequest,
			Errors: []commonModels.ErrorDetail{
				*err,
			},
		}
	}
	for _, exVendorPerson := range existingContact {
		deleteErr := VendorServiceObj.vendorServiceRepo.DeleteVendorContact(request.VendorId, exVendorPerson.ContactId)
		if deleteErr != nil {
			return commonModels.CommonResponse{
				ErrorMessage: fmt.Sprintf("error in deleting vendor contact id - %s for vendor id %s", exVendorPerson.ContactId, request.VendorId),
				StatusCode:   http.StatusBadRequest,
				Errors: []commonModels.ErrorDetail{
					*deleteErr,
				},
			}
		}
	}
	vendorDeleteErr := VendorServiceObj.vendorServiceRepo.DeleteVendor(request.VendorId)
	if vendorDeleteErr != nil {
		return commonModels.CommonResponse{
			ErrorMessage: fmt.Sprintf("error in deleting vendor id %s", request.VendorId),
			StatusCode:   http.StatusBadRequest,
			Errors: []commonModels.ErrorDetail{
				*vendorDeleteErr,
			},
		}
	}
	return commonModels.CommonResponse{
		StatusCode: http.StatusOK,
	}
}
func (service *VendorService) GetVendor(request commonModels.GetVendorRequestDto) commonModels.AddVendorResponse {
	vendor, err := service.vendorServiceRepo.GetVendor(request)
	if err != nil {
		return commonModels.AddVendorResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get vendor with vendor id  - %s", request.VendorId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	vendorPersons, err := service.vendorServiceRepo.GetPersonByVendorId(request)
	if err != nil {
		return commonModels.AddVendorResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("could not get vendor with vendor id  - %s", request.VendorId),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}
	return commonModels.AddVendorResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Data: commonModels.AddVendorRequest{
			VendorDto:      vendor,
			ContactPersons: vendorPersons,
		},
	}
}

func (service *VendorService) GetAll(request commonModels.VendorListRequest) commonModels.VendorListResponse {

	vendor, lastEvaluationKey, err := service.vendorServiceRepo.GetVendorByFilter(request)
	if err != nil {
		return commonModels.VendorListResponse{
			CommonListResponse: commonModels.CommonListResponse{
				CommonResponse: commonModels.CommonResponse{
					StatusCode:   http.StatusBadRequest,
					ErrorMessage: "could not get vendors",
					Errors: []commonModels.ErrorDetail{
						*err,
					},
				},
			},
		}
	}
	request.LastEvalutionKey = nil
	//request.PageSize = 100
	count, _ := service.vendorServiceRepo.GetVendorTotalByFilter(request)
	return commonModels.VendorListResponse{
		CommonListResponse: commonModels.CommonListResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode: http.StatusOK,
			},
			LastEvalutionKey: lastEvaluationKey,
			PageSize:         request.PageSize,
			Total:            count,
		},
		Data: vendor,
	}
}
