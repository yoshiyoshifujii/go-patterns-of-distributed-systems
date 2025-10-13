# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language Settings
- Default response language: Japanese
- Always respond in Japanese when I write to you in Japanese
- Technical terms and code examples can remain in English when necessary

## Project Overview

This repository implements distributed systems patterns from the book "分散システムのためのデザインパターン" (Patterns of Distributed Systems) by Maruzen Publishing in Go. The project is currently in its initial setup phase.

## Project Initialization

Since this is a new project, you'll need to initialize the Go module first:

```bash
go mod init github.com/yoshiyoshifujii/go-patterns-of-distributed-systems
```

## Expected Distributed Systems Patterns

The following patterns from the book should be implemented:

- **Consensus Algorithms**: Raft, Paxos, Multi-Paxos
- **Replication**: Leader-Follower, Chain Replication
- **Time & Ordering**: Lamport Timestamps, Vector Clocks, Hybrid Logical Clocks
- **Storage**: Write-Ahead Log (WAL), Segmented Log, Log-Structured Storage
- **Communication**: Request-Response, Gossip Protocol, Heartbeat
- **Failure Detection**: Phi Accrual Failure Detector, SWIM Protocol
- **Coordination**: Two-Phase Commit, Saga Pattern
- **Configuration Management**: Consistent Core, Lease-based Leadership

## Architecture Guidelines

When implementing patterns:

1. Each pattern should be in its own package under `internal/patterns/`
2. Provide example usage in `cmd/` for each major pattern
3. Use interfaces to allow different implementations of the same pattern
4. Include comprehensive tests that demonstrate failure scenarios
5. Document the trade-offs and use cases for each pattern

## Code Organization

Suggested structure for new implementations:

```
internal/
  patterns/
    consensus/
      raft/         # Raft consensus implementation
      paxos/        # Paxos implementation
    replication/
      leader/       # Leader-follower replication
    clock/
      lamport/      # Lamport logical clock
      vector/       # Vector clock
    storage/
      wal/          # Write-ahead log
```

## Testing Distributed Systems

When testing distributed patterns:
- Use table-driven tests for different scenarios
- Simulate network partitions and failures
- Test leader election scenarios
- Verify linearizability where applicable
- Use property-based testing for invariants

## Development Commands

Once the project is set up with a Makefile, common commands will be:

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run specific pattern tests
go test ./internal/patterns/consensus/raft/...

# Build examples
go build ./cmd/...
```

## Important Implementation Considerations

1. **Concurrency**: Use channels and goroutines appropriately for message passing between nodes
2. **Network Simulation**: Consider implementing a network simulator for testing distributed scenarios locally
3. **Time Handling**: Abstract time for testing (use interfaces for clocks)
4. **State Machines**: Many patterns involve state machines - keep state transitions explicit and testable
5. **Persistence**: Abstract storage interfaces to allow both in-memory and persistent implementations
