package main

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

const (
	testToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM1OTE1NDIsImlkIjoidGVzdCJ9.92PTS_X6shnvL0CmYoV5ms2UhLs7LyOKdux4787UFBTMSrpJLhsA8JhWg_CNw0UHdgaLkZzRXGG0L2jeJxnW_ylMnHe7WnyEHdLfj8iXT-uVaNTNU9VgofzciWuw8IFJoWlj1Rn20tuaIEEtIYqHGEC-Gxtu0-5S3dLdhDqVNIds0zSWc7fA2y55a5_PhRSdqeoUpopl2aDUxpRo8CaY7HNxlm_vq4Vi3xhNmAMaJlAXKnbYJwY7RysrVbKA52DP0ABDaAU3EghN8jtkatRLJOSp4RP3lC5QKvyL2hqquy25_N9QKGkodCuFx0K4CxUwt3EcGTvS_SQagqix0y2OgA"
)

func mockToken() context.Context {
	md := metadata.Pairs("token", testToken)
	return metadata.NewIncomingContext(context.Background(), md)
}

func prepareServerWithMocks(t *testing.T) (*Server, sqlmock.Sqlmock, *MockProducer) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	pubBytes, err := ioutil.ReadFile("signature.pub")
	if err != nil {
		t.Fatalf("failed to read public key: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}

	producer := &MockProducer{}

	return &Server{
		db:        db,
		producer:  producer,
		jwtPublic: pubKey,
	}, mock, producer
}

type MockProducer struct {
	Messages []*kafka.Message
}

func (p *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
	p.Messages = append(p.Messages, msg)
	return nil
}

func TestCreatePost(t *testing.T) {
	server, mock, _ := prepareServerWithMocks(t)
	ctx := mockToken()

	mock.ExpectExec(`INSERT INTO "posts"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "test", "Test title", "Test text").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := &CreatePostRequest{Title: "Test title", Text: "Test text"}
	_, err := server.CreatePost(ctx, r)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePost(t *testing.T) {
	server, mock, _ := prepareServerWithMocks(t)
	ctx := mockToken()
	postID := "post-id-123"

	mock.ExpectQuery(`SELECT author FROM "posts"`).
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"author"}).AddRow("test"))

	mock.ExpectExec(`DELETE FROM "posts"`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM "comments"`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := &DeletePostRequest{Id: postID}
	_, err := server.DeletePost(ctx, r)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLikePost(t *testing.T) {
	server, mock, producer := prepareServerWithMocks(t)
	ctx := mockToken()
	postID := "post-id-456"

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM "likes"`).
		WithArgs(postID, "test").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(`INSERT INTO "likes"`).
		WithArgs(sqlmock.AnyArg(), "test", postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectQuery(`SELECT author FROM "posts"`).
		WithArgs(postID).
		WillReturnRows(sqlmock.NewRows([]string{"author"}).AddRow("test"))

	r := &LikePostRequest{PostId: postID}
	_, err := server.LikePost(ctx, r)
	assert.NoError(t, err)
	assert.Len(t, producer.Messages, 1)
}
