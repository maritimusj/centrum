package ep6v2

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"

	Contract "github.com/maritimusj/chuanyan/gate/adapter/contract"
	"github.com/maritimusj/chuanyan/gate/device_types/ep6v2"
	L "github.com/maritimusj/chuanyan/gate/lang"
)

const (
	DefaultGetDataFreq     = 3
	DefaultGetPerfDataFreq = 6
	checkTimeout           = time.Second * 3

	maxDispatchRoutines    = 100
	deviceDataWaitChanSize = 10
	dataTransferChanSize   = 2

	perfDataTagName = "延迟"
)

type stats struct {
	packetSend     int64
	packetReceived int64
	storesOk       int32
	storesFailed   int32
	status         string
}

type ep6v2Adapter struct {
	uid    string
	device *ep6v2.Device

	getDataFreq     int
	getPerfDataFreq int

	dataStream      chan *Contract.MeasureData
	dataChannelsMap sync.Map
	dataChannels    int32

	stateChangeFNs []func(string)

	done context.CancelFunc
	ctx  context.Context

	stats stats

	wg sync.WaitGroup
	mu sync.Mutex
}

func New() Contract.Client {
	return &ep6v2Adapter{
		stateChangeFNs:  []func(string){},
		getPerfDataFreq: DefaultGetPerfDataFreq,
		getDataFreq:     DefaultGetDataFreq,
	}
}

func (adapter *ep6v2Adapter) GetUID() string {
	return adapter.uid
}

func (adapter *ep6v2Adapter) SetUID(uid string) {
	adapter.uid = uid
}

func (adapter *ep6v2Adapter) Stats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"status":         adapter.stats.status,
		"packetSend":     atomic.LoadInt64(&adapter.stats.packetSend),
		"packetReceived": atomic.LoadInt64(&adapter.stats.packetReceived),
		"storesOk":       atomic.LoadInt32(&adapter.stats.storesOk),
		"storesFailed":   atomic.LoadInt32(&adapter.stats.storesFailed),
	}, nil
}

func (adapter *ep6v2Adapter) Open(ctx context.Context, option Contract.Option) error {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()

	adapter.ctx, adapter.done = context.WithCancel(ctx)

	if v, ok := option["uid"]; ok {
		adapter.SetUID(v.(string))
	}

	if v, ok := option["stateChange"]; ok {
		if fn, ok := v.(func(string)); ok {
			adapter.stateChangeFNs = append(adapter.stateChangeFNs, fn)
		}
	}

	if v, ok := option["getDataFreq"]; ok {
		if freq, ok := v.(int); ok && freq > 0 {
			adapter.getDataFreq = freq
		}
	}

	if v, ok := option["getPerfDataFreq"]; ok {
		if freq, ok := v.(int); ok && freq > 0 {
			adapter.getPerfDataFreq = freq
		}
	}

	if v, ok := option["connStr"]; ok {
		if connStr, ok := v.(string); ok {
			if govalidator.IsMAC(connStr) {
				adapter.device = ep6v2.Inverse(GetInverseServer(), strings.ToLower(connStr))
			} else {
				adapter.device = ep6v2.Connect(connStr, 6*time.Second)
			}

			adapter.device.SetUID(adapter.uid)
			return nil
		}
	}

	return L.Error(L.ErrInvalidConnStr)
}

func (adapter *ep6v2Adapter) Status(status string) {
	if adapter.stateChangeFNs != nil {
		for _, fn := range adapter.stateChangeFNs {
			fn(status)
		}
	}

	adapter.stats.status = status
}

func (adapter *ep6v2Adapter) Close() error {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()

	log.WithFields(log.Fields{
		"src": adapter.uid,
	}).Trace("adapter closing...")

	adapter.Status(L.Str(L.StatusClosing))

	err := adapter.device.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"src": adapter.uid,
		}).Error("adapter's device client: ", err)
	}

	if adapter.done != nil {
		adapter.done()
	}

	adapter.wg.Wait()

	if adapter.dataStream != nil {
		close(adapter.dataStream)
		adapter.dataStream = nil
	}

	adapter.Status(L.Str(L.StatusClosed))

	log.WithFields(log.Fields{
		"src": adapter.uid,
	}).Trace("adapter closed.")

	return nil
}

//创建数据接收通道
func (adapter *ep6v2Adapter) Create() (Contract.MeasureDataChannel, error) {
	adapter.init()

	ch := make(chan *Contract.MeasureData, dataTransferChanSize)
	adapter.dataChannelsMap.Store(ch, new(int32))

	atomic.AddInt32(&adapter.dataChannels, 1)
	atomic.AddInt32(&adapter.stats.storesOk, 1)

	return ch, nil
}

//丢弃数据接收通道
func (adapter *ep6v2Adapter) Drop(channel Contract.MeasureDataChannel) {
	go func() {
		adapter.dataChannelsMap.Range(func(key, _ interface{}) bool {
			v, _ := key.(chan *Contract.MeasureData)
			if v == channel {
				adapter.dataChannelsMap.Delete(key)
				atomic.AddInt32(&adapter.dataChannels, -1)
				atomic.AddInt32(&adapter.stats.storesOk, -1)

				//close(v)
				return false
			}
			return true
		})
	}()
}

//插入控制数据通道
func (adapter *ep6v2Adapter) Plug(channel <-chan *Contract.CtrlData) error {
	adapter.wg.Add(1)

	go func() {
		defer func() {
			recover()
			adapter.wg.Done()
		}()
		for {
			select {
			case data := <-channel:
				if data != nil {
					for ch, v := range data.Values {
						if data.Error != nil {
							data.Error <- adapter.device.SetChannelValue(ch, v)
						} else {
							err := adapter.device.SetChannelValue(ch, v)
							if err != nil {
								log.Trace(err)
							}
						}
					}
				}
			case <-adapter.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (adapter *ep6v2Adapter) Set(options Contract.Option) error {
	for key, v := range options {
		switch key {
		case "stateChange":
			if fn, ok := v.(func(string)); ok {
				adapter.stateChangeFNs = append(adapter.stateChangeFNs, fn)
			}
		case "getDataFreq":
			if freq, ok := v.(int); ok && freq > 0 {
				adapter.getDataFreq = freq
			}
		case "getPerfDataFreq":
			if freq, ok := v.(int); ok && freq > 0 {
				adapter.getPerfDataFreq = freq
			}
		default:
			return L.Error(L.ErrUnknownParams, key)
		}
	}
	return nil
}

func (adapter *ep6v2Adapter) Get(keys ...string) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for _, key := range keys {
		switch key {
		case "getDataFreq":
			result["getDataFreq"] = adapter.getDataFreq
		case "getPerfDataFreq":
			result["getPerfDataFreq"] = adapter.getPerfDataFreq
		default:
			return nil, L.Error(L.ErrUnknownParams, key)
		}
	}
	return result, nil
}

func (adapter *ep6v2Adapter) init() {
	adapter.mu.Lock()
	defer adapter.mu.Unlock()

	if adapter.dataStream == nil {

		adapter.Status(L.Str(L.StatusInitialize))
		adapter.dataStream = make(chan *Contract.MeasureData, deviceDataWaitChanSize)

		adapter.wg.Add(4)
		{
			go adapter.fetchDataFromDevice(adapter.ctx)
			go adapter.fetchPerfData(adapter.ctx)
			go adapter.dispatchData(adapter.ctx)
			go adapter.timeoutCheck(adapter.ctx)
		}
	}
}
func (adapter *ep6v2Adapter) timeoutCheck(ctx context.Context) {
	defer func() {
		log.WithFields(log.Fields{
			"src": adapter.uid,
		}).Trace("timeoutCheck thread exit.")

		adapter.wg.Done()
	}()

	for {
		select {
		case <-time.After(checkTimeout):
			adapter.dataChannelsMap.Range(func(key, routines interface{}) bool {
				if atomic.LoadInt32(routines.(*int32)) > maxDispatchRoutines {
					adapter.dataChannelsMap.Delete(key)
					close(key.(chan *Contract.MeasureData))

					atomic.AddInt32(&adapter.dataChannels, -1)
					atomic.AddInt32(&adapter.stats.storesFailed, 1)
				}
				return true
			})
		case <-ctx.Done():
			return
		}
	}
}

func (adapter *ep6v2Adapter) fetchPerfData(ctx context.Context) {
	defer func() {
		log.WithFields(log.Fields{
			"src": adapter.uid,
		}).Trace("fetchPerfData thread exit.")

		adapter.wg.Done()
	}()

	for {
		select {
		case <-time.After(time.Duration(adapter.getPerfDataFreq) * time.Second):

			log.WithFields(log.Fields{
				"src": adapter.uid,
			}).Trace("fetchPerfData...")

			if adapter.device.IsValid() {
				delay := adapter.device.Perf().Delay("")

				atomic.AddInt64(&adapter.stats.packetReceived, 1)

				adapter.dataStream <- &Contract.MeasureData{
					Name: perfDataTagName,
					Tags: map[string]string{
						"ip":  adapter.device.GetAddr().Ip.String(),
						"mac": adapter.device.GetAddr().Mac.String(),
					},
					Fields: map[string]interface{}{
						"value": delay,
					},
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (adapter *ep6v2Adapter) dispatchData(ctx context.Context) {
	defer func() {
		log.WithFields(log.Fields{
			"src": adapter.uid,
		}).Trace("dispatchData thread exit.")

		adapter.wg.Done()
	}()

	for {
		select {
		case v, ok := <-adapter.dataStream:
			if !ok {
				return
			}
			if v != nil {
				adapter.dataChannelsMap.Range(func(ch, routines interface{}) bool {
					adapter.wg.Add(1)

					go func() {
						defer func() {
							recover()
							adapter.wg.Done()
						}()

						atomic.AddInt32(routines.(*int32), 1)
						defer atomic.AddInt32(routines.(*int32), -1)
						defer atomic.AddInt64(&adapter.stats.packetSend, 1)

						ch.(chan *Contract.MeasureData) <- v
					}()

					return true
				})
			}
		case <-ctx.Done():
			return
		}
	}
}

func (adapter *ep6v2Adapter) fetchDataFromDevice(ctx context.Context) {
	defer func() {
		log.WithFields(log.Fields{
			"src": adapter.uid,
		}).Trace("fetchDataFromDevice thread exit.")

		adapter.wg.Done()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(adapter.getDataFreq) * time.Second):
			if adapter.device.IsConnected() {
				if atomic.LoadInt32(&adapter.dataChannels) > 0 {
					log.WithFields(log.Fields{
						"src": adapter.uid,
					}).Trace("getCHData...")

					adapter.Status(L.Str(L.StatusRefresh))

					adapter.getCHData()
				} else {
					adapter.Status(L.Str(L.StatusIdle))
					adapter.heartBeat()
				}
			} else {
				adapter.Status(L.Str(L.StatusClosed))
			}
		}
	}
}

func (adapter *ep6v2Adapter) heartBeat() {
	adapter.device.GetCHNum()
}

func (adapter *ep6v2Adapter) getCHData() {
	adapter.device.RefreshCHData()

	for i := 0; adapter.ctx.Err() == nil && i < adapter.device.GetCHNum().DO; i++ {
		do := adapter.device.GetDO(i)
		if do != nil && do.IsEnabled() {
			atomic.AddInt64(&adapter.stats.packetReceived, 1)

			v, _ := do.GetValue()
			adapter.dataStream <- &Contract.MeasureData{
				Name: do.GetTitle(),
				Tags: map[string]string{
					"model":   adapter.device.GetModel().ID,
					"version": adapter.device.GetModel().Version,
					"ip":      adapter.device.GetAddr().Ip.String(),
					"mac":     adapter.device.GetAddr().Mac.String(),
					"tag":     do.GetTagName(),
					"ctl":     "true",
				},
				Fields: map[string]interface{}{
					"value": v,
				},
			}
		}
	}

	for i := 0; adapter.ctx.Err() == nil && i < adapter.device.GetCHNum().DI; i++ {
		di := adapter.device.GetDI(i)
		if di != nil && di.IsEnabled() {
			atomic.AddInt64(&adapter.stats.packetReceived, 1)

			v, _ := di.GetValue()
			adapter.dataStream <- &Contract.MeasureData{
				Name: di.GetTitle(),
				Tags: map[string]string{
					"model":   adapter.device.GetModel().ID,
					"version": adapter.device.GetModel().Version,
					"ip":      adapter.device.GetAddr().Ip.String(),
					"mac":     adapter.device.GetAddr().Mac.String(),
					"tag":     di.GetTagName(),
					"ctl":     "false",
				},
				Fields: map[string]interface{}{
					"value": v,
				},
			}
		}
	}

	for i := 0; adapter.ctx.Err() == nil && i < adapter.device.GetCHNum().AI; i++ {
		ai := adapter.device.GetAI(i)
		if ai != nil && ai.IsEnabled() {
			atomic.AddInt64(&adapter.stats.packetReceived, 1)

			v, s := ai.GetValue()
			adapter.dataStream <- &Contract.MeasureData{
				Name: ai.GetTitle(),
				Tags: map[string]string{
					"model":   adapter.device.GetModel().ID,
					"version": adapter.device.GetModel().Version,
					"ip":      adapter.device.IpAddr(), //adapter.device.GetAddr().Ip.String(),
					"mac":     adapter.device.GetAddr().Mac.String(),
					"tag":     ai.GetTagName(),
					"unit":    ai.GetUint(),
					"alarm":   s,
				},
				Fields: map[string]interface{}{
					"value": v,
				},
			}
		}
	}

	for i := 0; adapter.ctx.Err() == nil && i < adapter.device.GetCHNum().AO; i++ {
		ao := adapter.device.GetAO(i)
		if ao != nil && ao.IsEnabled() {
			atomic.AddInt64(&adapter.stats.packetReceived, 1)

			v, s := ao.GetValue()
			adapter.dataStream <- &Contract.MeasureData{
				Name: ao.GetTitle(),
				Tags: map[string]string{
					"model":   adapter.device.GetModel().ID,
					"version": adapter.device.GetModel().Version,
					"ip":      adapter.device.GetAddr().Ip.String(),
					"mac":     adapter.device.GetAddr().Mac.String(),
					"unit":    ao.GetUint(),
					"tag":     ao.GetTagName(),
					"alarm":   s,
				},
				Fields: map[string]interface{}{
					"value": v,
				},
			}
		}
	}

}
