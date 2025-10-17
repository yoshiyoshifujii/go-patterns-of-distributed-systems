package chapter_20

import "testing"

func TestCreateRangesFromSplitPoints(t *testing.T) {
	splits := []string{"d", "m"}

	ranges := createRangesFromSplitPoints(splits)

	if len(ranges) != len(splits)+1 {
		t.Fatalf("expected %d ranges, got %d", len(splits)+1, len(ranges))
	}

	if ranges[0].startKey != RangeMinKey || ranges[0].endKey != RangeKey("d") {
		t.Errorf("unexpected first range: %+v", ranges[0])
	}

	if ranges[1].startKey != RangeKey("d") || ranges[1].endKey != RangeKey("m") {
		t.Errorf("unexpected second range: %+v", ranges[1])
	}

	if ranges[2].startKey != RangeKey("m") || ranges[2].endKey != RangeMaxKey {
		t.Errorf("unexpected terminal range: %+v", ranges[2])
	}
}

func TestArrangePartitionsAssignsMembersRoundRobin(t *testing.T) {
	ranges := []Range{
		newRange(RangeMinKey, RangeKey("h")),
		newRange(RangeKey("h"), RangeKey("t")),
		newRange(RangeKey("t"), RangeMaxKey),
	}
	members := []Member{
		{address: "node-a"},
		{address: "node-b"},
	}

	table := arrangePartitions(ranges, members)

	if len(table.partitions) != len(ranges) {
		t.Fatalf("expected %d partitions, got %d", len(ranges), len(table.partitions))
	}

	expectedAssignments := map[RangeKey]string{
		ranges[0].startKey: "node-a",
		ranges[1].startKey: "node-b",
		ranges[2].startKey: "node-a",
	}

	seenStarts := make(map[RangeKey]struct{})

	for _, info := range table.partitions {
		if info.id == "" {
			t.Error("partition id should not be empty")
		}
		if info.status != PartitionStatusAssigned {
			t.Errorf("expected status Assigned, got %s", info.status)
		}

		expectedAddress, ok := expectedAssignments[info.keyRange.startKey]
		if !ok {
			t.Fatalf("unexpected range start key %q", info.keyRange.startKey)
		}
		if info.memberAddress != expectedAddress {
			t.Errorf("range starting at %q assigned to %q, want %q", info.keyRange.startKey, info.memberAddress, expectedAddress)
		}
		seenStarts[info.keyRange.startKey] = struct{}{}
	}

	if len(seenStarts) != len(expectedAssignments) {
		t.Fatalf("expected %d unique ranges, saw %d", len(expectedAssignments), len(seenStarts))
	}
}

func TestArrangePartitionsWithNoMembers(t *testing.T) {
	ranges := createRangesFromSplitPoints([]string{"k"})

	table := arrangePartitions(ranges, nil)

	if len(table.partitions) != 0 {
		t.Fatalf("expected no partitions when no members are live, got %d", len(table.partitions))
	}
}
