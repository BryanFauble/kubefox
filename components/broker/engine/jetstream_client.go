package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/xigxog/kubefox/components/broker/config"
	"github.com/xigxog/kubefox/libs/core/kubefox"
	"github.com/xigxog/kubefox/libs/core/logkf"
	"google.golang.org/protobuf/proto"
)

// Content types
const (
	name       = "jetstream-client"
	evtStream  = "events"
	compBucket = "components"
)

var (
	maxMsgSize     = int32(1024 * 1024 * 5) // 5 MiB
	EventStreamTTL = time.Hour * 24 * 3     // 3 days
	ComponentsTTL  = time.Hour * 12         // 12 hours
)

type RecvMsg func(*nats.Msg)

type JetStreamClient struct {
	nats.JetStreamContext

	nc       *nats.Conn
	eventsKV nats.KeyValue
	compKV   nats.KeyValue

	brk Broker

	ctx    context.Context
	cancel context.CancelFunc

	mutex sync.Mutex
	log   *logkf.Logger
}

func NewJetStreamClient(brk Broker) *JetStreamClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &JetStreamClient{
		log:    logkf.Global,
		brk:    brk,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *JetStreamClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.nc != nil && c.nc.IsConnected() {
		c.log.Debug("jetstream client already connected")
		return nil
	}

	c.log.Debug("jetstream client connecting")

	var err error

	name, _ := os.LookupEnv("POD_NAME")
	if name == "" {
		name, _ = os.Hostname()
	}
	if name == "" {
		name = "unknown"
	}

	c.nc, err = nats.Connect(
		fmt.Sprintf("nats://%s", config.NATSAddr),
		nats.Name("broker-"+name),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(3),
		nats.PingInterval(time.Second*30),
		// c.natsTLS(c.cfg.Namespace),
		nats.RootCAs(kubefox.PathCACert),
		nats.ClientCert(kubefox.PathTLSCert, kubefox.PathTLSKey),
	)
	if err != nil {
		return c.log.ErrorN("connecting to NATS failed: %v", err)
	}
	if c.JetStreamContext, err = c.nc.JetStream(); err != nil {
		return c.log.ErrorN("connecting to JetStream failed: %v", err)
	}

	if err := c.setupStream(); err != nil {
		return err
	}
	if err := c.setupCompsKV(); err != nil {
		return err
	}
	if err := c.setupEventsKV(); err != nil {
		return err
	}

	c.log.Info("jetstream client connected")
	return nil
}

func (c *JetStreamClient) setupStream() error {
	if _, err := c.StreamInfo(evtStream); err != nil {
		_, err = c.AddStream(&nats.StreamConfig{
			Name: evtStream,
			// cannot be updated
			Storage:   nats.FileStorage,
			Retention: nats.LimitsPolicy,
			//
			NoAck: true,
		})
		if err != nil {
			return c.log.ErrorN("unable to create events stream: %v", err)
		}
	}
	_, err := c.UpdateStream(&nats.StreamConfig{
		Name:        evtStream,
		Description: "Durable disk backed event stream. Events are retained for 3 days.",
		Subjects:    []string{"evt.>"},
		MaxMsgSize:  maxMsgSize,
		Discard:     nats.DiscardOld,
		MaxAge:      EventStreamTTL,
		NoAck:       true,
	})
	if err != nil {
		return c.log.ErrorN("unable to create events stream: %v", err)
	}

	return nil
}

func (c *JetStreamClient) setupCompsKV() (err error) {
	c.compKV, err = c.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      compBucket,
		Description: "Durable key/value store used by Brokers to register Components. Values are retained for 12 hours.",
		Storage:     nats.FileStorage,
		TTL:         ComponentsTTL,
	})
	if err != nil {
		return c.log.ErrorN("unable to create components key/value store: %v", err)
	}

	return nil
}

func (c *JetStreamClient) setupEventsKV() (err error) {
	c.eventsKV, err = c.CreateKeyValue(&nats.KeyValueConfig{
		Bucket:      evtStream,
		Description: "Durable disk backed key/value store for Event ids. Values are retained for 3 days.",
		Storage:     nats.FileStorage,
		TTL:         EventStreamTTL,
	})
	if err != nil {
		return c.log.ErrorN("unable to create archive key/value store: %w", err)
	}

	return nil
}

// func (c *Client) natsTLS(namespace string) nats.Option {
// 	return func(o *nats.Options) error {
// var err error
// var pool *x509.CertPool
// var cert tls.Certificate

// caFile := path.Join(platform.NATSTLSDir, platform.CACertFile)
// certFile := path.Join(platform.NATSTLSDir, platform.TLSCertFile)
// keyFile := path.Join(platform.NATSTLSDir, platform.TLSKeyFile)

// if pem, err := os.ReadFile(c.cfg.CACertPath); err == nil {
// 	c.log.Debugf("reading tls certs from file")
// 	pool = x509.NewCertPool()
// 	ok := pool.AppendCertsFromPEM(pem)
// 	if !ok {
// 		return fmt.Errorf("failed to parse root certificate from %s", c.cfg.CACertPath)
// 	}

// 	// cert, err = tls.LoadX509KeyPair(certFile, keyFile)
// 	// if err != nil {
// 	// 	return fmt.Errorf("nats: error loading client certificate: %v", err)
// 	// }

// } else {
// 	c.log.Debugf("reading tls certs from kubernetes secret")
// 	pool, err = utils.GetCAFromSecret(ktyps.NamespacedName{
// 		Namespace: c.cfg.Namespace,
// 		Name:      fmt.Sprintf("%s-%s", c.cfg.Platform, kfp.RootCASecret),
// 	})
// 	if err != nil {
// 		return err
// 	}
// }

// cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
// if err != nil {
// 	return err
// }

// if o.TLSConfig == nil {
// 	o.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
// }
// o.TLSConfig.Certificates = []tls.Certificate{cert}
// 		o.TLSConfig.RootCAs = pool
// 		o.Secure = true

// 		return nil
// 	}
// }

func (c *JetStreamClient) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.log.Info("jetstream client closing")

	if c.cancel != nil {
		c.cancel()
	}

	if c.nc != nil {
		c.nc.Close()
	}
}

func (c *JetStreamClient) EventsKV() nats.KeyValue {
	return c.eventsKV
}

func (c *JetStreamClient) ComponentsKV() nats.KeyValue {
	return c.compKV
}

func (c *JetStreamClient) Publish(subject string, evt *kubefox.Event) error {
	dataBytes, err := proto.Marshal(evt)
	if err != nil {
		return err
	}

	h := make(nats.Header)
	h.Set(nats.MsgIdHdr, evt.Id)
	h.Set("ce_specversion", "1.0")
	h.Set("ce_type", evt.Type)
	h.Set("ce_id", evt.Id)
	h.Set("ce_time", time.Now().Format(time.RFC3339))
	h.Set("ce_source", fmt.Sprintf("kubefox:component:%s", evt.Source.GetId()))
	h.Set("ce_dataschema", kubefox.DataSchemaKubefox)
	h.Set("ce_datacontenttype", kubefox.ContentTypeProtobuf)

	// Note use of NATS directly instead of JetStream. This is done for
	// performance and memory efficiency. The risk is a msg not getting
	// delivered as there is no ACK from the server.
	return c.nc.PublishMsg(&nats.Msg{
		Subject: subject,
		Header:  h,
		Data:    dataBytes,
	})
}

func (c *JetStreamClient) PullEvents(sub ReplicaSubscription) error {
	log := c.log.WithComponent(sub.Component())

	jsSub, err := c.PullSubscribe(
		sub.Component().Subject(),
		sub.Component().Key(),
		nats.DeliverNew(),
		nats.AckNone(),
		nats.InactiveThreshold(time.Minute),
		nats.Context(sub.Context()),
	)
	if err != nil {
		return log.ErrorN("unable to create JetStream pull subscription: %v", err)
	}

	var grpJSSub *nats.Subscription
	if sub.IsGroupEnabled() {
		grpConsumer := sub.Component().GroupKey()
		grpSubj := sub.Component().GroupSubject()
		grpCfg := &nats.ConsumerConfig{
			Name:              grpConsumer,
			Durable:           grpConsumer,
			AckPolicy:         nats.AckNonePolicy,
			FilterSubject:     grpSubj,
			InactiveThreshold: EventStreamTTL,
		}
		if _, err := c.AddConsumer(evtStream, grpCfg); err != nil {
			if _, err := c.UpdateConsumer(evtStream, grpCfg); err != nil {
				return log.ErrorN("unable to update JetStream consumer: %v", err)
			}

		}
		grpJSSub, err = c.PullSubscribe(
			grpSubj,
			grpConsumer,
			nats.Bind(evtStream, grpConsumer),
			nats.Context(sub.Context()),
			nats.AckNone(),
		)
		if err != nil {
			return log.ErrorN("unable to create JetStream pull subscription: %v", err)
		}
	}

	recvMsg := func(msg *nats.Msg) {
		evt := kubefox.NewEvent()
		if md, err := msg.Metadata(); err == nil { // success
			evt.ReduceTTL(md.Timestamp)
		}

		if err := proto.Unmarshal(msg.Data, evt); err != nil {
			log.Warn(err)
			return
		}

		rEvt := &ReceivedEvent{
			ActiveEvent:  kubefox.StartEvent(evt),
			Subscription: sub,
			Receiver:     ReceiverJetStream,
		}
		if err := c.brk.RecvEvent(rEvt); err != nil {
			c.log.WithEvent(evt).Debug(err)
			if evt.Target.Id == "" && errors.Is(err, ErrSubCanceled) {
				// Any component replica can process so do not remove msg from
				// subject and let another JetStream subscriber process it.
				// msg.Nak()
				// TODO, republish event since ACK disabled
				return
			}
		}
	}

	go c.pullEvents(log, jsSub, recvMsg)
	if grpJSSub != nil {
		go c.pullEvents(log, grpJSSub, recvMsg)
	}

	return nil
}

func (c *JetStreamClient) pullEvents(log *logkf.Logger, jsSub *nats.Subscription, recvMsg RecvMsg) {
	for {
		msgs, err := jsSub.Fetch(1)
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				// timeout waiting for msg, just go back to waiting
				continue
			}
			if !jsSub.IsValid() {
				log.Debug("jetstream pull subscription closed")
				return
			}
			if errors.Is(err, nats.ErrConnectionClosed) {
				// TODO deal with err
				log.Debug("jetstream connection closed")
				return
			}

			// TODO nats: Server Shutdown , should exit

			log.Error(err)
			// sub.errCount += 1
			// simple backoff, max 3 seconds
			// sleepTime := math.Min(3, float64(sub.errCount-1))
			// time.Sleep(time.Duration(sleepTime) * time.Second)
			time.Sleep(time.Second)
			continue
		}

		for _, msg := range msgs {
			recvMsg(msg)
		}
	}
}

func (c *JetStreamClient) IsHealthy(ctx context.Context) bool {
	return c.nc != nil && c.nc.IsConnected()
}

func (c *JetStreamClient) Name() string {
	return name
}
