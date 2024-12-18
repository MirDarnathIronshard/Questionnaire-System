package services

import (
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/internal/models"
	"github.com/QBG-P2/Voting-System/pkg/logging"
	"github.com/QBG-P2/Voting-System/pkg/service_errors"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}
type tokenDto struct {
	UserId int
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}
func (s *TokenService) GenerateToken(token tokenDto) (*models.TokenDetail, error) {
	var err error
	td := &models.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Second).Unix()
	act := jwt.MapClaims{}
	act["UserId"] = token.UserId
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, act)
	td.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}
	rtc := jwt.MapClaims{}
	rtc["UserId"] = token.UserId
	return td, nil
}
func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}
func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}
	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}
