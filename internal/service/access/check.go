package access

import (
	"auth/internal/utils"
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

func (s *AccessService) Check(ctx context.Context, handler_path string) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("medata is ot provided")
	}

	authHeader, ok := md["authorization"]

	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authprefix) {
		return errors.New("return invalid authorization header")

	}

	accessToken := strings.TrimPrefix(authHeader[0], authprefix)

	claims, err := utils.VerifyToken(accessToken, s.JWT.AccessTokenSecret)
	if err != nil {
		return errors.New("Access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[handler_path]
	if !ok {
		return errors.New("You don't have access 1")
	}

	if role == claims.Role {
		return nil
	}

	return errors.New("You don't have access")
}

func (s *AccessService) accessibleRoles(ctx context.Context) (map[string]string, error) {

	roles := map[string]string{
		"/user_v1.UserV1/Get":   "ADMIN",
		"/user_v1.UserV1/Store": "ADMIN",
	}

	return roles, nil

}
