package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

type nextFunc func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func ValidateJWTMiddleware(next nextFunc) nextFunc {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body:       "missing Auth Token",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		claims, err := parseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "user unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, err
		}
		expires := int64(claims["expires"].(float64))
		if time.Now().Unix() > expires {
			return events.APIGatewayProxyResponse{
				Body:       "token expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]
	if !ok {
		fmt.Println("no token in headers")
		return ""
	}
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		fmt.Printf("can not finde token: %v", authHeader)
		return ""
	}
	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		secret := "this_is_a_super_secure_secret"
		return []byte(secret), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}
	if !token.Valid {
		return jwt.MapClaims{}, fmt.Errorf("token is not valid")
	}
	claim, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("token is not valid - can not parse")
	}

	return claim, nil
}
