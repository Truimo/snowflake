package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

var (
	epoch    int64 = 1715231857029
	nodeBits uint8 = 10
	stepBits uint8 = 12

	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	nodeMask        = nodeMax << stepBits
	stepMask  int64 = -1 ^ (-1 << stepBits)
	timeShift       = nodeBits + stepBits
	nodeShift       = stepBits
)

type Node struct {
	mu    sync.Mutex // 添加互斥锁 确保并发安全
	node  int64      // 该节点的ID
	time  int64      // 记录时间戳
	step  int64      // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
	epoch time.Time
}

type ID int64

func NewNode(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}
	n := Node{
		node: node,
	}
	var curTime = time.Now()
	n.epoch = curTime.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(curTime))
	return &n, nil
}

func (n *Node) Generate() ID {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID(now<<timeShift | (n.node << nodeShift) | n.step)
	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}

func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

func (f ID) Uint() uint {
	return uint(f)
}
