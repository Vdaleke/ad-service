package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework10/internal/ads"
	"homework10/internal/users"
)

func AdSuccessResponse(ad *ads.Ad) *AdResponse {
	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		AuthorId:  ad.AuthorID,
		Published: ad.Published,
		CreatedAt: timestamppb.New(ad.CreatedAt),
		UpdatedAt: timestamppb.New(ad.UpdatedAt),
	}
}

func AdsSuccessResponse(ads *[]ads.Ad) *ListAdResponse {
	var adsResponseData []*AdResponse
	for _, ad := range *ads {
		adsResponseData = append(adsResponseData, AdSuccessResponse(&ad))
	}

	return &ListAdResponse{
		List: adsResponseData,
	}
}

func UserSuccessResponse(user *users.User) *UserResponse {
	return &UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
