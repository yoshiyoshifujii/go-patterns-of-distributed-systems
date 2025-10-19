package chapter_20

import (
	"strconv"
	"sync"
)

type Member struct {
	address string
}

type MembershipService struct {
	liveMembers []Member
}

type ClusterCoordinator struct {
	membership MembershipService

	mu      sync.Mutex
	counter uint64
}

func (m Member) getAddress() string {
	return m.address
}

func (m MembershipService) getLiveMembers() []Member {
	return append([]Member(nil), m.liveMembers...)
}

func (c *ClusterCoordinator) createPartitionTableFor(splits []string) PartitionTable {
	ranges := createRangesFromSplitPoints(splits)
	return c.arrangePartitions(ranges, c.membership.getLiveMembers())
}

func (c *ClusterCoordinator) arrangePartitions(ranges []Range, liveMembers []Member) PartitionTable {
	partitionTable := PartitionTable{
		partitions: make(map[string]PartitionInfo, len(ranges)),
	}

	if len(liveMembers) == 0 {
		return partitionTable
	}

	for i, rng := range ranges {
		member := liveMembers[i%len(liveMembers)]
		partitionID := c.newPartitionID()
		partitionInfo := newPartitionInfo(partitionID, member.getAddress(), PartitionStatusAssigned, rng)
		partitionTable.addPartition(partitionID, partitionInfo)
	}

	return partitionTable
}

func (c *ClusterCoordinator) newPartitionID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter++
	return "partition-" + strconv.FormatUint(c.counter, 10)
}
