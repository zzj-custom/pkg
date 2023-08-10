package pRedis

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Usage:
// Should call iRedis.NewSubscription to create a subscription, then call Subscribe method
// to listen a publishing channel with a callback function that needs an error as return value.
// After subscribe, `go sub.Start()` should be called to start a goroutine to process context
// callback, restart callback, ping, and so on.
//
// Example:
// sub, _ := iRedis.NewSubscription(opts)
// sub.Subscribe("policy", func(ctx context.Context, msg redis.Message) error {
//     log.Info("收到policy订阅消息", string(msg.Data))
//     return nil
// })
// go sub.Start()

type SubscriptionProcessor func(ctx context.Context, msg redis.Message) error

const (
	defaultRestartDuration = 3 * time.Minute
	defaultPingDuration    = 30 * time.Second
)

type Subscription struct {
	Context         context.Context
	RedisPool       *redis.Pool
	RestartDuration time.Duration
	PingDuration    time.Duration

	restartTicker *time.Ticker
	pingTicker    *time.Ticker

	channel   string
	processor SubscriptionProcessor
	psc       redis.PubSubConn
}

func (r *Subscription) Stop() {
	defer func() {
		_ = r.psc.Close()
	}()
	r.pingTicker.Stop()
	r.restartTicker.Stop()
	r.pingTicker = nil
	r.restartTicker = nil
}

func (r *Subscription) Subscribe(channel string, processor SubscriptionProcessor) error {
	r.channel = channel
	r.processor = processor
	if r.RestartDuration == 0 {
		r.RestartDuration = defaultRestartDuration
	}
	if r.PingDuration == 0 {
		r.PingDuration = defaultPingDuration
	}
	if r.channel == "" || r.processor == nil {
		return errors.Errorf("channel and processor must be not nil")
	}
	r.restartTicker = time.NewTicker(r.RestartDuration)
	r.pingTicker = time.NewTicker(r.PingDuration)
	return nil
}

func (r *Subscription) subscribe() error {
	conn := r.RedisPool.Get()
	r.psc = redis.PubSubConn{Conn: conn}
	defer func() {
		_ = r.psc.Close()
	}()

	if err := r.psc.Subscribe(r.channel); err != nil {
		return errors.Wrap(err, "subscribe failed")
	}

	for {
		switch n := r.psc.Receive().(type) {
		case redis.Message:
			r.processor(r.Context, n)
		case error:
			if strings.Contains(n.Error(), "redigo: connection closed") ||
				strings.Contains(n.Error(), "use of closed network connection") {
				log.WithError(n).Info("连接已断开，退出。")
				return nil
			}
		case redis.Subscription:
			if n.Kind == "unsubscribe" {
				return nil
			}
		default:
			//log.Infof("default, %+v", n)
		}
	}
}

func (r *Subscription) Start() {
	go func() {
		err := r.subscribe()
		if err != nil {
			log.Error("failed to subscribe: ", r.channel, ", error:", err)
		}
	}()
	for {
		select {
		case <-r.Context.Done():
			r.Stop()
			return
		case <-r.restartTicker.C:
			go func() {
				r.pingTicker.Stop()
				_ = r.psc.Unsubscribe()
				time.Sleep(1 * time.Second)
				r.pingTicker.Reset(r.PingDuration)
				_ = r.subscribe()
			}()
		case <-r.pingTicker.C:
			r.psc.Ping("PING")
		}
	}
}

type Options struct {
	Ctx             context.Context
	Pool            *redis.Pool
	RestartDuration time.Duration
	PingDuration    time.Duration
}

func NewSubscription(opts Options) (*Subscription, error) {
	if opts.Ctx == nil {
		return nil, errors.Errorf("context must be specified")
	}
	if opts.Pool == nil {
		return nil, errors.Errorf("redis pool must be specified")
	}
	sub := new(Subscription)
	sub.Context = opts.Ctx
	sub.RedisPool = opts.Pool
	if opts.RestartDuration > 0 {
		sub.RestartDuration = opts.RestartDuration
	}
	if opts.PingDuration > 0 {
		sub.PingDuration = opts.PingDuration
	}
	return sub, nil
}
