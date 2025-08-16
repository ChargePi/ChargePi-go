package grpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	tagsv1 "github.com/ChargePi/ChargePi-go/gen/proto/tags/v1"
	"github.com/ChargePi/ChargePi-go/internal/auth"
)

type AuthService struct {
	tagsv1.UnimplementedTagServiceServer
	tagManager auth.Service
}

func NewAuthService(tagManager auth.Service) *AuthService {
	return &AuthService{
		tagManager: tagManager,
	}
}

func (s *AuthService) GetAuthorizedCards(ctx context.Context, empty *empty.Empty) (*tagsv1.GetAuthorizedCardsResponse, error) {
	response := &tagsv1.GetAuthorizedCardsResponse{
		AuthorizedCards: []*tagsv1.AuthorizedCard{},
	}

	// Get all tags from the database
	tags, err := s.tagManager.GetTags()
	if err != nil {
		return response, nil
	}

	// Convert the tags to the gRPC response
	for _, tag := range tags {
		var timestamp *timestamppb.Timestamp
		if tag.IdTagInfo.ExpiryDate != nil {
			timestamp = timestamppb.New(tag.IdTagInfo.ExpiryDate.Time)
		}

		card := &tagsv1.AuthorizedCard{
			TagId:      tag.IdTag,
			Status:     string(tag.IdTagInfo.Status),
			ExpiryDate: timestamp,
		}
		response.AuthorizedCards = append(response.AuthorizedCards, card)
	}

	return response, nil
}

func (s *AuthService) AddAuthorizedCards(ctx context.Context, request *tagsv1.AddAuthorizedCardsRequest) (*tagsv1.AddAuthorizedCardsResponse, error) {
	response := &tagsv1.AddAuthorizedCardsResponse{Status: []string{}}

	for _, tag := range request.GetAuthorizedCards() {
		err := s.tagManager.CacheTag(tag.TagId, types.NewIdTagInfo(types.AuthorizationStatus(tag.Status)))
		if err != nil {
			response.Status = append(response.Status, "Failed")
			continue
		}

		response.Status = append(response.Status, "Success")
	}

	return response, nil
}

func (s *AuthService) RemoveAuthorizedCard(ctx context.Context, request *tagsv1.RemoveAuthorizedCardRequest) (*tagsv1.RemoveAuthorizedCardResponse, error) {
	response := &tagsv1.RemoveAuthorizedCardResponse{
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

func (s *UserHandler) mustEmbedUnimplementedTagServer() {
}
