package messagingrmq

import (
	"context"
	"encoding/json"
	"fmt"
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"
	"strconv"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	ch        *amqp091.Channel
	publishCh *amqp091.Channel
	publishMu sync.Mutex
	logger    *logrus.Entry
	retry     config.RabbitMQConsumerRetry
}

func NewConsumer(conn *Connection, logger *logrus.Entry, retry config.RabbitMQConsumerRetry) (*Consumer, error) {
	ch, err := conn.Conn.Channel()
	if err != nil {
		return nil, err
	}
	publishCh, err := conn.Conn.Channel()
	if err != nil {
		_ = ch.Close()
		return nil, err
	}
	return &Consumer{
		ch:        ch,
		publishCh: publishCh,
		logger:    logger,
		retry:     retry,
	}, nil
}

func (c *Consumer) Consume(
	ctx context.Context,
	topic string,
	handler portin.MessageHandler,
) error {

	if err := c.declareQueues(topic); err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		topic,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			msgCtx := c.contextWithMessageLogger(ctx, topic, msg)
			msgLogger := logging.FromContext(msgCtx)

			msgLogger.Debug("Message received")

			err := handler(msgCtx, msg.Body)
			if err != nil {
				msgLogger.WithError(err).Error("Handler failed")
				c.handleFailure(msgCtx, topic, msg, err)
				continue
			}
			_ = msg.Ack(false)
			msgLogger.Debug("Message handled")
		}
	}()

	return nil
}

func (c *Consumer) declareQueues(topic string) error {
	_, err := c.ch.QueueDeclare(
		topic,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare(
		retryQueueName(topic),
		true, // durable
		false,
		false,
		false,
		amqp091.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": topic,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.ch.QueueDeclare(
		dlqName(topic),
		true, // durable
		false,
		false,
		false,
		nil,
	)
	return err
}

func (c *Consumer) handleFailure(ctx context.Context, topic string, msg amqp091.Delivery, handlerErr error) {
	logger := logging.FromContext(ctx)
	retryCount := retryCountFromHeaders(msg.Headers)

	if retryCount >= len(c.retry.Backoff) {
		if err := c.publishDLQ(ctx, topic, msg, retryCount, handlerErr); err != nil {
			logger.WithError(err).Error("Failed to publish message to DLQ")
			_ = msg.Nack(false, true)
			return
		}
		_ = msg.Ack(false)
		logger.WithFields(logrus.Fields{
			"queue":       topic,
			"dlq":         dlqName(topic),
			"retry_count": retryCount,
		}).Warn("Message moved to DLQ")
		return
	}

	nextRetryCount := retryCount + 1
	delay := c.retryBackoff(nextRetryCount)
	if err := c.publishRetry(ctx, topic, msg, nextRetryCount, delay, handlerErr); err != nil {
		logger.WithError(err).Error("Failed to publish message to retry queue")
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
	logger.WithFields(logrus.Fields{
		"queue":       topic,
		"retry_queue": retryQueueName(topic),
		"retry_count": nextRetryCount,
		"delay":       delay.String(),
	}).Warn("Message scheduled for retry")
}

func (c *Consumer) publishRetry(ctx context.Context, topic string, msg amqp091.Delivery, retryCount int, delay time.Duration, handlerErr error) error {
	headers := copyHeaders(msg.Headers)
	headers["x-retry-count"] = retryCount
	headers["x-last-error"] = handlerErr.Error()

	publishing := publishingFromDelivery(msg, headers)
	publishing.Expiration = strconv.FormatInt(delayMillis(delay), 10)

	c.publishMu.Lock()
	defer c.publishMu.Unlock()

	return c.publishCh.PublishWithContext(
		ctx,
		"",
		retryQueueName(topic),
		false,
		false,
		publishing,
	)
}

func (c *Consumer) publishDLQ(ctx context.Context, topic string, msg amqp091.Delivery, retryCount int, handlerErr error) error {
	headers := copyHeaders(msg.Headers)
	headers["x-retry-count"] = retryCount
	headers["x-final-error"] = handlerErr.Error()
	headers["x-original-queue"] = topic

	c.publishMu.Lock()
	defer c.publishMu.Unlock()

	return c.publishCh.PublishWithContext(
		ctx,
		"",
		dlqName(topic),
		false,
		false,
		publishingFromDelivery(msg, headers),
	)
}

func (c *Consumer) retryBackoff(attempt int) time.Duration {
	return c.retry.Backoff[attempt-1]
}

func publishingFromDelivery(msg amqp091.Delivery, headers amqp091.Table) amqp091.Publishing {
	deliveryMode := msg.DeliveryMode
	if deliveryMode == 0 {
		deliveryMode = amqp091.Persistent
	}

	return amqp091.Publishing{
		Headers:         headers,
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    deliveryMode,
		Priority:        msg.Priority,
		CorrelationId:   msg.CorrelationId,
		ReplyTo:         msg.ReplyTo,
		Expiration:      msg.Expiration,
		MessageId:       msg.MessageId,
		Timestamp:       msg.Timestamp,
		Type:            msg.Type,
		UserId:          msg.UserId,
		AppId:           msg.AppId,
		Body:            msg.Body,
	}
}

func copyHeaders(headers amqp091.Table) amqp091.Table {
	copied := amqp091.Table{}
	for key, value := range headers {
		copied[key] = value
	}
	return copied
}

func retryCountFromHeaders(headers amqp091.Table) int {
	value, ok := headers["x-retry-count"]
	if !ok {
		return 0
	}

	switch count := value.(type) {
	case int:
		return count
	case int8:
		return int(count)
	case int16:
		return int(count)
	case int32:
		return int(count)
	case int64:
		return int(count)
	case uint:
		return int(count)
	case uint8:
		return int(count)
	case uint16:
		return int(count)
	case uint32:
		return int(count)
	case uint64:
		return int(count)
	case string:
		parsed, err := strconv.Atoi(count)
		if err == nil {
			return parsed
		}
	}

	return 0
}

func retryQueueName(topic string) string {
	return topic + ".retry"
}

func dlqName(topic string) string {
	return topic + ".dlq"
}

func delayMillis(delay time.Duration) int64 {
	millis := delay.Milliseconds()
	if millis < 1 {
		return 1
	}
	return millis
}

func (c *Consumer) contextWithMessageLogger(ctx context.Context, topic string, msg amqp091.Delivery) context.Context {
	var envelope struct {
		RID string `json:"rid"`
	}
	rid := msg.CorrelationId
	if rid == "" {
		rid = msg.MessageId
	}
	if rid == "" && json.Unmarshal(msg.Body, &envelope) == nil {
		rid = envelope.RID
	}
	if rid == "" {
		rid = fmt.Sprintf("msg-%d", msg.DeliveryTag)
	}
	ctx = logging.WithRequestID(ctx, rid)

	logger := c.logger
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	logger = logger.WithFields(logrus.Fields{
		"component":  "messaging.rabbitmq.consumer",
		"rid":        logging.RequestIDFromContext(ctx),
		"request_id": logging.RequestIDFromContext(ctx),
		"topic":      topic,
	})

	return logging.WithLogger(ctx, logger)
}
