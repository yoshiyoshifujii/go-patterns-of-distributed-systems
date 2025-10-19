package chapter_20

var (
	partitionTableProvider = getPartitionTable
	partitionValuesFetcher = sendGetRangeMessage
)

func getValuesInRange(rng Range) []string {
	partitionTable := partitionTableProvider()
	partitionsInRange := partitionTable.getPartitionsInRange(rng)
	var values []string
	for _, partitionInfo := range partitionsInRange {
		partitionValues := partitionValuesFetcher(partitionInfo.getPartitionID(), rng, partitionInfo.getAddress())
		values = append(values, partitionValues...)
	}
	return values
}

func getPartitionTable() PartitionTable {
	panic("implement me")
}

func sendGetRangeMessage(partitionID string, rng Range, address string) []string {
	panic("implement me")
}
