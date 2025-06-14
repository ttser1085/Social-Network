package main

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDB struct {
	mock.Mock
}

func (m *mockDB) Query(ctx context.Context, query string, args ...any) (driver.Rows, error) {
	call := m.Called(ctx, query, args)
	rows, _ := call.Get(0).(driver.Rows)
	return rows, call.Error(1)
}

func (m *mockDB) QueryRow(ctx context.Context, query string, args ...any) driver.Row {
	call := m.Called(ctx, query, args)
	return call.Get(0).(driver.Row)
}

func (m *mockDB) Exec(ctx context.Context, query string, args ...any) error {
	call := m.Called(ctx, query, args)
	return call.Error(0)
}

type mockRow struct {
	mock.Mock
}

func (m *mockRow) Scan(dest ...any) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *mockRow) ScanStruct(dest any) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *mockRow) Columns() ([]string, error) {
	return nil, nil
}

func (m *mockRow) ColumnTypes() ([]driver.ColumnType, error) {
	return nil, nil
}

func (m *mockRow) Err() error {
	return nil
}

type mockRows struct {
	mock.Mock
	index int
	data  [][]any
}

func (m *mockRows) Next() bool {
	if m.index < len(m.data) {
		return true
	}
	return false
}

func (m *mockRows) Scan(dest ...any) error {
	if m.index >= len(m.data) {
		return errors.New("no more rows")
	}
	for i := range dest {
		switch v := dest[i].(type) {
		case *uint32:
			*v = m.data[m.index][i].(uint32)
		case *uint64:
			*v = m.data[m.index][i].(uint64)
		case *string:
			*v = m.data[m.index][i].(string)
		case *time.Time:
			*v = m.data[m.index][i].(time.Time)
		default:
			return errors.New("unsupported scan type")
		}
	}
	m.index++
	return nil
}

func (m *mockRows) Close() error {
	return nil
}

func (m *mockRows) Err() error {
	return nil
}

func TestGetPostStatEmpty(t *testing.T) {
	mockDB := new(mockDB)
	mockRow := new(mockRow)

	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("no rows"))

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	s := &Server{db: mockDB}

	resp, err := s.GetPostStat(context.Background(), &PostRequest{PostId: "empty"})

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestGetPostStat(t *testing.T) {
	mockDB := new(mockDB)
	mockRow := new(mockRow)

	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(0).(*uint32) = 10
		*args.Get(1).(*uint32) = 5
		*args.Get(2).(*uint32) = 100
	}).Return(nil)

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

	s := &Server{db: mockDB}

	resp, err := s.GetPostStat(context.Background(), &PostRequest{PostId: "abc123"})
	assert.NoError(t, err)
	assert.Equal(t, int32(10), resp.Likes)
	assert.Equal(t, int32(5), resp.Comments)
	assert.Equal(t, int32(100), resp.Views)

	mockRow.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}

func TestProcessKafkaMessage(t *testing.T) {
	mockDB := new(mockDB)
	mockRow := new(mockRow)

	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		*args.Get(0).(*uint32) = 0
		*args.Get(1).(*uint32) = 0
		*args.Get(2).(*uint32) = 0
		*args.Get(3).(*time.Time) = time.Now()
		*args.Get(4).(*string) = "author1"
	}).Return(nil)

	mockDB.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	s := &Server{db: mockDB}

	msg := &KafkaMessage{
		PostID: "abc123",
		Author: "author1",
		Time:   time.Now().Format(time.RFC3339),
	}

	jsonBytes, _ := json.Marshal(msg)

	s.processMessage(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: ptr("posts-likes")},
		Value:          jsonBytes,
	})

	mockDB.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func ptr(s string) *string {
	return &s
}
