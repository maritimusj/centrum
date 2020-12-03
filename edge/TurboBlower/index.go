package TurboBlower

import (
	"strings"

	"github.com/maritimusj/centrum/edge/devices/measure"
)

const (
	Separator = "*"
)

var (
	alarmTitle1 = []string{"Alarm", "Normal"}
	alarmTitle2 = []string{"ON", "OFF"}
)

func isSet(v uint16, n int) bool {
	return v>>n&0x01 == 1
}

var (
	m = map[string]func(v uint16) []*measure.Data{
		"1": f1,
		"2": f2,
		"3": f3,
		"4": f4,
		"5": f5,
	}
)

func Match(name string) bool {
	return strings.Contains(name, Separator)
}

func Analysis(name string, v int) []*measure.Data {
	arr := strings.SplitN(name, Separator, 2)
	if len(arr) != 2 {
		return nil
	}
	fn, exists := m[arr[1]]
	if !exists {
		return nil
	}

	return fn(uint16(v))
}

func newData(name string, alarm bool, title []string) *measure.Data {
	data := measure.New(name)

	data.AddTag("tag", "AI-"+name)
	data.AddTag("title", name)

	var alarmTitle string
	if alarm {
		alarmTitle = title[0]
	} else {
		alarmTitle = title[1]
	}

	data.AddTag("alarm", alarmTitle)
	data.AddTag("threshold", title[1])

	var v = 0
	if alarm {
		v = 1
	}

	data.AddField("val", v)
	return data
}

/**
40007
	0	紧急停止
	1	EOCR跳闸
	4 	变频器反馈错误
	5	喘振跳闸
	14	变频器通信错误
	15	远程通信错误
*/
func f1(v uint16) []*measure.Data {
	dataList := make([]*measure.Data, 0, 6)
	dataList = append(dataList,
		newData("紧急停止", isSet(v, 0), alarmTitle1),
		newData("EOCR跳闸", isSet(v, 1), alarmTitle1),
		newData("变频器反馈错误", isSet(v, 4), alarmTitle1),
		newData("喘振跳闸", isSet(v, 5), alarmTitle1),
		newData("变频器通信错误", isSet(v, 14), alarmTitle1),
		newData("远程通信错误", isSet(v, 15), alarmTitle1),
	)
	return dataList
}

/**
40008
	1	出口压力超高跳闸
	2	过滤器超压跳闸
	3	泵压力超高跳闸
	4	泵压力过低跳闸
	5	吸气温度过高跳闸
	7	电机温度过高跳闸
	10	变频器超温跳闸
	12	吸入压力传感器断开
	13	出口压力传感器断开
	14	过滤器压力传感器断开
	15	泵压力传感器断开
*/
func f2(v uint16) []*measure.Data {
	dataList := make([]*measure.Data, 0, 11)
	dataList = append(dataList,
		newData("出口压力超高跳闸", isSet(v, 1), alarmTitle1),
		newData("过滤器超压跳闸", isSet(v, 2), alarmTitle1),
		newData("泵压力超高跳闸", isSet(v, 3), alarmTitle1),
		newData("泵压力过低跳闸", isSet(v, 4), alarmTitle1),
		newData("吸气温度过高跳闸", isSet(v, 5), alarmTitle1),
		newData("电机温度过高跳闸", isSet(v, 7), alarmTitle1),
		newData("变频器超温跳闸", isSet(v, 10), alarmTitle1),
		newData("吸入压力传感器断开", isSet(v, 12), alarmTitle1),
		newData("出口压力传感器断开", isSet(v, 13), alarmTitle1),
		newData("过滤器压力传感器断开", isSet(v, 14), alarmTitle1),
		newData("泵压力传感器断开", isSet(v, 15), alarmTitle1),
	)
	return dataList
}

/**
40009
	0	变频器未知故障
	1	变频器过电压
	2	变频器欠电压
	3	变频器直联打开
	4	变频器轮廓打开
	5	变频器过热
	6	变频器保险丝开路
	7	变频器过载
	8	变频器过电流
	9	变频器频率过高
	10	变频器零序电流
	11	变频器装置短路
	12	变频器modbus错误
	13	变频器风扇错误
	14	电机过电流
*/
func f3(v uint16) []*measure.Data {
	dataList := make([]*measure.Data, 0, 15)
	dataList = append(dataList,
		newData("变频器未知故障", isSet(v, 0), alarmTitle1),
		newData("变频器过电压", isSet(v, 1), alarmTitle1),
		newData("变频器欠电压", isSet(v, 2), alarmTitle1),
		newData("变频器直联打开", isSet(v, 3), alarmTitle1),
		newData("变频器轮廓打开", isSet(v, 4), alarmTitle1),
		newData("变频器过热", isSet(v, 5), alarmTitle1),
		newData("变频器保险丝开路", isSet(v, 6), alarmTitle1),
		newData("变频器过载", isSet(v, 7), alarmTitle1),
		newData("变频器过电流", isSet(v, 8), alarmTitle1),
		newData("变频器频率过高", isSet(v, 9), alarmTitle1),
		newData("变频器零序电流", isSet(v, 10), alarmTitle1),
		newData("变频器装置短路", isSet(v, 11), alarmTitle1),
		newData("变频器modbus错误", isSet(v, 12), alarmTitle1),
		newData("变频器风扇错误", isSet(v, 13), alarmTitle1),
		newData("电机过电流", isSet(v, 14), alarmTitle1),
	)
	return dataList
}

/**
40010

	0	本地准备状态
	1	远程准备状态
	2	风机运行状态
	3	风机报警状态
	4	风机故障状态
	5	电机运行状态
	8	定频率支行模式状态
	9	定流量运行模式状态
	10	定功率运行模式状态
	11	比例控制运行模式状态
	12	溶解氧运行模式状态
	13	恒压运行模式
	15	DCS 通讯检查脉冲

*/
func f4(v uint16) []*measure.Data {
	dataList := make([]*measure.Data, 0, 13)
	dataList = append(dataList,
		newData("本地准备状态", isSet(v, 0), alarmTitle2),
		newData("远程准备状态", isSet(v, 1), alarmTitle2),
		newData("风机运行状态", isSet(v, 2), alarmTitle2),
		newData("风机报警状态", isSet(v, 3), alarmTitle2),
		newData("风机故障状态", isSet(v, 4), alarmTitle2),
		newData("电机运行状态", isSet(v, 5), alarmTitle2),
		newData("定频率支行模式状态", isSet(v, 8), alarmTitle2),
		newData("定流量运行模式状态", isSet(v, 9), alarmTitle2),
		newData("定功率运行模式状态", isSet(v, 10), alarmTitle2),
		newData("比例控制运行模式状态", isSet(v, 11), alarmTitle2),
		newData("溶解氧运行模式状态", isSet(v, 12), alarmTitle2),
		newData("恒压运行模式", isSet(v, 13), alarmTitle2),
		newData("DCS 通讯检查脉冲", isSet(v, 15), alarmTitle2),
	)
	return dataList
}

/**
40011
	0	吸入压力过高报警
	1	排气压力过高报警
	2	过滤压力过高报警
	3	水泵压力过高报警
	4	水泵压力过低报警
	5	吸气温度过高报警
	6	排气温度过高报警
	7	电机温度过高报警
	8	外界温度过高报警
	9	外界温度过低报警
	10	变频器温度过高报警
	11	喘振控制器报警
	14	压力传感器断开警报
	15	温度传感器断开警报
*/
func f5(v uint16) []*measure.Data {
	dataList := make([]*measure.Data, 0, 14)
	dataList = append(dataList,
		newData("吸入压力过高报警", isSet(v, 0), alarmTitle1),
		newData("排气压力过高报警", isSet(v, 1), alarmTitle1),
		newData("过滤压力过高报警", isSet(v, 2), alarmTitle1),
		newData("水泵压力过高报警", isSet(v, 3), alarmTitle1),
		newData("水泵压力过低报警", isSet(v, 4), alarmTitle1),
		newData("吸气温度过高报警", isSet(v, 5), alarmTitle1),
		newData("排气温度过高报警", isSet(v, 6), alarmTitle1),
		newData("电机温度过高报警", isSet(v, 7), alarmTitle1),
		newData("外界温度过高报警", isSet(v, 8), alarmTitle1),
		newData("外界温度过低报警", isSet(v, 9), alarmTitle1),
		newData("变频器温度过高报警", isSet(v, 10), alarmTitle1),
		newData("喘振控制器报警", isSet(v, 11), alarmTitle1),
		newData("压力传感器断开警报", isSet(v, 14), alarmTitle1),
		newData("温度传感器断开警报", isSet(v, 15), alarmTitle1),
	)
	return dataList
}
