package services

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/QBG-P2/Voting-System/config"
	"github.com/QBG-P2/Voting-System/pkg/cache"
	"github.com/QBG-P2/Voting-System/pkg/common"
	"github.com/QBG-P2/Voting-System/pkg/logging"
	"github.com/QBG-P2/Voting-System/pkg/service_errors"
)

type OtpService struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type otpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {
	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpService{logger: logger, cfg: cfg, redisClient: redis}
}
func (u *OtpService) SendOtp(mobileNumber string) error {
	otp := common.GenerateOtp()
	err := u.SetOtp(mobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}
func (s *OtpService) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", "otp", mobileNumber)
	val := &otpDto{
		Value: otp,
		Used:  false,
	}
	res, err := cache.Get[otpDto](s.redisClient, key)
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OptExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	err = cache.Set(s.redisClient, key, val, s.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}
func (s *OtpService) ValidateOtp(email string, otp string) error {
	key := fmt.Sprintf("%s:%s", "otp", email)
	res, err := cache.Get[otpDto](s.redisClient, key)
	if err != nil {
		return err
	} else if res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpNotValid}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(s.redisClient, key, res, s.cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
