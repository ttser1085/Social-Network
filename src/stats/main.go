package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DB interface {
	Query(ctx context.Context, query string, args ...any) (driver.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) driver.Row
	Exec(ctx context.Context, query string, args ...any) error
}

type Server struct {
	UnimplementedStatsServer
	consumer *kafka.Consumer
	db       DB
}

const maxEntries = 100

func (s *Server) GetPostStat(ctx context.Context, req *PostRequest) (*PostStat, error) {
	const query = `
		SELECT likes, comments, views
		FROM post_stats
		WHERE post = ?
		ORDER BY event_time DESC
		LIMIT 1
	`

	var likes, comments, views uint32

	row := s.db.QueryRow(ctx, query, req.PostId)
	err := row.Scan(&likes, &comments, &views)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &PostStat{
				Likes:    0,
				Comments: 0,
				Views:    0,
			}, nil
		}
		log.Printf("failed to get post stats: %v", err)
		return nil, err
	}

	return &PostStat{
		Likes:    int32(likes),
		Comments: int32(comments),
		Views:    int32(views),
	}, nil
}

func (s *Server) GetViesDynamic(ctx context.Context, req *PostRequest) (*ViewsDynamic, error) {
	const query = `
		SELECT views, event_time
		FROM post_stats
		WHERE post = ?
		ORDER BY event_time DESC
		LIMIT ?
	`

	rows, err := s.db.Query(ctx, query, req.PostId, maxEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	views := make([]int32, 0, maxEntries)
	times := make([]*timestamppb.Timestamp, 0, maxEntries)

	for rows.Next() {
		var v uint32
		var eventTime time.Time
		if err := rows.Scan(&v, &eventTime); err != nil {
			return nil, err
		}

		views = append(views, int32(v))
		times = append(times, timestamppb.New(eventTime))
	}

	for i, j := 0, len(views)-1; i < j; i, j = i+1, j-1 {
		views[i], views[j] = views[j], views[i]
		times[i], times[j] = times[j], times[i]
	}

	return &ViewsDynamic{
		Views: views,
		Times: times,
	}, nil
}

func (s *Server) GetLikesDynamic(ctx context.Context, req *PostRequest) (*LikesDynamic, error) {
	const query = `
		SELECT likes, event_time
		FROM post_stats
		WHERE post = ?
		ORDER BY event_time DESC
		LIMIT ?
	`

	rows, err := s.db.Query(ctx, query, req.PostId, maxEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	likes := make([]int32, 0, maxEntries)
	times := make([]*timestamppb.Timestamp, 0, maxEntries)

	for rows.Next() {
		var l uint32
		var eventTime time.Time
		if err := rows.Scan(&l, &eventTime); err != nil {
			return nil, err
		}
		likes = append(likes, int32(l))
		times = append(times, timestamppb.New(eventTime))
	}

	for i, j := 0, len(likes)-1; i < j; i, j = i+1, j-1 {
		likes[i], likes[j] = likes[j], likes[i]
		times[i], times[j] = times[j], times[i]
	}

	return &LikesDynamic{
		Likes: likes,
		Times: times,
	}, nil
}

func (s *Server) GetCommentsDynamic(ctx context.Context, req *PostRequest) (*CommentsDynamic, error) {
	const query = `
		SELECT comments, event_time
		FROM post_stats
		WHERE post = ?
		ORDER BY event_time DESC
		LIMIT ?
	`

	rows, err := s.db.Query(ctx, query, req.PostId, maxEntries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]int32, 0, maxEntries)
	times := make([]*timestamppb.Timestamp, 0, maxEntries)

	for rows.Next() {
		var c uint32
		var eventTime time.Time
		if err := rows.Scan(&c, &eventTime); err != nil {
			return nil, err
		}
		comments = append(comments, int32(c))
		times = append(times, timestamppb.New(eventTime))
	}

	for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
		comments[i], comments[j] = comments[j], comments[i]
		times[i], times[j] = times[j], times[i]
	}

	return &CommentsDynamic{
		Comments: comments,
		Times:    times,
	}, nil
}

func (s *Server) GetTopPosts(ctx context.Context, req *Attribute) (*TopPosts, error) {
	const topN = 10

	var attr string
	switch req.Attr {
	case "likes", "comments", "views":
		attr = req.Attr
	default:
		return nil, fmt.Errorf("invalid attribute: %s", req.Attr)
	}

	query := fmt.Sprintf(`
		SELECT post
		FROM (
			SELECT post, %s, event_time,
				row_number() OVER (PARTITION BY post ORDER BY event_time DESC) AS rn
			FROM post_stats
		) WHERE rn = 1
		ORDER BY %s DESC
		LIMIT ?
	`, attr, attr)

	rows, err := s.db.Query(ctx, query, topN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]string, 0, topN)
	for rows.Next() {
		var post string
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &TopPosts{Posts: posts}, nil
}

func (s *Server) GetTopUsers(ctx context.Context, req *Attribute) (*TopUsers, error) {
	const topN = 10

	var attr string
	switch req.Attr {
	case "likes", "comments", "views":
		attr = req.Attr
	default:
		return nil, fmt.Errorf("invalid attribute: %s", req.Attr)
	}

	query := fmt.Sprintf(`
		SELECT author, SUM(%s) AS total
		FROM (
			SELECT author, %s,
				row_number() OVER (PARTITION BY post ORDER BY event_time DESC) AS rn
			FROM post_stats
		)
		WHERE rn = 1
		GROUP BY author
		ORDER BY total DESC
		LIMIT ?
	`, attr, attr)

	rows, err := s.db.Query(ctx, query, topN)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]string, 0, topN)
	for rows.Next() {
		var author string
		var total uint64
		if err := rows.Scan(&author, &total); err != nil {
			return nil, err
		}
		users = append(users, author)
	}

	return &TopUsers{Users: users}, nil
}

func connectClickHouse() clickhouse.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "mysecret",
		},
		DialTimeout: time.Second * 5,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect to ClickHouse: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatalf("ClickHouse ping error: %v", err)
	}

	log.Println("Connected to ClickHouse")
	return conn
}

type KafkaMessage struct {
	PostID string `json:"post_id"`
	Author string `json:"author"`
	UserID string `json:"user_id,omitempty"`
	Time   string `json:"time"`
}

func (s *Server) processMessage(msg *kafka.Message) {
	topic := *msg.TopicPartition.Topic
	payload := msg.Value

	log.Printf("Received message from %s: %s", topic, string(payload))

	var km KafkaMessage
	if err := json.Unmarshal(payload, &km); err != nil {
		log.Printf("failed to parse message JSON: %v", err)
		return
	}

	eventTime, err := time.Parse(time.RFC3339, km.Time)
	if err != nil {
		log.Printf("failed to parse event time: %v", err)
		return
	}

	ctx := context.Background()

	var lastLikes, lastComments, lastViews uint32
	var lastEventTime time.Time
	var lastAuthor string

	row := s.db.QueryRow(ctx, `
		SELECT likes, comments, views, event_time, author
		FROM post_stats
		WHERE post = ?
		ORDER BY event_time DESC
		LIMIT 1`, km.PostID)

	err = row.Scan(&lastLikes, &lastComments, &lastViews, &lastEventTime, &lastAuthor)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("failed to get last stats: %v", err)
		return
	}

	if err == sql.ErrNoRows {
		lastLikes = 0
		lastComments = 0
		lastViews = 0
		lastAuthor = km.Author
	}

	switch topic {
	case "posts-likes":
		lastLikes++
	case "posts-comments":
		lastComments++
	case "posts-views":
		lastViews++
	default:
		log.Printf("Unknown topic: %s", topic)
		return
	}

	err = s.db.Exec(ctx, `
		INSERT INTO post_stats (post, author, likes, comments, views, event_time)
		VALUES (?, ?, ?, ?, ?, ?)`,
		km.PostID, km.Author, lastLikes, lastComments, lastViews, eventTime)
	if err != nil {
		log.Printf("failed to insert stats: %v", err)
		return
	}
}

func initDB(db DB) {
	query := `
	DROP TABLE IF EXISTS post_stats
	`
	ctx := context.Background()
	err := db.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	query = `
	CREATE TABLE post_stats (
		post String,
		author String,
		likes UInt32,
		comments UInt32,
		views UInt32,
		event_time DateTime
	) 
	ENGINE = MergeTree()
	ORDER BY (post, event_time)
	`
	ctx = context.Background()
	err = db.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func main() {
	port := 8095
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	db := connectClickHouse()

	server := &Server{
		consumer: nil,
		db:       db,
	}

	initDB(server.db)

	server.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092", "group.id": "stats-group"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer server.consumer.Close()

	err = server.consumer.SubscribeTopics([]string{"posts-comments", "posts-views", "posts-likes"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	go func() {
		for {
			msg, err := server.consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Kafka read error: %v\n", err)
				continue
			}

			log.Printf("Received message: %s\n", string(msg.Value))
			server.processMessage(msg)
		}
	}()

	grpcServer := grpc.NewServer()

	RegisterStatsServer(grpcServer, server)

	fmt.Printf("Server is running on port :%d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
