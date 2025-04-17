package service

import (
	"context"
	"time"

	"github.com/ganis/okblog/profile/pkg/model"
	"github.com/go-kit/log"
)

// Middleware describes a service middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw *loggingMiddleware) CreateProfile(ctx context.Context, req model.CreateProfileRequest) (profile *model.Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CreateProfile",
			"username", req.Username,
			"email", req.Email,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return mw.next.CreateProfile(ctx, req)
}

func (mw *loggingMiddleware) GetProfile(ctx context.Context, id string) (profile *model.Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetProfile",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return mw.next.GetProfile(ctx, id)
}

func (mw *loggingMiddleware) UpdateProfile(ctx context.Context, id string, req model.UpdateProfileRequest) (profile *model.Profile, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "UpdateProfile",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return mw.next.UpdateProfile(ctx, id, req)
}

func (mw *loggingMiddleware) DeleteProfile(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "DeleteProfile",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())

	return mw.next.DeleteProfile(ctx, id)
}
