package grpc

import (
	"context"
	"fmt"
	validator "github.com/Vdaleke/ad-validation"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework10/internal/app"
	"log"
	"os"
	"time"
)

func InterceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, _ logging.Level, msg string, fields ...any) {
		l := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

		msg = fmt.Sprintf("INFO :%v", msg)
		l.Println(append([]any{msg}, fields...))
	})
}

func PanicInterceptor(p any) (err error) {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func NewService(a app.App) AdServiceServer {
	return &AdService{adApp: a}
}

type AdService struct {
	adApp app.App
}

func (a *AdService) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := a.adApp.CreateAd(request.Title, request.Text, request.UserId)

	if errors.Is(err, validator.ValidationError) || errors.Is(err, app.DefunctUser) {
		return &AdResponse{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &AdResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return AdSuccessResponse(&ad), nil
}

func (a *AdService) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := a.adApp.ChangeAdStatus(request.AdId, request.UserId, request.Published)

	if errors.Is(err, app.PermissionDenied) {
		return &AdResponse{}, status.New(codes.PermissionDenied, "the user does not have permission to edit the ad").Err()
	} else if errors.Is(err, validator.ValidationError) ||
		errors.Is(err, app.DefunctUser) || errors.Is(err, app.DefunctAd) {
		return &AdResponse{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &AdResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return AdSuccessResponse(&ad), nil
}

func (a *AdService) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := a.adApp.UpdateAd(request.AdId, request.UserId, request.Title, request.Text)

	if errors.Is(err, app.PermissionDenied) {
		return &AdResponse{}, status.New(codes.PermissionDenied, "the user does not have permission to edit the ad").Err()
	} else if errors.Is(err, validator.ValidationError) ||
		errors.Is(err, app.DefunctUser) || errors.Is(err, app.DefunctAd) {
		return &AdResponse{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &AdResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return AdSuccessResponse(&ad), nil
}

func (a *AdService) GetAd(ctx context.Context, request *GetAdRequest) (*AdResponse, error) {
	ad, err := a.adApp.GetAd(request.Id)

	if errors.Is(err, app.DefunctAd) {
		return &AdResponse{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &AdResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return AdSuccessResponse(&ad), nil
}

func (a *AdService) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	err := a.adApp.DeleteAd(request.AdId, request.AuthorId)

	if errors.Is(err, app.PermissionDenied) {
		return &emptypb.Empty{}, status.New(codes.PermissionDenied, "the user does not have permission to edit the ad").Err()
	} else if errors.Is(err, validator.ValidationError) || errors.Is(err, app.DefunctUser) {
		return &emptypb.Empty{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &emptypb.Empty{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return &emptypb.Empty{}, nil
}

func (a *AdService) ListAds(ctx context.Context, request *ListAdsRequest) (*ListAdResponse, error) {
	timeFilter, _ := time.Parse(time.RFC3339, request.CreationTime)

	ads, err := a.adApp.ListAds(request.Published, request.UserId, timeFilter)

	if err != nil {
		return &ListAdResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return AdsSuccessResponse(&ads), nil
}

func (a *AdService) SearchAds(ctx context.Context, request *SearchAdsRequest) (*ListAdResponse, error) {
	ads, err := a.adApp.SearchAds(request.Pattern)

	if err != nil {
		return &ListAdResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return AdsSuccessResponse(&ads), nil
}

func (a *AdService) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	user, err := a.adApp.CreateUser(request.Name, request.Email)

	if err != nil {
		return &UserResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return UserSuccessResponse(&user), nil
}

func (a *AdService) UpdateUser(ctx context.Context, request *UpdateUserRequest) (*UserResponse, error) {
	user, err := a.adApp.UpdateUser(request.Id, request.Name, request.Email)

	if errors.Is(err, app.DefunctUser) {
		return &UserResponse{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &UserResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return UserSuccessResponse(&user), nil
}

func (a *AdService) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	user, err := a.adApp.GetUser(request.Id)

	if errors.Is(err, app.DefunctUser) {
		return &UserResponse{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &UserResponse{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return UserSuccessResponse(&user), nil
}

func (a *AdService) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	err := a.adApp.DeleteUser(request.Id)

	if errors.Is(err, app.DefunctUser) {
		return &emptypb.Empty{}, status.New(codes.InvalidArgument, "invalid information received").Err()
	} else if err != nil {
		return &emptypb.Empty{}, status.New(codes.Unknown, "an unknown error has occurred").Err()
	}

	return &emptypb.Empty{}, nil
}
