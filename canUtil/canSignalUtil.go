package canUtil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

var motorolaArrays = [64]int{
	7, 6, 5, 4, 3, 2, 1, 0,
	15, 14, 13, 12, 11, 10, 9, 8,
	23, 22, 21, 20, 19, 18, 17, 16,
	31, 30, 29, 28, 27, 26, 25, 24,
	39, 38, 37, 36, 35, 34, 33, 32,
	47, 46, 45, 44, 43, 42, 41, 40,
	55, 54, 53, 52, 51, 50, 49, 48,
	63, 62, 61, 60, 59, 58, 57, 56,
}
var BYTEORDER_LSB = "lsb"
var BYTEORDER_MSB = "msb"

var BMS_BattTempFrame = []string{
	"BMS_BattTempFrame1",
	"BMS_BattTempFrame2",
	"BMS_BattTempFrame3",
	"BMS_BattTempFrame4",
	"BMS_BattTempFrame5",
}
var BMS_BattTempSensorNum_S = [][]string{
	{
		"BMS_BattTempSensorNum1",
		"BMS_BattTempSensorNum2",
		"BMS_BattTempSensorNum3",
		"BMS_BattTempSensorNum4",
		"BMS_BattTempSensorNum5",
		"BMS_BattTempSensorNum6",
		"BMS_BattTempSensorNum7",
	},
	{
		"BMS_BattTempSensorNum8",
		"BMS_BattTempSensorNum9",
		"BMS_BattTempSensorNum10",
		"BMS_BattTempSensorNum11",
		"BMS_BattTempSensorNum12",
		"BMS_BattTempSensorNum13",
		"BMS_BattTempSensorNum14",
	},
	{
		"BMS_BattTempSensorNum15",
		"BMS_BattTempSensorNum16",
		"BMS_BattTempSensorNum17",
		"BMS_BattTempSensorNum18",
		"BMS_BattTempSensorNum19",
	},
	{
		//填充位
	},
	{
		//填充位
	},
}

var BMS_BattTempSensorNum_D = [][]string{
	{
		"BMS_BattTempSensorNum36",
		"BMS_BattTempSensorNum37",
		"BMS_BattTempSensorNum38",
		"BMS_BattTempSensorNum39",
		"BMS_BattTempSensorNum40",
		"BMS_BattTempSensorNum41",
		"BMS_BattTempSensorNum42",
	},
	{
		"BMS_BattTempSensorNum43",
		"BMS_BattTempSensorNum44",
		"BMS_BattTempSensorNum45",
		"BMS_BattTempSensorNum46",
		"BMS_BattTempSensorNum47",
		"BMS_BattTempSensorNum48",
		"BMS_BattTempSensorNum49",
	},
	{
		"BMS_BattTempSensorNum50",
		"BMS_BattTempSensorNum51",
		"BMS_BattTempSensorNum52",
		"BMS_BattTempSensorNum53",
		"BMS_BattTempSensorNum54",
	},
}

var BMS_BattTempSensorNum_R = [][]string{
	{

	},
	{

	},
	{
		"BMS_BattTempSensorNum20",
		"BMS_BattTempSensorNum21",
	},
	{
		"BMS_BattTempSensorNum22",
		"BMS_BattTempSensorNum23",
		"BMS_BattTempSensorNum24",
		"BMS_BattTempSensorNum25",
		"BMS_BattTempSensorNum26",
		"BMS_BattTempSensorNum27",
		"BMS_BattTempSensorNum28",
	},
	{
		"BMS_BattTempSensorNum29",
		"BMS_BattTempSensorNum30",
		"BMS_BattTempSensorNum31",
		"BMS_BattTempSensorNum32",
		"BMS_BattTempSensorNum33",
		"BMS_BattTempSensorNum34",
		"BMS_BattTempSensorNum35",
	},
}

func Hex2bin(canMsg string) (binaryMsg []byte) {
	var ret string
	hexMsg, _ := hex.DecodeString(canMsg)
	for i := 0; i < len(hexMsg); i++ {
		r := fmt.Sprintf("%08b", hexMsg[i])
		ret += r
	}

	return []byte(ret)
}

/**
 * using CanSignalUtil.transformationDecimal
 * 4	218
 * 0	73
 * 0	125
 * 0	87
 * 0	68
 * 0	110
 * 0	101
 * 0	46
 * 0	18
 * 1	41
 * skip CanSignalUtil.transformationDecimal
 * 3	52
 * 0	10
 * 0	11
 * 0	11
 * 1	9
 * 1	7
 * 0	8
 * 0	8
 * 0	5
 * 0	16
 * 获取8字节内的信号解析详情
 * @param canMsg 16进制原始信号
 * @param matrix can矩阵表map，key为信号名，value为信号说明
 * @param signals 需要解析的信号名字列表
 * @param valid 是否过滤无效值,true过滤无效值，false不过滤无效值
 * @return 需要解析的信号列表
 */
func GetSignalMap(canMsg string, matrix map[string]StandardSignal, signals []string, invalid bool) map[string]string {

	m := make(map[string]string)
	//hex to byte
	binaryMsg := Hex2bin(canMsg)

	for i := 0; i < len(signals); i++ {
		signal := matrix[signals[i]]
		startBit, _ := strconv.Atoi(signal.StartBit)
		bitLen, _ := strconv.Atoi(signal.BitLength)
		result := transformationDecimal(binaryMsg, BYTEORDER_MSB, startBit, bitLen, signal.Resolution, signal.Offset)

		signalMin, _ := decimal.NewFromString(signal.SignalMin)
		signalMax, _ := decimal.NewFromString(signal.SignalMax)
		if invalid && (result.Cmp(signalMin) == -1 || result.Cmp(signalMax) == 1) {
			m[signal.SignalName] = "invalid"
		} else {
			m[signal.SignalName] = result.String()
		}
	}

	//电池监听映射。信号映射 todo：CLOUD-4500
	//for i := 0; i < len(BMS_BattTempFrame); i++ {
	//	if m[BMS_BattTempFrame[i]] == "1"{
	//		for j := 0; j < len(BMS_BattTempSensorNum_S[i]); j++ {
	//
	//		}
	//	}
	//}

	return m
}

/**
 * @param binaryString 二进制字符串
 * @param byteOrder    排序方式（LSB|MSB）
 * @param startBit     起始位
 * @param signalLength 信号长度
 * @param resolution   精度
 * @param offset       偏移量
 * @return 十进制信号, 物理值
 */
func transformationDecimal(binaryMsg []byte, byteOrder string, startBit int, bitLen int, resolution string, offset string) decimal.Decimal {
	values, _ := msgSplit(binaryMsg, byteOrder, startBit, bitLen)

	dbc_value := decimal.NewFromInt(values)
	dbc_resolution, _ := decimal.NewFromString(resolution)
	dbc_offset, _ := decimal.NewFromString(offset)

	result := dbc_value.Mul(dbc_resolution).Add(dbc_offset).Round(2)

	return result
}

func msgSplit(binaryMsg []byte, byteOrder string, startBit int, bitLen int) (int64, error) {
	index := motorolaArrays[startBit]
	signalBits := make([]byte, 0)

	if byteOrder == BYTEORDER_LSB {
		for i := index; i < index-bitLen; i-- {
			signalBits = append(signalBits, binaryMsg[i])
		}
	} else if byteOrder == BYTEORDER_MSB {
		for i := index; i < index+bitLen; i++ {
			signalBits = append(signalBits, binaryMsg[i])
		}
	} else {
		return 0, errors.New("the parameter byteOrder not specified or incorrect!")
	}

	retInt, _ := strconv.ParseInt(string(signalBits), 2, 64)

	return retInt, nil
}
