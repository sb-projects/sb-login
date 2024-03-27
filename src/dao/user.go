package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/sb-projects/sb-login/src/literals"
	"github.com/sb-projects/sb-login/src/models"
)

// RegisterUser implements Daolayer.
func (dao *Dao) RegisterUser(ctx context.Context, user models.RegisterUserReq) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), literals.DefaultTimeout)
	defer cancel()

	var id int

	err := dao.pg.GetContext(ctx, &id, insertUser, user.Name, user.Email, user.HashPass, time.Now())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", id), nil
}
