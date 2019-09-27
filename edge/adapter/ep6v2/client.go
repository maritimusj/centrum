package ep6v2

import (
	"context"
	"sync"
	"time"

	"github.com/maritimusj/centrum/edge"
	"github.com/maritimusj/centrum/edge/devices/ep6v2"
	"github.com/maritimusj/centrum/lang"
	"github.com/maritimusj/centrum/util"
	log "github.com/sirupsen/logrus"
)

type ep6v2Adapter struct {
	device *ep6v2.Device

	channels sync.Map

	cond *sync.Cond
	ctx  context.Context
}

func New() *ep6v2Adapter {
	return &ep6v2Adapter{
		cond:   sync.NewCond(&sync.Mutex{}),
		device: ep6v2.New(),
	}
}

func (e *ep6v2Adapter) Open(ctx context.Context, option map[string]interface{}) error {
	if connStr, ok := option["connStr"].(string); ok {
		err := e.device.Connect(ctx, connStr)
		if err != nil {
			return err
		}
		_, err = e.device.GetModel()
		if err != nil {
			return err
		}
		e.ctx = ctx
		return nil
	}
	return lang.Error(lang.ErrInvalidDeviceConnStr)
}

func (e *ep6v2Adapter) Close() error {
	if e.device != nil {
		e.device.Close()
	}
	return nil
}

func (e *ep6v2Adapter) fetchData(ch chan<- *edge.MeasureData) error {
	r, err := e.device.GetRealTimeData()
	if err != nil {
		return err
	}
	for i := 0; i < r.AINum(); i++ {
		ai, err := e.device.GetAI(i)
		if err != nil {
			return err
		}
		val, ok := r.GetAIValue(i, ai.GetConfig().Point)
		if ok {
			ch <- &edge.MeasureData{
				Name: ai.GetConfig().TagName,
				Tags: map[string]string{},
				Fields: map[string]interface{}{
					"value": val,
					"alarm": ai.CheckAlarm(val),
				},
			}
		}
	}

	for i := 0; i < r.AONum(); i++ {
		ao, err := e.device.GetAO(i)
		if err != nil {
			return err
		}
		val, ok := r.GetAOValue(i)
		if ok {
			ch <- &edge.MeasureData{
				Name: ao.GetConfig().TagName,
				Tags: map[string]string{},
				Fields: map[string]interface{}{
					"value": val,
				},
			}
		}
	}

	for i := 0; i < r.DINum(); i++ {
		di, err := e.device.GetDI(i)
		if err != nil {
			return err
		}
		val, ok := r.GetDIValue(i)
		if ok {
			ch <- &edge.MeasureData{
				Name: di.GetConfig().TagName,
				Tags: map[string]string{},
				Fields: map[string]interface{}{
					"value": util.If(val, 1, 0),
				},
			}
		}
	}

	for i := 0; i < r.DONum(); i++ {
		do, err := e.device.GetDO(i)
		if err != nil {
			return err
		}
		val, ok := r.GetDOValue(i)
		if ok {
			ch <- &edge.MeasureData{
				Name: do.GetConfig().TagName,
				Tags: map[string]string{},
				Fields: map[string]interface{}{
					"value": util.If(val, 1, 0),
				},
			}
		}
	}

	return nil
}

func (e *ep6v2Adapter) Create(option map[string]interface{}) (<-chan *edge.MeasureData, error) {
	delay := time.Second * 10
	if v, ok := option["delay"].(time.Duration); ok {
		delay = v
	}

	var ch = make(chan *edge.MeasureData)
	var done = make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-e.ctx.Done():
				return
			case <-time.After(delay):
				e.cond.L.Lock()
				e.cond.Wait()
				if err := e.fetchData(ch); err != nil {
					log.Error(err)
				}
				e.cond.L.Unlock()
				e.cond.Signal()
			}
		}
	}()

	e.channels.Store(ch, func() {
		close(done)
		wg.Wait()
	})
	return ch, nil
}

func (e *ep6v2Adapter) Drop(ch interface{}) {
	if fn, ok := e.channels.Load(ch); ok && fn != nil {
		fn.(func())()
	}
}

func (e *ep6v2Adapter) ProcessCtrlData(Values map[string]interface{}, ch chan error) {
	for k, v := range Values {
		err := e.device.SetCH(k, v)
		if err != nil {
			log.Errorf("ProcessCtrlData: %s", err)
			if ch != nil {
				ch <- err
			}
		}
	}
}

func (e *ep6v2Adapter) Plug(ch <-chan *edge.CtrlData) error {
	var done = make(chan struct{})
	var wg sync.WaitGroup

	e.channels.Store(ch, func() {
		close(done)
		wg.Wait()
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-e.ctx.Done():
				return
			case data := <-ch:
				if data != nil {
					e.ProcessCtrlData(data.Values, data.Error)
				} else {
					return
				}
			}
		}
	}()

	return nil
}

func (e *ep6v2Adapter) Stats() (map[string]interface{}, error) {
	panic("implement me")
}

func (e *ep6v2Adapter) Set(params map[string]interface{}) error {
	panic("implement me")
}

func (e *ep6v2Adapter) Get(keys ...string) (map[string]interface{}, error) {
	panic("implement me")
}
