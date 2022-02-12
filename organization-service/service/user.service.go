package service

import (
	commonModels "commonpkg/models"
	token "commonpkg/token"
	"fmt"
	"net/http"
	"organization-service/persistance"
	"time"
)

var userServiceObj *UserService

type UserService struct {
	repo *persistance.UserPersistance
}

func InitUserService() (*UserService, *commonModels.ErrorDetail) {
	if userServiceObj == nil {
		repo, err := persistance.InitUserPersistance()
		if err != nil {
			return nil, err
		}
		userServiceObj = &UserService{
			repo: repo,
		}
	}
	return userServiceObj, nil

}

func (svc *UserService) Login(loginRequest commonModels.LoginRequest) *commonModels.LoginResponse {
	user, err := svc.repo.GetUserByUserName(loginRequest.Username)
	if err != nil {
		return &commonModels.LoginResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("error in validating the credential"),
				Errors: []commonModels.ErrorDetail{
					*err,
				},
			},
		}
	}

	if user.Password != loginRequest.Password {
		return &commonModels.LoginResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Invalid credential"),
				Errors: []commonModels.ErrorDetail{
					commonModels.ErrorDetail{
						ErrorCode:    commonModels.ErrorNoDataFound,
						ErrorMessage: fmt.Sprintf("Invalid credential"),
					},
				},
			},
		}
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	token, tokenError := token.GenrateToken(&commonModels.JwtClaims{
		Username: user.UserName,
	}, expirationTime)
	if tokenError != nil {
		return &commonModels.LoginResponse{
			CommonResponse: commonModels.CommonResponse{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("error in validating the credential"),
				Errors: []commonModels.ErrorDetail{
					commonModels.ErrorDetail{
						ErrorCode:    commonModels.ErrorServer,
						ErrorMessage: tokenError.Error(),
					},
				},
			},
		}
	}
	return &commonModels.LoginResponse{
		CommonResponse: commonModels.CommonResponse{
			StatusCode: http.StatusOK,
		},
		Token: token,
	}
}
