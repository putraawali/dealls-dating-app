package repositories

import (
	"context"
	"dealls-dating-app/src/constants"
	"dealls-dating-app/src/models"
	"dealls-dating-app/src/pkg/response"
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

type UserRepository interface {
	Insert(ctx context.Context, user *models.User) (err error)
	FindByEmail(ctx context.Context, email string) (user models.User, err error)
	VerifyEmail(ctx context.Context, email string) (err error)
	FindByID(ctx context.Context, id int64) (user models.User, err error)
	FindBySex(ctx context.Context, sex string, limit, offset int) (users []models.User, err error)
	ActivatePremium(ctx context.Context, userID int64) (err error)
}

type userRepository struct {
	db       *gorm.DB
	response *response.Response
}

func NewUserRepository(di di.Container) UserRepository {
	return &userRepository{
		db:       di.Get(constants.PG_DB).(*gorm.DB),
		response: di.Get(constants.RESPONSE).(*response.Response),
	}
}

func (u *userRepository) Insert(ctx context.Context, user *models.User) (err error) {
	if err = u.db.Create(&user).WithContext(ctx).Error; err != nil {
		err = u.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (user models.User, err error) {
	if err = u.db.First(&user, "email = ?", email).WithContext(ctx).Error; err != nil {
		code := http.StatusInternalServerError
		msg := err

		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
			msg = fmt.Errorf("user with email %s not found", email)
		}

		err = u.response.NewError().
			SetContext(ctx).
			SetDetail(msg.Error()).
			SetMessage(msg).
			SetStatusCode(code)
	}

	return
}

func (u *userRepository) VerifyEmail(ctx context.Context, email string) (err error) {
	if err = u.db.Model(&models.User{}).Where("email = ?", email).Update("is_verified", true).WithContext(ctx).Error; err != nil {
		err = u.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}

func (u *userRepository) FindByID(ctx context.Context, id int64) (user models.User, err error) {
	if err = u.db.First(&user, "user_id = ?", id).WithContext(ctx).Error; err != nil {
		code := http.StatusInternalServerError
		msg := err

		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = http.StatusNotFound
			msg = errors.New("user not found")
		}

		err = u.response.NewError().
			SetContext(ctx).
			SetDetail(msg.Error()).
			SetMessage(msg).
			SetStatusCode(code)
	}

	return
}

func (u *userRepository) FindBySex(ctx context.Context, sex string, limit, offset int) (users []models.User, err error) {
	if err = u.db.
		Model(&models.User{}).
		Limit(limit).
		Offset(offset).
		Find(&users, "sex = ?", sex).
		WithContext(ctx).Error; err != nil {
		err = u.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}

func (u *userRepository) ActivatePremium(ctx context.Context, userID int64) (err error) {
	if err = u.db.Model(&models.User{}).Where("user_id = ?", userID).Update("is_premium", true).WithContext(ctx).Error; err != nil {
		err = u.response.NewError().
			SetContext(ctx).
			SetDetail(err.Error()).
			SetMessage(err).
			SetStatusCode(http.StatusInternalServerError)
	}

	return
}
