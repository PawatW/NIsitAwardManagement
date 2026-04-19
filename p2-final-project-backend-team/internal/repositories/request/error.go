package request

import (
	"errors"

	"github.com/472-68-AgileDevOps/p2-final-project-backend-team/internal/domain/exceptions"
	"gorm.io/gorm"
)

func HandleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exceptions.ErrRequestNotFound
	}

	return err
}
