package snowflake

import (
	"context"
	"strconv"
)

// SnowFlake twitter snowflake backend.
// It is based on https://github.com/bwmarrin/snowflake and
// - remove some unused methods
// - add GenerateBatch methods
//
// **How it Works**.
// Each time you generate an ID, it works, like this.
//
// A timestamp with millisecond precision is stored using 41 bits of the ID.
// Then the NodeID is added in subsequent bits.
// Then the Sequence Number is added, starting at 0 and incrementing for each ID generated in the same millisecond. If you generate enough IDs in the same millisecond that the sequence would roll over or overfill then the generate function will pause until the next millisecond.
//
// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+
type SnowFlake struct {
	serverID string
	node     *Node
}

// NewSnowFlake created a new SnowFlake with given configs.
func NewSnowFlake(serverID int64, epoch int64, nodeBits, stepBits uint8) *SnowFlake {
	node, err := NewNode(serverID, epoch, nodeBits, stepBits)
	if err != nil {
		panic("can't create snowflake node: " + err.Error())
	}

	return &SnowFlake{
		serverID: strconv.FormatInt(serverID, 10),
		node:     node,
	}
}

// GetServerID returns server id.
// Ignores the request (dummy).
func (sf *SnowFlake) GetServerID(ctx context.Context, dummy string, serverID *string) error {
	*serverID = sf.serverID
	return nil
}

// Get gets a unique ID.
// Ignores the request (dummy).
func (sf *SnowFlake) Get(ctx context.Context, dummy int8, id *int64) error {
	*id = sf.node.Generate()
	return nil
}

// GetBatch gets a batch of ids.
func (sf *SnowFlake) GetBatch(ctx context.Context, count uint16, ids *[]int64) error {
	*ids = sf.node.GenerateBatch(count)
	return nil
}
