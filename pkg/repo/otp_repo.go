package repo

import (
	"context"
	"errors"
	"github.com/byteintellect/go_commons/db"
	"github.com/byteintellect/otp_svc/pkg/domain"
	"gorm.io/gorm"
)

type OtpRepo interface {
	db.BaseRepository
	GetOtp(ctx context.Context, externalId, phone, token string) (*domain.Otp, error)
	InvalidateOtp(ctx context.Context, otp *domain.Otp) error
}

type otpGormRepo struct {
	db.BaseRepository
}

func (ogr *otpGormRepo) GetOtp(ctx context.Context, externalId, phone, token string) (*domain.Otp, error) {
	db := ogr.GetDb().(*gorm.DB)
	var res domain.Otp
	if err := db.WithContext(ctx).Model(&res).Where("external_id = ? AND phone = ? AND token = ? AND status = 1", externalId, phone, token).Find(&res).Error; err != nil || res.Id == 0 {
		return nil, errors.New("invalid otp")
	}
	return &res, nil
}

func (ogr *otpGormRepo) InvalidateOtp(ctx context.Context, otp *domain.Otp) error {
	db := ogr.GetDb().(*gorm.DB)
	if err := db.WithContext(ctx).Model(otp).Where("external_id = ?", otp.ExternalId).Updates(map[string]interface{}{"status": "0"}).Error; err != nil {
		return err
	}
	return nil
}

func NewOtpRepo(db db.BaseRepository) OtpRepo {
	return &otpGormRepo{
		db,
	}
}
