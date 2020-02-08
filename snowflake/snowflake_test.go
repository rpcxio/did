package snowflake

import (
	"testing"
	"time"
)

func TestNewNode(t *testing.T) {
	_, err := NewNode(0, 1580601600000, 10, 12)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	_, err = NewNode(5000, 1580601600000, 10, 12)
	if err == nil {
		t.Fatalf("no error creating NewNode, %s", err)
	}
}

// lazy check if Generate will create duplicate IDs
// would be good to later enhance this with more smarts
func TestGenerateDuplicateID(t *testing.T) {
	node, _ := NewNode(1, 1580601600000, 10, 12)

	var x, y int64
	for i := 0; i < 1000000; i++ {
		y = node.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

func TestGenerateBatch(t *testing.T) {
	node, _ := NewNode(1, 1580601600000, 10, 12)

	count := uint16(node.stepMask - 10)
	var x, y []int64
	for i := 0; i < 1000; i++ {
		y = node.GenerateBatch(count)
		if len(x) > 0 && x[count-1] >= y[0] {
			t.Errorf("x(%d) & y(%d) are overlapped", x[count-1], y[0])
		}
		for i := 0; i < int(count-1); i++ {
			if y[i] >= y[i+1] {
				t.Errorf("not incremental id: y[%d](%d), y[%d](%d)", i, y[i], i+1, y[i+1])
			}
		}
		x = y
	}

	count = uint16(100)
	for i := 0; i < 100000; i++ {
		y = node.GenerateBatch(count)
		if len(x) > 0 && x[count-1] >= y[0] {
			t.Errorf("x(%d) & y(%d) are overlapped", x[count-1], y[0])
		}
		for i := 0; i < int(count-1); i++ {
			if y[i] >= y[i+1] {
				t.Errorf("not incremental id: y[%d](%d), y[%d](%d)", i, y[i], i+1, y[i+1])
			}
		}
		x = y
	}
}

// I feel like there's probably a better way
func TestRace(t *testing.T) {
	node, _ := NewNode(1, 1580601600000, 10, 12)
	go func() {
		for i := 0; i < 1000000000; i++ {
			NewNode(1, 1580601600000, 10, 12)
		}
	}()

	for i := 0; i < 4000; i++ {
		node.Generate()
	}
}

func TestParse(t *testing.T) {
	node, _ := NewNode(1, 1580601600000, 10, 12)
	ids := node.GenerateBatch(100)
	id := ids[99]
	tt := node.ParseTime(id)
	t.Logf("id:%d, time:%v", id, time.Unix(tt/1000, tt%1000*1000000))
	t.Logf("id:%d, serverID:%d", id, node.ParseServerID(id))
	t.Logf("id:%d, step:%d", id, node.ParseStep(id))
}
