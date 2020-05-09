package pkg

import (
	"errors"
	"longhorn/proxy/internal/global"
	"strconv"
	"sync"
	"time"
)

const (
	encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"
	encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
)

var Generator *GeneratorSnowFlake

var decodeBase32Map [256]byte
var decodeBase58Map [256]byte

func NewSnowflake(config global.SnowflakeConfig) *GeneratorSnowFlake {
	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}

	nodeMax := int64(-1 ^ (-1 << config.NodeBits))
	generator := &GeneratorSnowFlake{
		Epoch:            config.Epoch,
		NodeID:           config.BaseNodeID,
		NodeCount:        config.NodeCount,
		NodeBits:         config.NodeBits,
		StepBits:         config.StepBits,
		nodeMax:          nodeMax,
		nodeMask:         nodeMax << config.StepBits,
		stepMask:         -1 ^ (-1 << config.StepBits),
		timeShift:        config.NodeBits + config.StepBits,
		nodeShift:        config.StepBits,
		currentNodeIndex: 0,
	}
	generator.initGenerator()

	return generator
}

type ID int64

// Int64 returns an int64 of the snowflake ID
func (f ID) Int64() int64 {
	return int64(f)
}

func (f ID) Uint64() uint64 {
	return uint64(f)
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// Base2 returns a string base2 of the snowflake ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// Base36 returns a base36 string of the snowflake ID
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// Base32 uses the z-base-32 character set but encodes and decodes similar
// to base58, allowing it to create an even smaller result string.
// NOTE: There are many different base32 implementations so becareful when
// doing any interoperation interop with other packages.
func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

type Node struct {
	mu        sync.Mutex
	time      int64
	node      int64
	step      int64
	generator *GeneratorSnowFlake
}

func (n *Node) Generate() ID {

	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixNano() / 1000000

	if n.time == now {
		n.step = (n.step + 1) & n.generator.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID((now-n.generator.Epoch)<<n.generator.timeShift |
		(n.node << n.generator.nodeShift) |
		(n.step),
	)

	return r
}

type GeneratorSnowFlake struct {
	NodeID    int64
	NodeCount int64
	Epoch     int64
	NodeBits  uint8
	StepBits  uint8
	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
	nodes     []*Node

	lock             sync.Mutex
	currentNodeIndex int
}

func (g *GeneratorSnowFlake) newNode(nodeID int64) (*Node, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	// re-calc in case custom NodeBits or StepBits were set
	g.nodeMax = -1 ^ (-1 << g.NodeBits)
	g.nodeMask = g.nodeMax << g.StepBits
	g.stepMask = -1 ^ (-1 << g.StepBits)
	g.timeShift = g.NodeBits + g.StepBits
	g.nodeShift = g.StepBits

	if nodeID < 0 || nodeID > g.nodeMax {
		return nil, errors.New("node节点数已达上限")
	}

	n := &Node{
		time:      0,
		node:      nodeID,
		step:      0,
		generator: g,
	}

	g.nodes = append(g.nodes, n)
	return n, nil
}

func (g *GeneratorSnowFlake) getNextNode() *Node {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.currentNodeIndex >= len(g.nodes)-1 {
		g.currentNodeIndex = 0
	} else {
		g.currentNodeIndex += 1
	}

	return g.nodes[g.currentNodeIndex]
}

func (g *GeneratorSnowFlake) initGenerator() (err error) {
	for i := 0; int64(i) < g.NodeCount; i++ {
		_, err = g.newNode(g.NodeID + int64(i))
		if err != nil {
			return
		}
	}
	return
}

func (g *GeneratorSnowFlake) GenerateUniqueID() (uint64, error) {
	node := g.getNextNode()
	id := node.Generate()
	return id.Uint64(), nil
}
