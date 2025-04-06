package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func connectToDB() (*sql.DB, error) {
	fmt.Println("Connecting to database...")
	connStr := "host=posts_db port=5433 user=posts password=password dbname=postsdb sslmode=disable"
	return sql.Open("postgres", connStr)
}

func initDB(db *sql.DB) {
	fmt.Println("Init tables...")

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			id          TEXT PRIMARY KEY,
    		created     TIMESTAMP NOT NULL,
    		author      TEXT NOT NULL,
    		title       TEXT,
			text        TEXT
		);`)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init table:", err)
		os.Exit(1)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "comments" (
			id          TEXT PRIMARY KEY,
			post_id		TEXT,
			created     TIMESTAMP NOT NULL,
			author      TEXT NOT NULL,
			text        TEXT
		);`)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to init table:", err)
		os.Exit(1)
	}
}

type Server struct {
	UnimplementedPostsServer
	db        *sql.DB
	jwtPublic *rsa.PublicKey
}

func (s *Server) CreatePost(ctx context.Context, req *CreatePostRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata in context")
	}

	token, err := jwt.Parse(md.Get("token")[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.jwtPublic, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	postId := uuid.New().String()

	_, err = s.db.Exec(`
		INSERT INTO "posts" (id, created, author, title, text) 
		VALUES ($1, $2, $3, $4, $5);
	`, postId, time.Now().UTC(), userId, req.GetTitle(), req.GetText())
	if err != nil {
		return nil, fmt.Errorf("error connecting with db: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateComment(ctx context.Context, req *CreateCommentRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata in context")
	}

	token, err := jwt.Parse(md.Get("token")[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.jwtPublic, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	commId := uuid.New().String()

	_, err = s.db.Exec(`
		INSERT INTO "comments" (id, post_id, created, author, text) 
		VALUES ($1, $2, $3, $4, $5);
	`, commId, req.GetPostId(), time.Now().UTC(), userId, req.GetText())
	if err != nil {
		return nil, fmt.Errorf("error connecting with db: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeletePost(ctx context.Context, req *DeletePostRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata in context")
	}

	token, err := jwt.Parse(md.Get("token")[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.jwtPublic, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	postId := req.GetId()

	var author string
	err = s.db.QueryRow(`
		SELECT author FROM "posts" WHERE id = $1
	`, postId).Scan(&author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("error querying post: %v", err)
	}

	if author != userId {
		return nil, fmt.Errorf("user is not the author of the post")
	}

	_, err = s.db.Exec(`
		DELETE FROM "posts" WHERE id = $1
	`, postId)
	if err != nil {
		return nil, fmt.Errorf("error deleting post: %v", err)
	}

	_, err = s.db.Exec(`
		DELETE FROM "comments" WHERE post_id = $1
	`, postId)
	if err != nil {
		return nil, fmt.Errorf("error deleting comments: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteComment(ctx context.Context, req *DeleteCommentRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata in context")
	}

	token, err := jwt.Parse(md.Get("token")[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}

		return s.jwtPublic, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	commentId := req.GetId()

	var author string
	err = s.db.QueryRow(`
		SELECT author FROM "comments" WHERE id = $1
	`, commentId).Scan(&author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("error querying comment: %v", err)
	}

	if author != userId {
		return nil, fmt.Errorf("user is not the author of the comment")
	}

	_, err = s.db.Exec(`
		DELETE FROM "comments" WHERE id = $1
	`, commentId)
	if err != nil {
		return nil, fmt.Errorf("error deleting comment: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ModifyPost(ctx context.Context, req *ModifyPostRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata in context")
	}

	token, err := jwt.Parse(md.Get("token")[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtPublic, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	var author string
	err = s.db.QueryRow(`
		SELECT author FROM "posts" WHERE id = $1
	`, req.GetId()).Scan(&author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("error querying post: %v", err)
	}

	if author != userId {
		return nil, fmt.Errorf("user is not the author of the post")
	}

	_, err = s.db.Exec(`
		UPDATE "posts" 
		SET title = $1, text = $2
		WHERE id = $3
	`, req.GetTitle(), req.GetText(), req.GetId())

	if err != nil {
		return nil, fmt.Errorf("error updating post: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ModifyComment(ctx context.Context, req *ModifyCommentRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no metadata in context")
	}

	token, err := jwt.Parse(md.Get("token")[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid signing method")
		}
		return s.jwtPublic, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userId, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	var author string
	err = s.db.QueryRow(`
		SELECT author FROM "comments" WHERE id = $1
	`, req.Id).Scan(&author)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("error querying comment: %v", err)
	}

	if author != userId {
		return nil, fmt.Errorf("user is not the author of the comment")
	}

	_, err = s.db.Exec(`
		UPDATE "comments" 
		SET text = $1
		WHERE id = $2
	`, req.Text, req.Id)

	if err != nil {
		return nil, fmt.Errorf("error updating comment: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetPosts(req *GetPostsRequest, stream Posts_GetPostsServer) error {
	rows, err := s.db.Query(`
		SELECT id, created, author, title, text 
		FROM "posts"
		WHERE author = $1
	`, req.UserId)
	if err != nil {
		return fmt.Errorf("error querying posts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post PostInfo
		var created time.Time

		err := rows.Scan(&post.Id, &created, &post.Author, &post.Title, &post.Text)
		if err != nil {
			return fmt.Errorf("error scanning post row: %v", err)
		}

		post.Created = timestamppb.New(created)

		if err := stream.Send(&post); err != nil {
			return fmt.Errorf("error sending post: %v", err)
		}
	}

	return nil
}

func (s *Server) GetComments(req *GetCommentsRequest, stream Posts_GetCommentsServer) error {
	rows, err := s.db.Query(`
		SELECT id, post_id, created, author, text
		FROM "comments"
		WHERE post_id = $1
	`, req.PostId)
	if err != nil {
		return fmt.Errorf("error querying comments: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comment CommentInfo
		var created time.Time

		err := rows.Scan(&comment.Id, &comment.PostId, &created, &comment.Author, &comment.Text)
		if err != nil {
			return fmt.Errorf("error scanning comment row: %v", err)
		}

		comment.Created = timestamppb.New(created)

		if err := stream.Send(&comment); err != nil {
			return fmt.Errorf("error sending comment: %v", err)
		}
	}

	return nil
}

func main() {
	port := 8093
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := &Server{}

	publicPath := "signature.pub"
	public, err := os.ReadFile(publicPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	server.jwtPublic, err = jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	server.db, err = connectToDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to the database:", err)
		os.Exit(1)
	}
	defer server.db.Close()

	initDB(server.db)

	RegisterPostsServer(grpcServer, server)

	fmt.Printf("Server is running on port :%d\n", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
