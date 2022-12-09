package login_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/login"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/entity"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestCreate() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	login := domainmodel.Login{}

	newLogin := domainmodel.Login{}

	errorType := customerror.NoType

	sqlQuery := `INSERT INTO "logins" ("user_id","username","password","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingTheLogin",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainmodelfactory.NewLogin(args)

				args = map[string]interface{}{
					"userID":   login.UserID,
					"username": login.Username,
					"password": login.Password,
				}

				newLogin = domainmodelfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.NewV4()))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheLogin",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainmodelfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheLoginBecauseTheUserLoginIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				login = domainmodelfactory.NewLogin(args)

				mock.ExpectBegin()

				mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.Conflict.New("logins_user_id_key"))

				mock.ExpectRollback()

				errorType = customerror.Conflict
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			loginDatastoreRepository := logindatastorerepository.New(db)

			returnedLogin, err := loginDatastoreRepository.Create(login)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, newLogin.UserID, returnedLogin.UserID)
				assert.Equal(t, newLogin.Username, returnedLogin.Username)
				security := securitypkg.New()
				err := security.VerifyPasswords(returnedLogin.Password, newLogin.Password)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedLogin)
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
