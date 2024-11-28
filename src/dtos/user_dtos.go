package dtos

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type RegisterParam struct {
	Email     string `json:"email"  valid:"required~Email wajib diisi,email~Format email tidak sesuai"`
	Password  string `json:"password" valid:"required~Password wajib diisi,minstringlength(8)~Password minimal memiliki 8 karakter"`
	Sex       string `json:"sex" valid:"required~Jenis kelamin wajib dipilih"`
	FirstName string `json:"first_name"  valid:"required~First Name wajib diisi"`
	LastName  string `json:"last_name"`
}

func (r *RegisterParam) Validate() (err error) {
	if _, err = govalidator.ValidateStruct(r); err != nil {
		return
	}

	sex := map[string]bool{
		"male":   true,
		"female": true,
	}

	if !sex[r.Sex] {
		return errors.New("mohon untuk jenis kelamin antara male atau female")
	}

	return
}

type LoginParam struct {
	Email    string `json:"email" valid:"required~Email wajib diisi,email~Format email tidak sesuai"`
	Password string `json:"password" valid:"required~Password wajib diisi,minstringlength(8)~Password minimal memiliki 8 karakter"`
}

func (l *LoginParam) Validate() (err error) {
	_, err = govalidator.ValidateStruct(l)
	return
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type VerifyEmailParam struct {
	Email string `json:"email" valid:"required~Email wajib diisi,email~Format email tidak sesuai"`
}

func (v *VerifyEmailParam) Validate() (err error) {
	_, err = govalidator.ValidateStruct(v)
	return
}
