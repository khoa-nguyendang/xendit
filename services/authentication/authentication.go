package services

import (
	"context"
	"errors"
	"net/http"
	"time"
	vm "xendit/view_models/login"
)

func (s *service) LogIn(ctx context.Context, model *vm.LogInRequest) (*vm.LogInResponse, error) {
	s.logger.Infof("LogIn triggered: %#v", model)
	span := s.tracer.StartSpan("AuthenticationService.LogIn")
	defer span.Finish()

	if model == nil {
		s.logger.Info("LogIn empty model")
		return nil, errors.New("LogIn empty model")
	}

	userEntity, err := s.repository.GetUser(ctx, model)
	if err != nil {
		return nil, err
	}

	token, err := s.jwtManager.GenerateTokenForUser(userEntity)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(model.UserName)
	token_expiry := time.Now().UTC().Add(time.Duration(s.cfg.Server.RefreshTokenExpireInHour * int(time.Hour)))
	refresh_token_expiry := time.Now().UTC().Add(time.Duration(s.cfg.Server.RefreshTokenExpireInHour * int(time.Hour)))
	if err != nil {
		return nil, err
	}

	return &vm.LogInResponse{
		Code:               http.StatusOK,
		Message:            "",
		Error:              "",
		Token:              token,
		RefreshToken:       refreshToken,
		TokenExpiry:        token_expiry.UnixMilli(),
		RefreshTokenExpiry: refresh_token_expiry.UnixMilli(),
	}, nil
}
