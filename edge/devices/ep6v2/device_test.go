package ep6v2

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestDevice(t *testing.T) {
	device := New()
	err := device.Connect(context.Background(), "192.168.1.50:502")
	if err != nil {
		log.Fatal(err)
	}

	defer device.Close()

	model, err := device.GetModel()
	if err != nil {
		log.Fatal(err)
	}

	println(model.ID)

	r, err := device.GetRealTimeData()
	if err != nil {
		log.Fatal(err)
	}
	println("AI:", r.AINum())
	for i := 0; i < r.AINum(); i++ {
		v, ok := r.GetAIValue(i, 2)
		ai, _ := device.GetAI(i)
		v2, alarm, _ := ai.GetValue()
		fmt.Printf("%.2f => %.2f (%s), %#v\r\n", v, v2, FormatAlarm(alarm), ok)
	}

	println("DI:", r.DINum())
	for i := 0; i < r.DINum(); i++ {
		v, ok := r.GetDIValue(i)
		di, _ := device.GetDI(i)
		v2, _ := di.GetValue()
		fmt.Printf("%#v => %#v, %#v\r\n", v, v2, ok)
	}

	println("AO:", r.AONum())
	for i := 0; i < r.AONum(); i++ {
		v, ok := r.GetAOValue(i)
		fmt.Printf("%#v , %#v\r\n", v, ok)
	}

	//data, err := device.client.ReadInputRegisters(106, 4)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//println(ToSingle(data))
	//
	//data, err = device.client.ReadInputRegisters(0, 5)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%d\r\n", binary.BigEndian.Uint32(data))

	//for index := 0; index < chNum.AI; index ++ {
	//	ai, err := device.GetAI(index)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	v, err := ai.GetValue(device.client)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("%s (%s)=> %.2f\r\n", ai.config.Title, ai.config.TagName, v)
	//}
	//
	//for index := 0; index < chNum.AO; index ++ {
	//	ao, err := device.GetAO(index)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	v, err := ao.GetValue(device.client)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("%s (%s)=> %.2f\r\n", ao.config.Title, ao.config.TagName, v)
	//}
	//
	//for index := 0; index < chNum.DI; index ++ {
	//	di, err := device.GetDI(index)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	v, err := di.GetValue(device.client)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("%s (%s)=> %d\r\n", di.config.Title, di.config.TagName, v)
	//}

	//for index := 0; index < chNum.DO; index ++ {
	//	do, err := device.GetDO(index)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	v, err := do.GetValue(device.client)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("old:%s (%s)=> %#v\r\n", do.config.Title, do.config.TagName, v)
	//	v, err = do.SetValue(device.client, !v)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("new:%s (%s)=> %#v\r\n", do.config.Title, do.config.TagName, v)
	//}
}
