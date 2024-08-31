package dep

import (
	"os"
	"os/signal"
	"syscall"

	queue "github.com/online-bnsp/backend/util/queue"
	"github.com/spf13/viper"
)

type MBI struct {
	q         queue.Queuer
	config    nsqConsumerConfig
	consumers []queue.Consumer
}

// nsqConsumerConfig consumer config
type nsqConsumerConfig struct {
	prefix, nsqLookUpds               string
	maxInFlight, maxAttempts, workers int
}

func InitMBI(configFile string) (*MBI, error) {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	q, err := queue.NewQueue("nsq", viper.GetString("nsqd"))
	if err != nil {
		return nil, err
	}

	consumer := &MBI{
		q: q,
		config: nsqConsumerConfig{
			maxInFlight: viper.GetInt("nsq_max_inflight"),
			maxAttempts: viper.GetInt("nsq_delay_time"),
			nsqLookUpds: viper.GetString("nsqlookupd"),
			workers:     viper.GetInt("nsq_workers"),
		},
		consumers: make([]queue.Consumer, 0, 10),
	}

	return consumer, nil
}

// register the consumer
func (mbi *MBI) Register(name, topic, channel string, handlerFn queue.ConsumerPayloadHandlerFn) {
	consumer, err := mbi.q.NewConsumer(queue.NsqConsumerArgs{
		Name:        name,
		Topic:       topic,
		Channel:     channel,
		Prefix:      mbi.config.prefix,
		MaxInFlight: mbi.config.maxInFlight,
		MaxAttempts: mbi.config.maxAttempts,
		NsqLookUpds: mbi.config.nsqLookUpds,
		Workers:     mbi.config.workers,
		HandlerFn:   handlerFn,
	})
	if err != nil {
		return
	}

	mbi.consumers = append(mbi.consumers, consumer)
}

// run the consumers
func (mbi *MBI) Run() {
	for _, c := range mbi.consumers {
		err := c.Run()
		if err != nil {
			return
		}
	}
}

// wait signal terminate
func (mbi *MBI) Wait() {
	// channel to listen for SIGINT (Ctrl+C) to signal to our application to gracefully shutdown all consumers.
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		for _, c := range mbi.consumers {
			c.Stop()
		}
	}
}
