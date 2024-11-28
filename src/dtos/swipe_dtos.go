package dtos

import "errors"

type SwipePartnerParams struct {
	UserID       int64
	TargetUserID int64  `json:"target_user_id"`
	Action       string `json:"action"`
}

func (s *SwipePartnerParams) Validate() (err error) {
	if s.UserID == s.TargetUserID {
		return errors.New("tidak dapat melakukan swipe pada diri sendiri")
	}

	if s.TargetUserID <= 0 {
		return errors.New("partner tidak valid")
	}

	if s.Action != "pass" && s.Action != "like" {
		return errors.New("hanya dapat melakukan swipe right ataupun left")
	}

	return nil
}
