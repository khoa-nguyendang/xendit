package repositories

import (
	"context"
	ent "xendit/entities"
	vm "xendit/view_models/login"

	"github.com/opentracing/opentracing-go"
)

//
func (r *authRepo) GetUser(ctx context.Context, model *vm.LogInRequest) (*ent.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AuthenticationRepository.GetUser")
	defer span.Finish()

	usr := &ent.User{}

	if err := r.db.QueryRowContext(ctx,
		FindUserByUsernameQuery,
		model.UserName,
	).Scan(
		&usr.Id,
		&usr.Username,
		&usr.Password,
	); err != nil || usr == nil {
		r.logger.Errorf("GetUser.err: %v", err)
		return nil, err
	}

	return usr, nil
}
