package auth

import (
	"context"
	"net/http"
	"strings"

	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	"gorm.io/gorm"
)

// Auth is the function that wraps a http.Handler to evaluate the authentication of API based on a JWT token.
func Auth(db *gorm.DB, authN authpkg.IAuth, timeBeforeExpTimeInSec int) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			hdrAuth := r.Header.Get("Authorization")
			if len(hdrAuth) == 0 {
				errorMessage := "the auth header must be informed along with the token"
				responsehttputilpkg.RespondErrorWithJson(w, customerror.BadRequest.New(errorMessage))
				return
			}

			bearerToken := strings.Split(hdrAuth, " ")
			if len(bearerToken) != 2 {
				errorMessage := "the token must be associated with the auth header"
				responsehttputilpkg.RespondErrorWithJson(w, customerror.BadRequest.New(errorMessage))
				return
			}

			token, err := authN.DecodeToken(bearerToken[1])
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, customerror.Unauthorized.New(err.Error()))
				return
			}

			auth, err := authN.FetchAuthFromToken(token)
			if err != nil {
				responsehttputilpkg.RespondErrorWithJson(w, err)
				return
			}

			// Before proceeding is necessary to check if the user who is performing operations is logged
			// based on the authentication details inserted within in the token.
			authAux := domainmodel.Auth{}

			result := db.Find(&authAux, "id=?", auth.ID)
			if result.Error != nil {
				responsehttputilpkg.RespondErrorWithJson(w, result.Error)
				return
			}

			if authAux.IsEmpty() {
				errorMessage := "you are not logged in, then perform a login to get a token before proceeding"
				responsehttputilpkg.RespondErrorWithJson(w, customerror.BadRequest.New(errorMessage))
				return
			}

			if auth.UserID.String() != authAux.UserID.String() {
				errorMessage := "the token's auth_id and user_id are not associated"
				responsehttputilpkg.RespondErrorWithJson(w, customerror.BadRequest.New(errorMessage))
				return
			}

			// It is necessary to set auth details that can be used for performing authenticated operations.
			ctx := r.Context()
			var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
			ctx = context.WithValue(ctx, authDetailsKey, auth)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}
