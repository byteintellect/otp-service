package svc

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/otp_svc/config"
	"github.com/byteintellect/otp_svc/pkg/domain"
	"github.com/byteintellect/otp_svc/pkg/repo"
	commonsv1 "github.com/byteintellect/protos_go/commons/v1"
	otpsv1 "github.com/byteintellect/protos_go/otps/v1"
	"go.uber.org/zap"
	"math/big"
)

const (
	seed = "0123456789"
)

type OtpSvc struct {
	repo.OtpRepo
	otpsv1.UnimplementedOtpServiceServer
	cfg    *config.OtpServiceConfig
	logger *zap.Logger
}

func NewOtpSvc(db db.BaseRepository, cfg *config.OtpServiceConfig, logger *zap.Logger) otpsv1.OtpServiceServer {
	return &OtpSvc{
		OtpRepo:                       repo.NewOtpRepo(db),
		UnimplementedOtpServiceServer: otpsv1.UnimplementedOtpServiceServer{},
		cfg:                           cfg,
		logger:                        logger,
	}
}

func (o OtpSvc) GenerateOTPCode() (string, error) {
	byteSlice := make([]byte, o.cfg.OtpLength)
	for i := 0; i < o.cfg.OtpLength; i++ {
		max := big.NewInt(int64(len(seed)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		byteSlice[i] = seed[num.Int64()]
	}
	return string(byteSlice), nil
}

func (o OtpSvc) GetOtp(ctx context.Context, request *commonsv1.GetOtpForPhoneRequest) (*commonsv1.GetOtpForPhoneResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		otp, err := o.GenerateOTPCode()
		if err != nil {
			return nil, err
		}
		otpDomain := domain.NewOtp(request.PhoneNumber, otp)
		err, createdOtpDomain := o.OtpRepo.Create(ctx, otpDomain)
		if err != nil {
			return nil, err
		}
		return &commonsv1.GetOtpForPhoneResponse{
			Response: createdOtpDomain.(*domain.Otp).ToDto().(*commonsv1.AuthOtpDto),
		}, nil
	}
}

func (o OtpSvc) ValidateOtp(ctx context.Context, request *commonsv1.AuthValidateOtpRequest) (*commonsv1.AuthValidateOtpResponse, error) {
	select {
	case <-ctx.Done():
		return nil, errors.New("timed out")
	default:
		if otp, err := o.OtpRepo.GetOtp(ctx, request.OtpId, request.PhoneNumber, request.Otp); err != nil || otp == nil {
			return nil, err
		} else {
			o.logger.Info("otp validation", zap.String("otp_id", otp.ExternalId))
			err := o.OtpRepo.InvalidateOtp(ctx, otp)
			if err != nil {
				o.logger.Error("error invalidating otp", zap.Error(err))
				return nil, err
			}
			return &commonsv1.AuthValidateOtpResponse{
				Valid: true,
			}, nil
		}
	}
}
