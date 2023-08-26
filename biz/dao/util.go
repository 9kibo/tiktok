package dao

import (
	"errors"
	"gorm.io/gorm"
	"tiktok/pkg/errno"
)

// ofExists only ofExists return true, otherwise return false, if it has db err return with err
func ofExists(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
func ofGet(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errno.NotExists
	}
	return err
}
func ofUpdate1(tx *gorm.DB) error {
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errno.Update
	}
	return nil
}
