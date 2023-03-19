package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	authdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/auth"
	logindatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/login"
	userdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/user"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
	mockauth "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockauth"
	mocksecuritypkg "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockvalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestModifyPassword() {
	id := ""

	passwords := securitypkg.Passwords{}

	login := domainentity.Login{}

	updatedLogin := domainentity.Login{}

	errorType := customerror.NoType

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInModifyingThePassword",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				passwords = securitypkg.PasswordsFactory(nil)

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainentity.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: passwords.CurrentPassword,
				}

				updatedLogin = login
				updatedLogin.Password = passwords.NewPassword

				newLogin := updatedLogin

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{nil},
					{newLogin, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheIDIsNotValid",
			SetUp: func(t *testing.T) {
				id = ""

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{nil},
					{domainentity.Login{}, nil},
					{nil},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheEvaluatedPasswordsValuesAreNotValid",
			SetUp: func(t *testing.T) {
				passwords = securitypkg.Passwords{}

				returnArgs = ReturnArgs{
					{nil},
					{customerror.New("failed")},
					{domainentity.Login{}, nil},
					{nil},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{domainentity.Login{}, customerror.New("failed")},
					{nil},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheIDIsNotRegistered",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{domainentity.Login{}, nil},
					{nil},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.NotFound
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOfInvalidPasswordHappensWhenVerifyingThePasswords",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainentity.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{customerror.Unauthorized.New("the password is invalid")},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.Unauthorized
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnotherErrorHappensWhenVerifyingThePasswords",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainentity.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{customerror.New("failed")},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheNewPasswordISTheSameAsTheOneCurrentlyRegistered",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := currentPassword

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainentity.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{nil},
					{domainentity.Login{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenUpdatingTheLogin",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				loginID := uuid.NewV4()
				userID := uuid.NewV4()
				username := fake.Username()

				login = domainentity.Login{
					ID:       loginID,
					UserID:   userID,
					Username: username,
					Password: currentPassword,
				}

				updatedLogin = login
				updatedLogin.Password = newPassword

				returnArgs = ReturnArgs{
					{nil},
					{nil},
					{login, nil},
					{nil},
					{domainentity.Login{}, customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			validator := new(mockvalidator.Validator)
			validator.On("ValidateWithTags", id, "nonzero, uuid").Return(returnArgs[0]...)
			validator.On("Validate", passwords).Return(returnArgs[1]...)

			persistentLoginRepository := new(logindatastoremockrepository.Repository)
			persistentLoginRepository.On("GetByUserID", id).Return(returnArgs[2]...)

			security := new(mocksecuritypkg.Security)
			security.On("VerifyPasswords", login.Password, passwords.CurrentPassword).Return(returnArgs[3]...)

			persistentLoginRepository.On("Update", updatedLogin.ID.String(), updatedLogin).Return(returnArgs[4]...)

			persistentAuthRepository := new(authdatastoremockrepository.Repository)

			persistentUserRepository := new(userdatastoremockrepository.Repository)

			authN := new(mockauth.Auth)

			authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
				authN, security, validator, tokenExpTimeInSec)

			err := authService.ModifyPassword(id, passwords)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}
