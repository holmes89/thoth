package database

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/holmes89/thoth/internal"
	"github.com/rs/zerolog/log"
)

func (conn *conn) ListGames(ctx context.Context) ([]internal.Game, error) {

	params := &dynamodb.ScanInput{
		TableName:        tableName,
		FilterExpression: aws.String("#DYNOBASE_SK = :SK"),
		ExpressionAttributeNames: map[string]string{
			"#DYNOBASE_SK": "SK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":SK": &types.AttributeValueMemberS{Value: "game"},
		},
	}
	resp, err := conn.db.Scan(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("unable to find games")
		return nil, errors.New("unable to fetch games")
	}
	var games []internal.Game
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &games); err != nil {
		log.Error().Err(err).Msg("unable to unmarshal games")
		return nil, errors.New("failed to scan games")
	}

	for _, g := range games {
		g.Path = nil
	}
	log.Info().Int("entries", len(games)).Msg("found games")
	return games, nil
}

func (conn *conn) FindGame(ctx context.Context, id string) (internal.Game, error) {

	params := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
			"SK": &types.AttributeValueMemberS{Value: "game"},
		},
		TableName: tableName,
	}
	resp, err := conn.db.GetItem(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("unable to find game")
		return internal.Game{}, errors.New("unable to fetch game")
	}
	var game internal.Game
	if err := attributevalue.UnmarshalMap(resp.Item, &game); err != nil {
		log.Error().Err(err).Msg("unable to unmarshal games")
		return internal.Game{}, errors.New("failed to scan game")
	}

	obj, err := conn.bucket.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: bucketName,
		Key:    game.Path,
	})

	if err != nil {
		log.Error().Err(err).Msg("unable to get object url")
		return internal.Game{}, errors.New("failed to get object url")
	}

	game.Path = &obj.URL
	return game, nil
}
