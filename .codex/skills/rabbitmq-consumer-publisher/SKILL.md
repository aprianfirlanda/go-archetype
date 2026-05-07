---
name: rabbitmq-consumer-publisher
description: Use when adding, changing, or debugging asynchronous RabbitMQ flows in the go-archetype repo, including message publisher ports, consumer ports, RabbitMQ adapter code, message handlers, topic registration, payload mapping, context propagation, and command wiring.
---

# RabbitMQ Consumer Publisher

Use this skill for async messaging work.

## Files And Packages

- RabbitMQ adapter root: `internal/adapters/messaging/rabbitmq`, package `messagingrmq`.
- RabbitMQ handlers: `internal/adapters/messaging/rabbitmq/handler/<domain>`.
- Task message handler example package: `taskhandlermq`.
- Publisher port: `internal/ports/output/message_publisher.go`.
- Consumer port: `internal/ports/input/message_consumer.go`.
- Consumer registry: `internal/bootstrap/consumer_registry.go`.
- Consumer command wiring: `cmd/consumer.go`.
- RabbitMQ config: `internal/infrastructure/config/schema.go`, `cmd/root.go`, and YAML config.

## Publisher Flow

- Application services should depend on `portout.MessagePublisher`, not RabbitMQ directly.
- Publish through the output port with `context.Context`.
- Derive logger from context in publisher methods with `logging.ComponentLogger(logging.FromContext(ctx), "<component>")`.
- Set AMQP `CorrelationId` and `MessageId` from context `rid` to preserve trace continuity.
- Use JSON payloads unless an existing flow uses another format.
- Keep topic names stable and centralized near wiring or domain-specific flow code.

## Consumer Flow

1. Add or update a message handler under `internal/adapters/messaging/rabbitmq/handler/<domain>`.
2. Keep handlers thin: unmarshal payload, map to application command/query, call input port, return error.
3. Register the topic and handler in `cmd/consumer.go` via `bootstrap.NewConsumerRegistry()`.
4. Let errors return to the consumer so it can run retry/backoff and DLQ routing; successful handling should `Ack`.
5. Preserve context propagation.
6. For each consumed message, create/propagate context `rid` from `CorrelationId` (fallback `MessageId`) before calling handlers/services.

## Retry Backoff And DLQ

- Queue names:
  - Main queue: `<topic>`
  - Retry queue: `<topic>.retry`
  - DLQ: `<topic>.dlq`
- Retry queue must dead-letter to main queue (`x-dead-letter-exchange=""`, `x-dead-letter-routing-key=<topic>`).
- Track retry state with header `x-retry-count`.
- Track failure details with:
  - `x-last-error` during retry scheduling
  - `x-final-error` and `x-original-queue` when moving to DLQ
- Use configured backoff list from `cfg.Messaging.RabbitMQ.Consumer.Retry.Backoff`.
- Behavior:
  - If retry slot exists: publish to retry queue with expiration delay, then `Ack` original.
  - If retry slots exhausted: publish to DLQ, then `Ack` original.
  - If publish retry/DLQ fails: `Nack(false, true)` original.
- Preserve message metadata when republishing (`CorrelationId`, `MessageId`, headers, body, content metadata).

## Boundaries

- Do not import RabbitMQ adapter code into domain or application packages.
- Do not put business rules in RabbitMQ handlers.
- Do not create background `context.Background()` inside request or consumer processing unless starting a true independent process.

## Verification

- Run `gofmt` on changed Go files.
- Run focused application and handler tests if present.
- Run `go test ./...` when dependencies are available.
