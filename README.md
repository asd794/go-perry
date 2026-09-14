# go-perry

A collection of Golang experiments, and topic-based examples.
Each directory generally represents a technical topic, implementation exercise, or experiment focused on a specific Go mechanism.The project is intended not only for learning Go, but also for exploring topics such as concurrency, networking, runtime internals, performance, systems programming, and other related areas.

## Topics

| Directory | Description |
| --- | --- |
| `alogrithm` | Algorithms and data structures, including design, search, sort, and structure topics. |
| `async` | Asynchronous processing examples, including producers, consumers, and background tasks. |
| `assembly` | Experiments involving assembly, CPU registers, instructions, and memory models. |
| `background_faktory` | Background job queue examples using Faktory. |
| `background_gocraft` | Background job queue examples using gocraft/work. |
| `barrier` | Goroutine synchronization using barriers. |
| `benchmark` | Go benchmarks, CPU profiling, and memory profiling experiments. |
| `bce` | Bound Check Elimination experiments. |
| `boomfilter` | Bloom Filter implementation examples. |
| `breakcircut` | Circuit Breaker implementation examples. |
| `build_tag` | Go build tags and conditional compilation examples. |
| `cache` | Cache implementations and cached HTTP client examples. |
| `channel` | Basic Go channel examples. |
| `channel_buffered` | Buffered channel examples. |
| `channel_unbuffered` | Unbuffered channel examples. |
| `context` | Go context examples, including cancellation, timeouts, and value propagation. |
| `convert` | Type conversion and data format conversion examples. |
| `decorator` | Decorator design pattern examples. |
| `dependency` | Dependency injection and dependency design examples. |
| `dump` | Program dump generation and debugging experiments. |
| `embed-demo` | Examples using Go's `embed` feature. |
| `errgroup` | Managing goroutines and collecting errors with `errgroup`. |
| `escape` | Escape analysis and memory allocation experiments. |
| `falseshare` | False sharing and performance experiments. |
| `generator` | Generator patterns and goroutine-based generator examples. |
| `gocron` | Scheduled job examples using gocron. |
| `gorilla_websocket` | WebSocket server and client examples using gorilla/websocket. |
| `gorm` | GORM examples, including models, migrations, and database operations. |
| `goroutine_confinement` | Goroutine confinement patterns for controlling access to shared data. |
| `goroutine_mutex` | Goroutine synchronization using mutexes. |
| `graphql` | GraphQL server examples. |
| `grpc` | gRPC client/server and Protocol Buffers examples. |
| `ipfs` | IPFS and decentralized storage experiments. |
| `iter` | Iterator patterns and iterator-style programming. |
| `jwt` | JWT authentication and authorization examples. |
| `kafka` | Kafka producer, consumer, and worker examples. |
| `linkname-demo` | Experimental examples using `//go:linkname`. |
| `logger` | Logging and file output examples. |
| `lsm` | LSM Tree storage engine experiments. |
| `memcache` | Memcache usage examples. |
| `module` | Go module usage and module organization. |
| `multipartload` | Multipart file upload client/server examples. |
| `nosplit-demo` | Experiments with `//go:nosplit`. |
| `orDone&T` | Go pipeline patterns such as `orDone` and `tee`. |
| `p2p` | Peer-to-peer networking experiments. |
| `paseto` | PASETO authentication and security examples. |
| `pin_thread` | Goroutine and OS thread binding experiments. |
| `pipeline` | Pipeline-based concurrent data processing. |
| `pool` | Resource pool and worker pool examples. |
| `protobuf` | Protocol Buffers serialization and code generation examples. |
| `rabbitmq` | RabbitMQ producer and consumer examples. |
| `ratelimit` | Rate limiting implementations and examples. |
| `redis_lock` | Redis distributed lock examples. |
| `reflect` | Go reflection experiments. |
| `rss` | RSS generation and parsing examples. |
| `runtime` | Go runtime concepts and runtime behavior experiments. |
| `schnorr` | Schnorr signature and cryptography-related examples. |
| `singleflight` | Request deduplication using `singleflight`. |
| `sort` | Sorting algorithms and sorting examples. |
| `sync_map` | `sync.Map` examples. |
| `test` | Go testing examples and experiments. |
| `unsafe` | `unsafe` package and low-level memory experiments. |
| `worker_pool` | Worker pool pattern examples. |
| `workerpool` | Additional worker pool implementations and experiments. |

The topic list is continuously evolving. New topics are welcome.

## Contributing

Contributions are welcome.

You are not limited to improving existing examples. You are encouraged to create your own concept, experiment, or technical topic.

The recommended workflow is to fork then create a feature branch:

```bash
git checkout -b feature/<your-concept>
```

For example:

```bash
git checkout -b feature/tracing
```

After implementing and documenting your concept, open a Pull Request.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the contribution workflow and guidelines.

## Authors

Contributors can add their name and ID to [AUTHORS.md](AUTHORS.md).

The author list is intended to give credit to everyone who contributes to the project.

## Philosophy

The goal of `go-perry` is not simply to collect finished code.

It is a place to:

- Experiment
- Build
- Document
- Share

A small idea can become a new topic.

**Create your own concept. Implement it. Document it. Share it.**