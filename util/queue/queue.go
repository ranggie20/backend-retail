package queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/online-bnsp/backend/constant"

	nsq "github.com/nsqio/go-nsq"
)

// Queuer generic method to handle queue
type Queuer interface {
	NewProducer(v interface{}) (Producer, error)
	NewConsumer(v interface{}) (Consumer, error)
	setup(address string)
}

// Producer publish message to queue
type Producer interface {
	Publish(topic string, data interface{}) error
	DeferredPublish(topic string, data interface{}, delaySecond int) error
	DeferredPublishWithoutPrefix(topic string, data interface{}, delaySecond int) error
	Ping() error
}

// Consumer consumer message from queue
type Consumer interface {
	Run() error
	Stop()
}

// New create new queue service
func NewQueue(service, address string) (Queuer, error) {
	q, err := getQueueService(service)
	if err != nil {
		return nil, err
	}

	q.setup(address)

	return q, nil
}

func getQueueService(name string) (Queuer, error) {
	if name == "" {
		return nil, errors.New("queue service is required, please define which queue service to use")
	}

	switch name {
	case "nsq":
		return &Nsq{}, nil
	case "kafka":
		return nil, nil
	}

	return nil, errors.New("queue service is not available")
}

// Nsq is an object use to create its producer or consumer
type Nsq struct {
	Address string
}

// NsqProducerArgs option for nsq producer
type NsqProducerArgs struct {
	Prefix      string
	DelaySecond int
}

// NsqConsumerArgs arguments for nsq consumer
type NsqConsumerArgs struct {
	Name, Topic, Channel, Prefix, NsqLookUpds string
	MaxInFlight, MaxAttempts, Workers         int
	HandlerFn                                 ConsumerPayloadHandlerFn
}

// ConsumerPayloadHandlerFn is the signature for all payload handlers for NSQ consumers
type ConsumerPayloadHandlerFn func(context.Context, *nsq.Message) error

// MessageHandler consumer message handler
type MessageHandler struct {
	PayloadHandlerFn ConsumerPayloadHandlerFn
}

// nsqProducer contain producer connection and config
type nsqProducer struct {
	producer    *nsq.Producer
	prefix      string
	delaySecond int
}

// nsqConsumer contain list of consumers
type nsqConsumer struct {
	consumer    *nsq.Consumer
	workers     int
	handlerFn   ConsumerPayloadHandlerFn
	nsqLookUpds string
}

// NewProducer cretae new producer
func (q *Nsq) NewProducer(v interface{}) (Producer, error) {
	opt, ok := v.(NsqProducerArgs)
	if !ok {
		return nil, errors.New("invalid producer opt config")
	}

	p, err := nsq.NewProducer(q.Address, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	producer := &nsqProducer{
		producer:    p,
		prefix:      opt.Prefix,
		delaySecond: opt.DelaySecond,
	}

	return producer, nil
}

// Publish message to queue
func (p *nsqProducer) Publish(topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	topicWithPrefix := p.prefix + topic

	err = p.producer.Publish(topicWithPrefix, payload)
	if err != nil {
		return err
	}

	return nil
}

// DeferredPublish message to queue
func (p *nsqProducer) DeferredPublish(topic string, data interface{}, delaySecond int) error {
	topicWithPrefix := p.prefix + topic
	p.deferredPublish(topicWithPrefix, data, delaySecond)

	return nil
}

// DeferredPublish message to queue without prefix
func (p *nsqProducer) DeferredPublishWithoutPrefix(topic string, data interface{}, delaySecond int) error {
	p.deferredPublish(topic, data, delaySecond)

	return nil
}

// DeferredPublish message to queue
func (p *nsqProducer) deferredPublish(topic string, data interface{}, delaySecond int) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	delay := time.Second * time.Duration(p.delaySecond)
	if delaySecond != 0 {
		delay = time.Second * time.Duration(delaySecond)
	}

	err = p.producer.DeferredPublish(topic, delay, payload)
	if err != nil {
		return err
	}

	return nil
}

// Ping nsq producer
func (p *nsqProducer) Ping() error {
	return p.producer.Ping()
}

// NewConsumer create new consumer
func (q *Nsq) NewConsumer(v interface{}) (Consumer, error) {
	args, ok := v.(NsqConsumerArgs)
	if !ok {
		return nil, errors.New("invalid consumer opt config")
	}

	nsqConsumerConfig := nsq.NewConfig()
	nsqConsumerConfig.MaxInFlight = args.MaxInFlight
	nsqConsumerConfig.MaxAttempts = uint16(args.MaxAttempts)

	// Create a NewConsumer with the name of our topic, the channel, and our config
	c, err := nsq.NewConsumer(args.Prefix+args.Topic, args.Channel, nsqConsumerConfig)
	if err != nil {
		return nil, err
	}

	consumer := &nsqConsumer{
		consumer:    c,
		workers:     args.Workers,
		handlerFn:   args.HandlerFn,
		nsqLookUpds: args.NsqLookUpds,
	}

	return consumer, nil
}

// Run assign handlers for consumer and run it
func (c *nsqConsumer) Run() error {
	c.consumer.AddConcurrentHandlers(
		&MessageHandler{
			PayloadHandlerFn: c.handlerFn,
		},
		c.workers,
	)

	if err := c.consumer.ConnectToNSQLookupds([]string{c.nsqLookUpds}); err != nil {
		return err
	}

	return nil
}

// Stop the consumer
func (c *nsqConsumer) Stop() {
	c.consumer.Stop()
}

// HandleMessage is the implementation of nsq.Handler interface
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()

	t := time.Now()

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.ContextMessageID, string(m.ID[:]))
	ctx = context.WithValue(ctx, constant.ContextBirthTime, t)

	err := h.PayloadHandlerFn(ctx, m)
	if err != nil {
		m.Requeue(2 * time.Second)
		return err
	}

	m.Finish()

	return nil
}

// Setup set queue address
func (q *Nsq) setup(address string) {
	q.Address = address
}
