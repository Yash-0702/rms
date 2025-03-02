package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	dbhelper "rms/database/dbHelper"
	"rms/models"
	"rms/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKeys string

const (
	userContext ContextKeys = "userContext"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		// Check if the user is authenticated
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ResponseError(res, http.StatusUnauthorized, nil, "missing or invalid Authorization header")
			return
		}

		//Bearer remove
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if the token signature is valid
		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// Check if the token is valid and delete session if token is vaild but expired
		if parseErr != nil || !token.Valid {

			if errors.Is(parseErr, jwt.ErrTokenExpired) {
				utils.ResponseError(res, http.StatusUnauthorized, parseErr, "token expired")

				// delete session if token is expired
				sessionId := token.Claims.(jwt.MapClaims)["sessionID"].(string)
				sessionErr := dbhelper.DeleteSession(sessionId)
				if sessionErr != nil {
					utils.ResponseError(res, http.StatusInternalServerError, sessionErr, "failed to delete session")
				}
				return
			}

			utils.ResponseError(res, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		// Get the claims from the token
		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ResponseError(res, http.StatusUnauthorized, nil, "error getting claims")
			return
		}

		// check if session exist or not (if log out but token is still valid not expired)
		exist, existErr := dbhelper.IsSessionExist(claimValues["sessionID"].(string))
		if existErr != nil {
			utils.ResponseError(res, http.StatusInternalServerError, existErr, "failed to check if session exists")
			return
		}
		if !exist {
			utils.ResponseError(res, http.StatusUnauthorized, nil, "session is  expired")
			return
		}

		// Get the email and user_id from the claims
		email := claimValues["email"].(string)
		userId := claimValues["user_id"].(string)
		sessionId := claimValues["sessionID"].(string)
		role := claimValues["role"].(string)

		// fmt.Println(userId)

		user := &models.UserCtx{
			UserId:    userId,
			Email:     email,
			SessionId: sessionId,
			Role:      role,
		}

		ctx := context.WithValue(req.Context(), userContext, user)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}

func UserContext(req *http.Request) *models.UserCtx {
	if user, ok := req.Context().Value(userContext).(*models.UserCtx); ok {
		return user
	}
	return nil
}
