package chapter_20

import (
	"strconv"
	"sync"
)

type (
	RangeKey string

	Range struct {
		startKey RangeKey
		endKey   RangeKey
	}

	PartitionStatus string

	PartitionInfo struct {
		id            string
		memberAddress string
		status        PartitionStatus
		keyRange      Range
	}

	PartitionTable struct {
		partitions map[string]PartitionInfo
	}

	Member struct {
		address string
	}

	MembershipService struct {
		liveMembers []Member
	}

	ClusterCoordinator struct {
		membership MembershipService
	}
)

var (
	RangeMinKey                             = RangeKey("")
	RangeMaxKey                             = RangeKey("\uffff")
	PartitionStatusAssigned PartitionStatus = "Assigned"

	partitionCounter   uint64
	partitionCounterMu sync.Mutex
)

func newRange(startKey, endKey RangeKey) Range {
	return Range{startKey: startKey, endKey: endKey}
}

func newPartitionInfo(partitionId string, memberAddress string, status PartitionStatus, rng Range) PartitionInfo {
	return PartitionInfo{
		id:            partitionId,
		memberAddress: memberAddress,
		status:        status,
		keyRange:      rng,
	}
}

func (m Member) getAddress() string {
	return m.address
}

func (m MembershipService) getLiveMembers() []Member {
	return append([]Member(nil), m.liveMembers...)
}

func (pt *PartitionTable) addPartition(partitionId string, partitionInfo PartitionInfo) {
	if pt.partitions == nil {
		pt.partitions = make(map[string]PartitionInfo)
	}
	pt.partitions[partitionId] = partitionInfo
}

func (c ClusterCoordinator) createPartitionTableFor(splits []string) PartitionTable {
	ranges := createRangesFromSplitPoints(splits)
	return arrangePartitions(ranges, c.membership.getLiveMembers())
}

func createRangesFromSplitPoints(splits []string) []Range {
	var ranges []Range
	startKey := RangeMinKey
	for _, split := range splits {
		rangeKey := RangeKey(split)
		ranges = append(ranges, newRange(startKey, rangeKey))
		startKey = rangeKey
	}
	ranges = append(ranges, newRange(startKey, RangeMaxKey))
	return ranges
}

func arrangePartitions(ranges []Range, liveMembers []Member) PartitionTable {
	partitionTable := PartitionTable{
		partitions: make(map[string]PartitionInfo, len(ranges)),
	}

	if len(liveMembers) == 0 {
		return partitionTable
	}

	for i, rng := range ranges {
		member := liveMembers[i%len(liveMembers)]
		partitionId := newPartitionId()
		partitionInfo := newPartitionInfo(partitionId, member.getAddress(), PartitionStatusAssigned, rng)
		partitionTable.addPartition(partitionId, partitionInfo)
	}

	return partitionTable
}

func newPartitionId() string {
	partitionCounterMu.Lock()
	defer partitionCounterMu.Unlock()
	partitionCounter++
	return "partition-" + strconv.FormatUint(partitionCounter, 10)
}
