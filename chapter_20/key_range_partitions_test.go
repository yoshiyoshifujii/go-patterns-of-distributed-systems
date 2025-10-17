package chapter_20

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRangesFromSplitPoints(t *testing.T) {
	// Given
	splits := []string{"d", "m"}

	// When
	actual := createRangesFromSplitPoints(splits)

	// Then
	assert.Len(t, actual, len(splits)+1, "range count mismatch")
	assert.Equal(t, RangeMinKey, actual[0].startKey, "first range start")
	assert.Equal(t, RangeKey("d"), actual[0].endKey, "first range end")
	assert.Equal(t, RangeKey("d"), actual[1].startKey, "second range start")
	assert.Equal(t, RangeKey("m"), actual[1].endKey, "second range end")
	assert.Equal(t, RangeKey("m"), actual[2].startKey, "last range start")
	assert.Equal(t, RangeMaxKey, actual[2].endKey, "last range end")
}

func TestArrangePartitionsAssignsMembersRoundRobin(t *testing.T) {
	// Given
	ranges := []Range{
		newRange(RangeMinKey, RangeKey("h")),
		newRange(RangeKey("h"), RangeKey("t")),
		newRange(RangeKey("t"), RangeMaxKey),
	}
	members := []Member{
		{address: "node-a"},
		{address: "node-b"},
	}

	// When
	sut := &ClusterCoordinator{}
	actual := sut.arrangePartitions(ranges, members)

	// Then
	assert.Len(t, actual.partitions, len(ranges), "partition count mismatch")

	expectedAssignments := map[RangeKey]string{
		ranges[0].startKey: "node-a",
		ranges[1].startKey: "node-b",
		ranges[2].startKey: "node-a",
	}

	seenStarts := make(map[RangeKey]struct{})

	for _, info := range actual.partitions {
		assert.NotEmpty(t, info.id, "partition id should not be empty")
		assert.Equal(t, PartitionStatusAssigned, info.status, "partition status")

		expectedAddress, ok := expectedAssignments[info.keyRange.startKey]
		assert.Truef(t, ok, "unexpected range start key %q", info.keyRange.startKey)
		assert.Equalf(t, expectedAddress, info.memberAddress, "assigned member mismatch for range %q", info.keyRange.startKey)
		seenStarts[info.keyRange.startKey] = struct{}{}
	}

	assert.Len(t, seenStarts, len(expectedAssignments), "unexpected number of unique range starts")
}

func TestArrangePartitionsWithNoMembers(t *testing.T) {
	// Given
	ranges := createRangesFromSplitPoints([]string{"k"})

	// When
	sut := &ClusterCoordinator{}
	actual := sut.arrangePartitions(ranges, nil)

	// Then
	assert.Empty(t, actual.partitions, "partitions should be empty when no members live")
}
