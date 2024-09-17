package templates

import (
	"context"
	"log"

	user_model "code.gitea.io/gitea/models/user"
)

type ConverterUtil struct {
	ctx context.Context
}

func NewConverterUtils(ctx context.Context) *ConverterUtil {
	return &ConverterUtil{ctx: ctx}
}

func (cu *ConverterUtil) UserIdToUser(id int64) *user_model.User {
	user, err := user_model.GetUserByID(cu.ctx, id)
	if err != nil {
		log.Fatalf("failed to load user %v", err)
		return nil
	}
	return user
}
