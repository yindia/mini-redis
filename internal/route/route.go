package route

import (
	"context"
	"fmt"
	"log"
	"strconv"

	v1 "redis/internal/gen/cloud/v1"
	"redis/internal/store"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
)

//go:generate mockery --output=./mocks --case=underscore --all --with-expecter
type RedisServiceHandler interface {
	Get(ctx context.Context, req *connect.Request[v1.GetRequest]) (*connect.Response[v1.GetResponse], error)
	Set(ctx context.Context, req *connect.Request[v1.SetRequest]) (*connect.Response[v1.SetResponse], error)
	Del(ctx context.Context, req *connect.Request[v1.DelRequest]) (*connect.Response[v1.DelResponse], error)
	Incr(ctx context.Context, req *connect.Request[v1.IncrRequest]) (*connect.Response[v1.IncrResponse], error)
	Expire(ctx context.Context, req *connect.Request[v1.ExpireRequest]) (*connect.Response[v1.ExpireResponse], error)
	Ping(ctx context.Context, req *connect.Request[v1.PingRequest]) (*connect.Response[v1.PingResponse], error)
	Backup(ctx context.Context, req *connect.Request[v1.BackupRequest]) (*connect.Response[v1.BackupResponse], error)
	Restore(ctx context.Context, req *connect.Request[v1.RestoreRequest]) (*connect.Response[v1.RestoreResponse], error)
	Join(ctx context.Context, req *connect.Request[v1.JoinRequest]) (*connect.Response[v1.JoinResponse], error)
}

// RedisServer represents the server handling Redis-like operations.
// It implements the v1.RedisServiceHandler interface.
type RedisServer struct {
	validator *protovalidate.Validator
	logger    *log.Logger
	store     *store.Store
}

// NewRedisServer creates and returns a new instance of RedisServer.
// It initializes the validator and sets up the logger.
func NewRedisServer(store *store.Store) RedisServiceHandler {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}

	server := &RedisServer{
		validator: validator,
		store:     store,
		logger:    log.New(log.Writer(), "RedisServer: ", log.LstdFlags|log.Lshortfile),
	}

	server.logger.Println("RedisServer initialized successfully")
	return server
}

// Get retrieves the value for a given key.
func (s *RedisServer) Get(ctx context.Context, req *connect.Request[v1.GetRequest]) (*connect.Response[v1.GetResponse], error) {
	if err := s.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	value, err := s.store.Get(req.Msg.Key)
	if err != nil {
		s.logger.Printf("Error getting key %s: %v", req.Msg.Key, err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetResponse{Value: []byte(value)}), nil
}

// Set stores a key-value pair.
func (s *RedisServer) Set(ctx context.Context, req *connect.Request[v1.SetRequest]) (*connect.Response[v1.SetResponse], error) {
	if err := s.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err := s.store.Set(req.Msg.Key, string(req.Msg.Value))
	if err != nil {
		s.logger.Printf("Error setting key %s: %v", req.Msg.Key, err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.SetResponse{Success: true}), nil
}

// Del deletes one or more keys.
func (s *RedisServer) Del(ctx context.Context, req *connect.Request[v1.DelRequest]) (*connect.Response[v1.DelResponse], error) {
	if err := s.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	deletedCount := 0
	for _, key := range req.Msg.Keys {
		err := s.store.Delete(key)
		if err != nil {
			s.logger.Printf("Error deleting key %s: %v", key, err)
		} else {
			deletedCount++
		}
	}

	return connect.NewResponse(&v1.DelResponse{DeletedCount: int32(deletedCount)}), nil
}

// Incr increments the integer value of a key.
func (s *RedisServer) Incr(ctx context.Context, req *connect.Request[v1.IncrRequest]) (*connect.Response[v1.IncrResponse], error) {
	if err := s.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	// This operation is not directly supported by the provided store interface.
	// We need to implement it using Get and Set operations.
	value, err := s.store.Get(req.Msg.Key)
	if err != nil {
		s.logger.Printf("Error getting key %s for increment: %v", req.Msg.Key, err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("value is not an integer"))
	}

	intValue++
	err = s.store.Set(req.Msg.Key, strconv.Itoa(intValue))
	if err != nil {
		s.logger.Printf("Error setting incremented value for key %s: %v", req.Msg.Key, err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.IncrResponse{Value: int64(intValue)}), nil
}

// Expire sets a timeout on a key.
func (s *RedisServer) Expire(ctx context.Context, req *connect.Request[v1.ExpireRequest]) (*connect.Response[v1.ExpireResponse], error) {
	// The provided store interface doesn't support expiration.
	// This would require additional implementation in the store.
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("expire operation not supported"))
}

// Ping checks if the server is responsive.
func (s *RedisServer) Ping(ctx context.Context, req *connect.Request[v1.PingRequest]) (*connect.Response[v1.PingResponse], error) {
	return connect.NewResponse(&v1.PingResponse{Message: "PONG"}), nil
}

// Backup creates a backup of the current dataset.
func (s *RedisServer) Backup(ctx context.Context, req *connect.Request[v1.BackupRequest]) (*connect.Response[v1.BackupResponse], error) {
	// The provided store interface doesn't support backup operations.
	// This would require additional implementation in the store.
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("backup operation not supported"))
}

// Restore rebuilds the dataset from a backup file.
func (s *RedisServer) Restore(ctx context.Context, req *connect.Request[v1.RestoreRequest]) (*connect.Response[v1.RestoreResponse], error) {
	// The provided store interface doesn't support restore operations.
	// This would require additional implementation in the store.
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("restore operation not supported"))
}

// Join adds a new node to the cluster
func (s *RedisServer) Join(ctx context.Context, req *connect.Request[v1.JoinRequest]) (*connect.Response[v1.JoinResponse], error) {
	if err := s.validator.Validate(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	err := s.store.Join(req.Msg.NodeId, req.Msg.RemoteAddr)
	if err != nil {
		s.logger.Printf("Error joining node %s at %s: %v", req.Msg.NodeId, req.Msg.RemoteAddr, err)
		return connect.NewResponse(&v1.JoinResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}), nil
	}

	return connect.NewResponse(&v1.JoinResponse{
		Success: true,
	}), nil
}
