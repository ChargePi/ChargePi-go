package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/internal/api/grpc"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthService struct {
	grpc.UnimplementedTagServer
	tagManager auth.TagManager
}

func NewAuthService(tagManager auth.TagManager) *AuthService {
	return &AuthService{
		tagManager: tagManager,
	}
}

func (s *AuthService) GetAuthorizedCards(ctx context.Context, empty *empty.Empty) (*grpc.GetAuthorizedCardsResponse, error) {
	response := &grpc.GetAuthorizedCardsResponse{
		AuthorizedCards: []*grpc.AuthorizedCard{},
	}

	tags := s.tagManager.GetTags()

	for _, tag := range tags {
		var timestamp *timestamppb.Timestamp
		if tag.IdTagInfo.ExpiryDate != nil {
			timestamp = timestamppb.New(tag.IdTagInfo.ExpiryDate.Time)
		}

		card := &grpc.AuthorizedCard{
			TagId:      tag.IdTag,
			Status:     string(tag.IdTagInfo.Status),
			ExpiryDate: timestamp,
		}
		response.AuthorizedCards = append(response.AuthorizedCards, card)
	}

	return response, nil
}

func (s *AuthService) AddAuthorizedCards(ctx context.Context, request *grpc.AddAuthorizedCardsRequest) (*grpc.AddAuthorizedCardsResponse, error) {
	response := &grpc.AddAuthorizedCardsResponse{Status: []string{}}

	for _, tag := range request.GetAuthorizedCards() {
		err := s.tagManager.AddTag(tag.TagId, types.NewIdTagInfo(types.AuthorizationStatus(tag.Status)))
		if err != nil {
			response.Status = append(response.Status, "Fail")
			continue
		}

		response.Status = append(response.Status, "Success")
	}

	return response, nil
}

func (s *AuthService) RemoveAuthorizedCard(ctx context.Context, request *grpc.RemoveCardRequest) (*grpc.RemoveCardResponse, error) {
	response := &grpc.RemoveCardResponse{
		Status: "Fail",
	}

	return response, nil
}

func (s *UserService) mustEmbedUnimplementedTagServer() {
}
