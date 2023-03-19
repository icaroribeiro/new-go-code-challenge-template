package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	persistententity "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/perentity"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/user"
	authhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/auth"
	authpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/auth"
	requesthttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignOut() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	persistentAuth := persistententity.Auth{}
	persistentLogin := persistententity.Login{}
	persistentUser := persistententity.User{}

	auth := domainentity.Auth{}

	authDetailsCtxValue := domainentity.Auth{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningOut",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				persistentUser = persistententity.User{
					Username: username,
				}

				result := dbTrx.Create(&persistentUser)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				persistentLogin = persistententity.Login{
					UserID:   persistentUser.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&persistentLogin)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				persistentAuth = persistententity.Auth{
					UserID: persistentUser.ID,
				}

				result = dbTrx.Create(&persistentAuth)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				auth = persistentAuth.ToDomain()

				authDetailsCtxValue = auth
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				authDetailsCtxValue = domainentity.Auth{}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				authN = authpkg.New(ts.RSAKeys)

				auth = domainentity.Auth{}

				authDetailsCtxValue = auth
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentAuthRepository := authdatastorerepository.New(dbTrx)
			persistentLoginRepository := logindatastorerepository.New(dbTrx)
			persistentUserRepository := userdatastorerepository.New(dbTrx)

			authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "SignOut",
				Method:      http.MethodPost,
				Path:        "/sign_out",
				HandlerFunc: authHandler.SignOut,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			ctx := req.Context()
			ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
			req = req.WithContext(ctx)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedMessage := responsehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, returnedMessage.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}
