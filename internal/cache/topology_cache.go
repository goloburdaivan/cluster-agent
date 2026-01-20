package cache

import (
	"cluster-agent/internal/services/graph"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	cacheKeyPrefix = "topology:"
	ttl            = time.Hour
)

var (
	ErrNotFound = errors.New("topology not found in cache")
)

type TopologyCache struct {
	redisClient *redis.Client
}

func NewTopologyCache(redisClient *redis.Client) *TopologyCache {
	return &TopologyCache{
		redisClient: redisClient,
	}
}

func (c *TopologyCache) Get(ctx context.Context, namespace string) (*graph.Graph, error) {
	var topology graph.Graph

	cacheKey := cacheKeyPrefix + namespace
	result, err := c.redisClient.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	err = json.Unmarshal(result, &topology)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal topology from cache: %w", err)
	}

	return &topology, nil
}

func (c *TopologyCache) Set(ctx context.Context, namespace string, topology *graph.Graph) error {
	bytes, err := json.Marshal(topology)
	if err != nil {
		return fmt.Errorf("failed to marshal topology to cache: %w", err)
	}

	cacheKey := cacheKeyPrefix + namespace
	err = c.redisClient.Set(ctx, cacheKey, bytes, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save topology to cache: %w", err)
	}

	return nil
}
