package domain

import (
	"database/sql"
	"encoding/json"
	"github.com/byteintellect/go_commons/entity"
	commonsv1 "github.com/byteintellect/protos_go/commons/v1"
)

type Otp struct {
	entity.BaseDomain
	Phone string
	Token string
}

func (o *Otp) GetTable() entity.DomainName {
	return entity.DomainName("otps")
}

func (o *Otp) ToDto() interface{} {
	return &commonsv1.AuthOtpDto{
		ExternalId:  o.ExternalId,
		Status:      commonsv1.Status(o.Status),
		PhoneNumber: o.Phone,
		Otp:         o.Token,
	}
}

func (o *Otp) FromDto(dto interface{}) (entity.Base, error) {
	otpDto := dto.(commonsv1.AuthOtpDto)
	o.Phone = otpDto.PhoneNumber
	o.Token = otpDto.Otp
	return o, nil
}

func (o *Otp) Merge(other interface{}) {
	otherOtp := other.(*Otp)
	o.Status = otherOtp.Status
}

func (o *Otp) FromSqlRow(rows *sql.Rows) (entity.Base, error) {
	var err error
	for rows.Next() {
		err = rows.Scan(&o.Id, &o.CreatedAt, &o.UpdatedAt, &o.DeletedAt, &o.ExternalId, &o.Phone, &o.Token)
	}
	return o, err
}

func (o *Otp) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Otp) UnmarshalBinary(buffer []byte) error {
	return json.Unmarshal(buffer, o)
}

func (o *Otp) Invalidate() {
	o.Status = int(commonsv1.Status_STATUS_INVALID)
}

func NewOtp(phone, token string) *Otp {
	return &Otp{
		Token: token,
		Phone: phone,
		BaseDomain: entity.BaseDomain{
			Status: 1,
		},
	}
}
