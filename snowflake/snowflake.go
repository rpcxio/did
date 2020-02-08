package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// A Node struct holds the basic information needed for a snowflake generator
// node
type Node struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewNode(node int64, epoch int64, nodeBits, stepBits uint8) (*Node, error) {
	n := Node{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << nodeBits)
	n.nodeMask = n.nodeMax << stepBits
	n.stepMask = -1 ^ (-1 << stepBits)
	n.timeShift = nodeBits + stepBits
	n.nodeShift = stepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(n.nodeMax, 10))
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(curTime))

	return &n, nil
}

// Generate creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Node) Generate() int64 {
	n.mu.Lock()

	now := time.Since(n.epoch).Nanoseconds() / 1000000
	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}
	n.time = now

	r := (now)<<n.timeShift |
		(n.node << n.nodeShift) |
		(n.step)

	n.mu.Unlock()
	return r
}

// GenerateBatch creates and returns a batch of unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Node) GenerateBatch(c uint16) []int64 {
	n.mu.Lock()

	var rt = make([]int64, 0, c)
	for c > 0 {
		generated := n.generateBatchInCurrentTime(c, &rt)
		c = c - generated
	}
	n.mu.Unlock()
	return rt
}

func (n *Node) generateBatchInCurrentTime(c uint16, rt *[]int64) uint16 {
	var generated = c
	var startStep = n.step
	now := time.Since(n.epoch).Nanoseconds() / 1000000
	if now == n.time {
		startStep = (n.step + 1) & n.stepMask
		if startStep == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
			n.step = 0
		}

	} else {
		n.step = 0
		startStep = 0
	}
	n.time = now

	if n.step+int64(c) > n.stepMask {
		generated = uint16(n.stepMask - n.step)
	}

	n.step = (n.step + int64(generated)) & n.stepMask

	for i := int64(0); i < int64(generated); i++ {
		r := (now)<<n.timeShift |
			(n.node << n.nodeShift) |
			(startStep + i)
		*rt = append(*rt, r)
	}

	return generated
}
