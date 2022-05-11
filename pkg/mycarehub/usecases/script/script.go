package main

import (
	"context"
	"errors"
	"log"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	"github.com/savannahghi/serverutils"
	"github.com/sirupsen/logrus"
	libGorm "gorm.io/gorm"
)

var (
	getStreamAPIKey    = serverutils.MustGetEnvVar("GET_STREAM_KEY")
	getStreamAPISecret = serverutils.MustGetEnvVar("GET_STREAM_SECRET")
)

func main() {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		log.Fatalf("can't instantiate repository in resolver: %v", err)
	}

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)

	streamClient, err := stream.NewClient(getStreamAPIKey, getStreamAPISecret)
	if err != nil {
		log.Fatalf("failed to start getstream client: %v", err)
	}

	getStream := streamService.NewServiceGetStream(streamClient)

	cleanUnknowns(getStream, db)
}

func cleanUnknowns(getStream streamService.ServiceGetStream, db *postgres.MyCareHubDb) {
	ctx := context.Background()

	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"role": "user",
		},
	}

	getStreamUserResponse, err := getStream.ListGetStreamUsers(ctx, query)
	if err != nil {
		log.Fatalf("failed to list getstream users: %v", err)
	}

	toRemove := []string{}

	for _, user := range getStreamUserResponse.Users {
		userTypeDeterminer := func(user *stream.User) string {
			_, err := db.GetStaffProfileByStaffID(ctx, user.ID)
			if err != nil {
				if errors.Is(err, libGorm.ErrRecordNotFound) {
					_, err := db.GetClientProfileByClientID(ctx, user.ID)
					if err != nil {
						toRemove = append(toRemove, user.ID)
						return ""
					}
					return "CLIENT"
				}
				toRemove = append(toRemove, user.ID)
				return ""
			}
			return "STAFF"
		}

		var metadata domain.MemberMetadata
		err := mapstructure.Decode(user.ExtraData, &metadata)
		if err != nil {
			logrus.Errorln("error decoding user extra data", err)
			continue
		}

		if metadata.UserType != "" && metadata.UserID != "" && metadata.Username != "" {
			continue
		}

		if metadata.UserType == "" {
			metadata.UserType = userTypeDeterminer(user)
		}

		var user *stream.User

		switch metadata.UserType {
		case "STAFF":
			staffProfile, err := db.GetStaffProfileByStaffID(ctx, user.ID)
			if err != nil {
				if errors.Is(err, libGorm.ErrRecordNotFound) {
					toRemove = append(toRemove, user.ID)
					continue
				}
				logrus.Errorln(err)
				continue
			}

			user = &stream.User{
				ID: *staffProfile.ID,
				ExtraData: map[string]interface{}{
					"username": staffProfile.User.Username,
					"userID":   staffProfile.User.ID,
					"userType": "STAFF",
				},
			}

		case "CLIENT":
			clientProfile, err := db.GetClientProfileByClientID(ctx, user.ID)
			if err != nil {
				if errors.Is(err, libGorm.ErrRecordNotFound) {
					toRemove = append(toRemove, user.ID)
					continue
				}
				logrus.Errorln(err)
				continue
			}

			user = &stream.User{
				ID: *clientProfile.ID,
				ExtraData: map[string]interface{}{
					"username": clientProfile.User.Username,
					"userID":   clientProfile.User.ID,
					"userType": "CLIENT",
				},
			}

		}

		log.Println("updating user details", user)
		_, err = getStream.UpsertUser(ctx, user)
		if err != nil {
			logrus.Errorln("error upserting staff", err)
		}
	}

	if len(toRemove) > 0 {
		log.Println("removing redundant users:", toRemove)
		opts := stream.DeleteUserOptions{
			User:          stream.HardDelete,
			Messages:      stream.HardDelete,
			Conversations: stream.HardDelete,
		}

		_, err = getStream.DeleteUsers(ctx, toRemove, opts)
		if err != nil {
			log.Fatalf("failed to delete getstream users: %v", err)
		}
	}

	log.Println("clean up completed")
}
