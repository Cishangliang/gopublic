package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/jiehuzen/gopublic"
)

func main() {
	test()
}

func test() {
	url := "http://localgw.cloud.enovatemotors.com/can-matrix/api/can/signal/standard"
	canUtil.LoadCanMatrix(url)
	log.Info(len(canUtil.SignalMapME5))
	log.Info(len(canUtil.SignalMapME7))
	signals := []string{"BCM_TurnRightLight_Sts", "BCM_RR_Door_Sts", "BCM_RL_Door_Sts", "BCM_RearViewHeatSts", "BCM_DriverLock_Sts", "BCM_SupplyVoltage"}
	//binaryMsg, _ := hex.DecodeString("240202AA0200007D")

	//bin := canUtil.Hex2bin("240202AA0200007D")
	//log.Info(string(bin))

	log.Info("----------------------")
	result := canUtil.GetSignalMap("240202AA0200007D", canUtil.SignalMapME7, signals, true)
	log.Info(result)

	signals1 := []string{"TPMS_FR_Pressure_value", "TPMS_RL_Pressure_value", "TPMS_RR_Pressure_value"}
	result = canUtil.GetSignalMap("AEA6B6A8524D4B52", canUtil.SignalMapME7, signals1, true)
	log.Info(result)

	signals2 := []string{"BMS_CCTemp1", "BMS_MinTemp"}
	result = canUtil.GetSignalMap("25680665533939F7", canUtil.SignalMapME7, signals2, false)
	log.Info(result)

	//bytes := []byte{49, 48}
	//parseInt, _ := strconv.ParseInt(string(bytes), 2, 64)
	//log.Info("--------", parseInt)
}
