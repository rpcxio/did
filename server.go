package did

import "context"

// IDGen is the interface of ID generator.
// You can get one ID by `Get` or get some IDs by GetBatch().
type IDGen interface {
	GetServerID(ctx context.Context, dummy string, serverID *string) error
	Get(ctx context.Context, dummy int8, id *int64) error
	GetBatch(ctx context.Context, count uint16, ids *[]int64) error
}
