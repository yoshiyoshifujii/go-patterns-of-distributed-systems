package chapter_20

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

	Partition struct {
	}
)

var (
	RangeMinKey             = RangeKey("")
	RangeMaxKey             = RangeKey("\uffff")
	PartitionStatusAssigned = PartitionStatus("Assigned")
)

func (rk RangeKey) compareTo(other RangeKey) int {
	if rk < other {
		return -1
	}
	if rk > other {
		return 1
	}
	return 0
}

func newRange(startKey, endKey RangeKey) Range {
	return Range{startKey: startKey, endKey: endKey}
}

func (r Range) contains(key RangeKey) bool {
	return key.compareTo(r.startKey) >= 0 && (r.endKey == RangeMaxKey || r.endKey.compareTo(key) > 0)
}

func (r Range) isOverlapping(other Range) bool {
	return r.contains(other.startKey) || other.contains(r.startKey)
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

func newPartitionInfo(partitionID string, memberAddress string, status PartitionStatus, rng Range) PartitionInfo {
	return PartitionInfo{
		id:            partitionID,
		memberAddress: memberAddress,
		status:        status,
		keyRange:      rng,
	}
}

func (pi *PartitionInfo) getRange() Range {
	return pi.keyRange
}

func (pi *PartitionInfo) getPartitionID() string {
	return pi.id
}

func (pi *PartitionInfo) getAddress() string {
	return pi.memberAddress
}

func (pt *PartitionTable) addPartition(partitionID string, partitionInfo PartitionInfo) {
	if pt.partitions == nil {
		pt.partitions = make(map[string]PartitionInfo)
	}
	pt.partitions[partitionID] = partitionInfo
}

func (pt *PartitionTable) getPartitionsInRange(rng Range) []PartitionInfo {
	allPartitions := pt.partitions
	var partitionsInRange []PartitionInfo
	for _, partitionInfo := range allPartitions {
		if partitionInfo.getRange().isOverlapping(rng) {
			partitionsInRange = append(partitionsInRange, partitionInfo)
		}
	}
	return partitionsInRange
}

func (p *Partition) getAllInRange(rng Range) []string {
	panic("implement me")
}
