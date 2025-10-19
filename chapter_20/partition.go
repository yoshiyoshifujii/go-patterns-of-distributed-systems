package chapter_20

type RangeKey string

type Range struct {
	startKey RangeKey
	endKey   RangeKey
}

type PartitionStatus string

type PartitionInfo struct {
	id            string
	memberAddress string
	status        PartitionStatus
	keyRange      Range
}

type PartitionTable struct {
	partitions map[string]PartitionInfo
}

var (
	RangeMinKey             = RangeKey("")
	RangeMaxKey             = RangeKey("\uffff")
	PartitionStatusAssigned = PartitionStatus("Assigned")
)

func newRange(startKey, endKey RangeKey) Range {
	return Range{startKey: startKey, endKey: endKey}
}

func newPartitionInfo(partitionID string, memberAddress string, status PartitionStatus, rng Range) PartitionInfo {
	return PartitionInfo{
		id:            partitionID,
		memberAddress: memberAddress,
		status:        status,
		keyRange:      rng,
	}
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

func (pt *PartitionTable) addPartition(partitionID string, partitionInfo PartitionInfo) {
	if pt.partitions == nil {
		pt.partitions = make(map[string]PartitionInfo)
	}
	pt.partitions[partitionID] = partitionInfo
}
