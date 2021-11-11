package canUtil

import (
	"encoding/json"
	"fmt"
	"gopublic/requester"
	"net/http"
)

type StandardSignal struct {
	Id                   int    `orm:"column(id);auto" json:"id"`
	MsgId                string `orm:"column(msg_id);null" json:"msgId"`
	SignalCode           string `orm:"column(signal_code)" json:"signalCode"` //内部定义code，每一code都是唯一的
	Status               int    `orm:"column(status)" json:"status"`          //是否是关键信号，已经在使用的为1，未使用未-1
	EcuName              string `orm:"column(ecu_name);null" json:"ecuName"`
	SignalName           string `orm:"column(signal_name);null" json:"signalName"`
	SignalDesc           string `orm:"column(signal_desc);null" json:"signalDesc"`
	ByteOrder            string `orm:"column(byte_order);null" json:"byteOrder"`
	StartByte            string `orm:"column(start_byte);null" json:"startByte"`
	StartBit             string `orm:"column(start_bit);null" json:"startBit"`
	SignalSendType       string `orm:"column(signal_send_type);null" json:"signalSendType"`
	BitLength            string `orm:"column(bit_length);null" json:"bitLength"`
	DataType             string `orm:"column(data_type);null" json:"dataType"`
	Resolution           string `orm:"column(resolution);null" json:"resolution"`
	Offset               string `orm:"column(offset);null" json:"offset"`
	SignalMin            string `orm:"column(signal_min);null" json:"signalMin"`
	SignalMax            string `orm:"column(signal_max);null" json:"signalMax"`
	SignalBusMin         string `orm:"column(signal_bus_min);null" json:"signalBusMin"`
	SignalBusMax         string `orm:"column(signal_bus_max);null" json:"signalBusMax"`
	InitialVal           string `orm:"column(initial_val);null" json:"initialVal"`
	InvalidVal           string `orm:"column(invalid_val);null" json:"invalidVal"`
	InactiveVal          string `orm:"column(inactive_val);null" json:"inactiveVal"`
	Unit                 string `orm:"column(unit);null" json:"unit"`
	SignalValDesc        string `orm:"column(signal_val_desc);null" json:"signalValDesc"`
	SignalValDescStr     string `orm:"column(signal_val_desc_str)" json:"signalValDescStr"`
	GroupMark            string `orm:"column(group_mark)" json:"groupMark"`
	InnerModeName        string `orm:"column(inner_mode_name)" json:"innerModeName"`
	SecondPackageVersion string `orm:"column(second_package_version);null" json:"secondPackageVersion"`
	Repeat               string `orm:"column(repeat)" json:"repeat"` //重复标记，空为无重复，非空有有重复，重复的signal_code
}

type CanIDSignals struct {
	Totals  int64            `json:"totals"`
	Version int64            `json:"version"`
	Signals []StandardSignal `json:"signals"`
}

type CanSignalRet struct {
	Code        int          `json:"code"`
	Message     string       `json:"message"`
	BusinessObj CanIDSignals `json:"businessObj"`
}

var (
	SignalMapME7 map[string]StandardSignal
	SignalMapME5 map[string]StandardSignal
)

func LoadCanMatrix(url string) {
	SignalMapME7, _ = GetSignalMapByInnerModel(url, "ME7")
	SignalMapME5, _ = GetSignalMapByInnerModel(url, "ME5")
}

func GetSignalMapByInnerModel(url string, innerModel string) (map[string]StandardSignal, error) {
	ret, err := requester.Instance().GetWithQuery(url, map[string]string{
		"innerModeName": innerModel,
	})
	if nil != err {
		return nil, err
	}
	if ret.Response().StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not 200")
	}
	var retSignal CanSignalRet
	if err := json.Unmarshal(ret.Bytes(), &retSignal); err != nil {
		return nil, err
	}

	sMap := make(map[string]StandardSignal, len(retSignal.BusinessObj.Signals))
	for k, v := range retSignal.BusinessObj.Signals {
		sMap[v.SignalName] = retSignal.BusinessObj.Signals[k]
	}

	return sMap, nil
}
