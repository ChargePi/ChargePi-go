package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/pkg/grpc"
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

	// Get all tags from the database
	tags := s.tagManager.GetTags()

	// Convert the tags to the gRPC response
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
		// Add the tag to the database
		err := s.tagManager.AddTag(tag.TagId, types.NewIdTagInfo(types.AuthorizationStatus(tag.Status)))
		if err != nil {
			response.Status = append(response.Status, "Failed")
			continue
		}

		response.Status = append(response.Status, "Success")
	}

	return response, nil
}

func (s *AuthService) RemoveAuthorizedCard(ctx context.Context, request *grpc.RemoveCardRequest) (*grpc.RemoveCardResponse, error) {
	response := &grpc.RemoveCardResponse{
		Status: "Failed",
	}

	// Remove the tag from the database
	err := s.tagManager.RemoveTag(request.GetTagId())
	if err != nil {
		return response, nil
	}

	response.Status = "Success"
	return response, nil
}

func (s *UserService) mustEmbedUnimplementedTagServer() {
}
