package auth

import (
	"context"
	"errors"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	userIDKey       contextKey = "userID"
	userRoleIDKey   contextKey = "userRoleID"
	userRoleNameKey contextKey = "userRoleName"
	userKey         contextKey = "user"
	claimsKey       contextKey = "claims"
)

func SetUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (*uint, error) {
	userID, ok := ctx.Value(userIDKey).(uint)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}
	return &userID, nil
}

func SetUserRoleID(ctx context.Context, roleID uint) context.Context {
	return context.WithValue(ctx, userRoleIDKey, roleID)
}

func GetUserRoleID(ctx context.Context) (*uint, error) {
	roleID, ok := ctx.Value(userRoleIDKey).(uint)
	if !ok {
		return nil, errors.New("user role ID not found in context")
	}
	return &roleID, nil
}

func SetUserRoleName(ctx context.Context, roleName string) context.Context {
	return context.WithValue(ctx, userRoleNameKey, roleName)
}

func GetUserRoleName(ctx context.Context) (*string, error) {
	roleName, ok := ctx.Value(userRoleNameKey).(string)
	if !ok {
		return nil, errors.New("user role name not found in context")
	}
	return &roleName, nil
}

func SetUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func SetClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetClaims(ctx context.Context) (jwt.MapClaims, error) {
	claims, ok := ctx.Value(claimsKey).(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims not found in context")
	}
	return claims, nil
}

func HasPermission(ctx context.Context, permission string) (bool, error) {
	user := GetUser(ctx)
	if user == nil {
		return false, errors.New("user not found in context")
	}

	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permission {
				return true, nil
			}
		}
	}

	return false, nil
}

func IsAuthenticated(ctx context.Context) bool {
	_, err := GetUserID(ctx)
	return err == nil
}

func HasRole(ctx context.Context, roleName string) (bool, error) {
	user := GetUser(ctx)
	if user == nil {
		return false, errors.New("user not found in context")
	}

	for _, role := range user.Roles {
		if role.Name == roleName {
			return true, nil
		}
	}

	return false, nil
}

func GetUserPermissions(ctx context.Context) ([]string, error) {
	user := GetUser(ctx)
	if user == nil {
		return nil, errors.New("user not found in context")
	}

	var permissions []string
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Name)
		}
	}

	return permissions, nil
}

func CreateAuthContext(ctx context.Context, user *models.User, claims jwt.MapClaims) context.Context {
	ctx = SetUser(ctx, user)
	ctx = SetUserID(ctx, user.ID)
	ctx = SetUserRoleID(ctx, user.Roles[0].ID)
	ctx = SetUserRoleName(ctx, user.Roles[0].Name)
	ctx = SetClaims(ctx, claims)
	return ctx
}

func ClearAuthContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, userKey, nil)
	ctx = context.WithValue(ctx, userIDKey, nil)
	ctx = context.WithValue(ctx, userRoleIDKey, nil)
	ctx = context.WithValue(ctx, userRoleNameKey, nil)
	ctx = context.WithValue(ctx, claimsKey, nil)
	return ctx
}

func IsAdmin(ctx context.Context) (bool, error) {
	return HasRole(ctx, "admin")
}

func IsSuperAdmin(ctx context.Context) (bool, error) {
	return HasRole(ctx, "super_admin")
}

func GetUserRoles(ctx context.Context) ([]string, error) {
	user := GetUser(ctx)
	if user == nil {
		return nil, errors.New("user not found in context")
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	return roles, nil
}
