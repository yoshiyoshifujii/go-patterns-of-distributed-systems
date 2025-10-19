package chapter_20

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValuesInRangeAggregatesPartitionResponses(t *testing.T) {
	originalPartitionTableProvider := partitionTableProvider
	originalPartitionValuesFetcher := partitionValuesFetcher
	t.Cleanup(func() {
		partitionTableProvider = originalPartitionTableProvider
		partitionValuesFetcher = originalPartitionValuesFetcher
	})

	// Given
	partitionTable := PartitionTable{
		partitions: map[string]PartitionInfo{
			"partition-1": newPartitionInfo("partition-1", "node-a", PartitionStatusAssigned, newRange(RangeMinKey, RangeKey("m"))),
			"partition-2": newPartitionInfo("partition-2", "node-b", PartitionStatusAssigned, newRange(RangeKey("m"), RangeMaxKey)),
		},
	}

	partitionTableProvider = func() PartitionTable {
		return partitionTable
	}

	expectedResponses := map[string][]string{
		"partition-1": {"alpha", "beta"},
		"partition-2": {"gamma"},
	}

	type fetchCall struct {
		partitionID string
		rng         Range
		address     string
	}

	var calls []fetchCall

	partitionValuesFetcher = func(partitionID string, rng Range, address string) []string {
		calls = append(calls, fetchCall{
			partitionID: partitionID,
			rng:         rng,
			address:     address,
		})
		return expectedResponses[partitionID]
	}

	requestRange := newRange(RangeKey("a"), RangeKey("z"))

	// When
	actual := getValuesInRange(requestRange)

	// Then
	assert.ElementsMatch(t, []string{"alpha", "beta", "gamma"}, actual, "aggregated values mismatch")
	assert.Len(t, calls, len(expectedResponses), "fetcher call count mismatch")

	expectedAddresses := map[string]string{
		"partition-1": "node-a",
		"partition-2": "node-b",
	}

	for _, call := range calls {
		assert.Equal(t, requestRange, call.rng, "unexpected range passed to fetcher")

		expectedAddress, ok := expectedAddresses[call.partitionID]
		assert.Truef(t, ok, "unexpected partition id %q", call.partitionID)
		assert.Equalf(t, expectedAddress, call.address, "address mismatch for partition %q", call.partitionID)
	}
}

func TestGetValuesInRangeWhenNoPartitionsOverlap(t *testing.T) {
	originalPartitionTableProvider := partitionTableProvider
	originalPartitionValuesFetcher := partitionValuesFetcher
	t.Cleanup(func() {
		partitionTableProvider = originalPartitionTableProvider
		partitionValuesFetcher = originalPartitionValuesFetcher
	})

	// Given
	partitionTable := PartitionTable{
		partitions: map[string]PartitionInfo{
			"partition-99": newPartitionInfo("partition-99", "node-x", PartitionStatusAssigned, newRange(RangeKey("x"), RangeMaxKey)),
		},
	}

	partitionTableProvider = func() PartitionTable {
		return partitionTable
	}

	partitionValuesFetcher = func(partitionID string, rng Range, address string) []string {
		t.Fatalf("partition fetcher should not be called when no partitions overlap")
		return nil
	}

	requestRange := newRange(RangeKey("a"), RangeKey("b"))

	// When
	actual := getValuesInRange(requestRange)

	// Then
	assert.Empty(t, actual, "expected empty result when no partitions overlap request range")
}
