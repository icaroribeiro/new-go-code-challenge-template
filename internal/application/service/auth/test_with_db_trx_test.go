package auth_test

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/auth"
	authdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/mockrepository/auth"
	logindatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/mockrepository/login"
	userdatastoremockrepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/infrastructure/storage/datastore/mockrepository/user"
	mockauth "github.com/icaroribeiro/new-go-code-challenge-template/tests/mocks/pkg/mockauth"
	mocksecurity "github.com/icaroribeiro/new-go-code-challenge-template/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/new-go-code-challenge-template/tests/mocks/pkg/mockvalidator"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestWithDBTrx() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	authDatastoreRepositoryWithDBTrx := &authdatastoremockrepository.Repository{}
	userDatastoreRepositoryWithDBTrx := &userdatastoremockrepository.Repository{}
	loginDatastoreRepositoryWithDBTrx := &logindatastoremockrepository.Repository{}

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSettingRepositoriesWithDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = db

				authDatastoreRepositoryWithDBTrx = &authdatastoremockrepository.Repository{}
				userDatastoreRepositoryWithDBTrx = &userdatastoremockrepository.Repository{}
				loginDatastoreRepositoryWithDBTrx = &logindatastoremockrepository.Repository{}

				returnArgs = ReturnArgs{
					{authDatastoreRepositoryWithDBTrx},
					{userDatastoreRepositoryWithDBTrx},
					{loginDatastoreRepositoryWithDBTrx},
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := new(authdatastoremockrepository.Repository)
			authDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[0]...)

			userDatastoreRepository := new(userdatastoremockrepository.Repository)
			userDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[1]...)

			loginDatastoreRepository := new(logindatastoremockrepository.Repository)
			loginDatastoreRepository.On("WithDBTrx", dbTrx).Return(returnArgs[2]...)

			authN := new(mockauth.Auth)
			security := new(mocksecurity.Security)
			validator := new(mockvalidator.Validator)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedAuthService := authService.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedAuthService, "Service interface is empty.")
				assert.Equal(t, authService, returnedAuthService, "Service interfaces are not the same.")
			}
		})
	}
}
