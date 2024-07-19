package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func readFileContent(t *testing.T, path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	return content
}

func TestGetURLWorksForIP(t *testing.T) {
	u, e := getURL("192.168.1.1")
	assert.Nil(t, e)
	assert.Equal(t, "http://192.168.1.1/rawData", u)
}

func TestGetURLWorksWithSingleSlashPathProvided(t *testing.T) {
	u, e := getURL("192.168.1.1/")
	assert.Nil(t, e)
	// should use the single slash as the path, assuming we're using some kind of proxy here or something
	assert.Equal(t, "http://192.168.1.1/", u)
}

func TestGetURLWorksWithWholePathProvided(t *testing.T) {
	u, e := getURL("192.168.1.1/some/sub/path")
	assert.Nil(t, e)
	// should use the single slash as the path, assuming we're using some kind of proxy here or something
	assert.Equal(t, "http://192.168.1.1/some/sub/path", u)
}

func TestGetURLWorksWithHTTPS(t *testing.T) {
	u, e := getURL("https://192.168.1.1")
	assert.Nil(t, e)
	assert.Equal(t, "https://192.168.1.1/rawData", u)
}

func TestGetURLWorksWithHTTP(t *testing.T) {
	u, e := getURL("http://192.168.1.1")
	assert.Nil(t, e)
	assert.Equal(t, "http://192.168.1.1/rawData", u)
}

func TestGetURLWorksWithHostname(t *testing.T) {
	u, e := getURL("http://sbms.local")
	assert.Nil(t, e)
	assert.Equal(t, "http://sbms.local/rawData", u)
}

func TestGetURLWorksWithHostnameAndPath(t *testing.T) {
	u, e := getURL("https://sbms.local/")
	assert.Nil(t, e)
	assert.Equal(t, "https://sbms.local/", u)
}

func TestRawData6(t *testing.T) {
	content := readFileContent(t, "./__source__/rawData6")
	out := decodeResponse(content)
	log.Printf("%+#v", out)
	assert.Equal(t, &SBMSData{ts: "2024-02-20T13:32:56",
		soc:            69,
		batteryVoltage: 26479,
		batteryPower:   -139.06770799999998,
		cells: []Cell{
			{mV: 3310, isBalancing: false},
			{mV: 3313, isBalancing: false},
			{mV: 3312, isBalancing: false},
			{mV: 3311, isBalancing: false},
			{mV: 3310, isBalancing: false},
			{mV: 3308, isBalancing: false},
			{mV: 3307, isBalancing: false},
			{mV: 3308, isBalancing: false},
		},
		minMV:               2500,
		maxMV:               3750,
		internalTemperature: 26,
		externalTemperature: -45,
		batteryCurrent:      -5252,
		pv1Current:          39,
		pv2Current:          0,
		externalCurrent:     5234,
		adc2:                0,
		adc3:                0,
		adc4:                0,
		heat1:               0,
		heat2:               11457,
		flags: Flags{
			OverVoltage:             false,
			OverVoltageLock:         false,
			UnderVoltage:            false,
			UnderVoltageLock:        false,
			InternalOverTemperature: false,
			ChargeOverCurrent:       false,
			DischargeOverCurrent:    false,
			DischargeShortCircuit:   false,
			CellFail:                false,
			OpenCellWire:            false,
			LowVoltageCell:          false,
			EEPROMFail:              false,
			ChargeFETActive:         true,
			EndOfCharge:             false,
			DischargeFETActive:      true},
		batteryEnergyWh: 398056.7,
		batteryEnergyAh: 14945.571,
		pV1EnergyWh:     725043.8,
		pV1EnergyAh:     26846.452,
		pV2EnergyWh:     0,
		pV2EnergyAh:     0,
		loadEnergyWh:    256912.7,
		loadEnergyAh:    9567.035,
		extLoadEnergyWh: 395565.5,
		extLoadEnergyAh: 14852.17,
		dmpptEnergyWh:   0,
		dmpptEnergyAh:   0,
		cellType:        1,
		capacity:        280,
		status:          20480,
	}, out)
}

func TestRawData7(t *testing.T) {
	content := readFileContent(t, "./__source__/rawData7")
	out := decodeResponse(content)
	log.Printf("%+#v", out)
	assert.Equal(t, &SBMSData{
		ts:             "2024-03-03T07:44:24",
		soc:            96,
		batteryVoltage: 26612,
		batteryPower:   -114.24531599999999,
		cells: []Cell{
			{mV: 3326, isBalancing: false},
			{mV: 3328, isBalancing: false},
			{mV: 3330, isBalancing: false},
			{mV: 3328, isBalancing: false},
			{mV: 3325, isBalancing: false},
			{mV: 3326, isBalancing: false},
			{mV: 3325, isBalancing: false},
			{mV: 3324, isBalancing: false},
		},
		minMV:               2500,
		maxMV:               3750,
		internalTemperature: 28.9,
		externalTemperature: -45,
		batteryCurrent:      -4293,
		pv1Current:          1367,
		pv2Current:          0,
		externalCurrent:     4271,
		adc2:                0,
		adc3:                0,
		adc4:                0,
		heat1:               0,
		heat2:               11370,
		flags: Flags{
			OverVoltage:             false,
			OverVoltageLock:         false,
			UnderVoltage:            false,
			UnderVoltageLock:        false,
			InternalOverTemperature: false,
			ChargeOverCurrent:       false,
			DischargeOverCurrent:    false,
			DischargeShortCircuit:   false,
			CellFail:                false,
			OpenCellWire:            false,
			LowVoltageCell:          false,
			EEPROMFail:              false,
			ChargeFETActive:         true,
			EndOfCharge:             false,
			DischargeFETActive:      true},
		batteryEnergyWh: 424697,
		batteryEnergyAh: 15949.933,
		pV1EnergyWh:     774548.8,
		pV1EnergyAh:     28687.449,
		pV2EnergyWh:     0,
		pV2EnergyAh:     0,
		loadEnergyWh:    276205.4,
		loadEnergyAh:    10287.851,
		extLoadEnergyWh: 422096.5,
		extLoadEnergyAh: 15852.42,
		dmpptEnergyWh:   0,
		dmpptEnergyAh:   0,
		cellType:        1,
		capacity:        280,
		status:          20480,
	}, out)
}

func TestRawData8(t *testing.T) {
	content := readFileContent(t, "./__source__/rawData8")
	out := decodeResponse(content)
	log.Printf("%+#v", out)
	assert.Equal(t, &SBMSData{
		ts:             "2024-03-03T06:44:24",
		soc:            100,
		batteryVoltage: 27017,
		batteryPower:   15.885995999999999,
		cells: []Cell{
			{mV: 3375, isBalancing: false},
			{mV: 3376, isBalancing: false},
			{mV: 3381, isBalancing: false},
			{mV: 3377, isBalancing: false},
			{mV: 3389, isBalancing: false},
			{mV: 3351, isBalancing: true},
			{mV: 3394, isBalancing: false},
			{mV: 3374, isBalancing: false},
		},
		minMV:               2500,
		maxMV:               3750,
		internalTemperature: 30,
		externalTemperature: -45,
		batteryCurrent:      588,
		pv1Current:          5512,
		pv2Current:          0,
		externalCurrent:     0,
		adc2:                0,
		adc3:                0,
		adc4:                0,
		heat1:               0,
		heat2:               11106,
		flags: Flags{
			OverVoltage:             false,
			OverVoltageLock:         false,
			UnderVoltage:            false,
			UnderVoltageLock:        false,
			InternalOverTemperature: false,
			ChargeOverCurrent:       false,
			DischargeOverCurrent:    false,
			DischargeShortCircuit:   false,
			CellFail:                false,
			OpenCellWire:            false,
			LowVoltageCell:          false,
			EEPROMFail:              false,
			ChargeFETActive:         true,
			EndOfCharge:             false,
			DischargeFETActive:      true,
		},
		batteryEnergyWh: 424445.2,
		batteryEnergyAh: 15940.414,
		pV1EnergyWh:     774464.5,
		pV1EnergyAh:     28684.309,
		pV2EnergyWh:     0,
		pV2EnergyAh:     0,
		loadEnergyWh:    276130.2,
		loadEnergyAh:    10285.05,
		extLoadEnergyWh: 421845.2,
		extLoadEnergyAh: 15842.915,
		dmpptEnergyWh:   0,
		dmpptEnergyAh:   0,
		cellType:        1,
		capacity:        280,
		status:          20480,
	}, out)
}

func TestRawData9(t *testing.T) {
	content := readFileContent(t, "./__source__/rawData9")
	out := decodeResponse(content)
	log.Printf("%+#v", out)
	assert.Equal(t, &SBMSData{ts: "2024-07-10T13:42:12",
		soc:            41,
		batteryVoltage: 23160,
		batteryPower:   16.37412,
		cells: []Cell{
			{mV: 2893, isBalancing: false},
			{mV: 2959, isBalancing: false},
			{mV: 2814, isBalancing: false},
			{mV: 2888, isBalancing: false},
			{mV: 2957, isBalancing: false},
			{mV: 2863, isBalancing: false},
			{mV: 2849, isBalancing: false},
			{mV: 2937, isBalancing: false},
		},
		minMV:               2500,
		maxMV:               3750,
		internalTemperature: 24.8,
		externalTemperature: -45,
		batteryCurrent:      707,
		pv1Current:          29,
		pv2Current:          0,
		externalCurrent:     0,
		adc2:                0,
		adc3:                0,
		adc4:                0,
		heat1:               0,
		heat2:               11318,
		flags: Flags{
			OverVoltage:             false,
			OverVoltageLock:         false,
			UnderVoltage:            true,
			UnderVoltageLock:        false,
			InternalOverTemperature: false,
			ChargeOverCurrent:       false,
			DischargeOverCurrent:    false,
			DischargeShortCircuit:   false,
			CellFail:                false,
			OpenCellWire:            false,
			LowVoltageCell:          false,
			EEPROMFail:              false,
			ChargeFETActive:         true,
			EndOfCharge:             false,
			DischargeFETActive:      false,
		},
		batteryEnergyWh: 580629.2,
		batteryEnergyAh: 21862.908,
		pV1EnergyWh:     999996.2,
		pV1EnergyAh:     37061.873,
		pV2EnergyWh:     0,
		pV2EnergyAh:     0,
		loadEnergyWh:    373157.6,
		loadEnergyAh:    13890.788,
		extLoadEnergyWh: 577024.6,
		extLoadEnergyAh: 21727.536,
		dmpptEnergyWh:   0,
		dmpptEnergyAh:   0,
		cellType:        1,
		capacity:        280,
		status:          4100},
		out)
}

func TestDebug1(t *testing.T) {
	content := readFileContent(t, "./__source__/debug1")
	out := decodeDebugResponse(content)
	assert.Equal(t, []SystemTaskInfo{
		{name: "async_tcp", state: 1, priority: 3, runTimeCounter: 7.37905955e+08, runTimePercent: 32},
		{name: "loopTask", state: 1, priority: 1, runTimeCounter: 4.14155237e+08, runTimePercent: 18},
		{name: "IDLE0", state: 1, priority: 0, runTimeCounter: 3.402661223e+09, runTimePercent: 150},
		{name: "IDLE1", state: 1, priority: 0, runTimeCounter: 438596, runTimePercent: 0},
		{name: "tiT", state: 2, priority: 18, runTimeCounter: 6.7312419e+08, runTimePercent: 29},
		{name: "Tmr Svc", state: 2, priority: 1, runTimeCounter: 68, runTimePercent: 0},
		{name: "network_event", state: 2, priority: 19, runTimeCounter: 305, runTimePercent: 0},
		{name: "ipc1", state: 2, priority: 24, runTimeCounter: 1.130147e+06, runTimePercent: 0},
		{name: "ipc0", state: 2, priority: 24, runTimeCounter: 5.165267e+06, runTimePercent: 0},
		{name: "mdns", state: 2, priority: 1, runTimeCounter: 5.165444e+06, runTimePercent: 0},
		{name: "esp_timer", state: 2, priority: 22, runTimeCounter: 1.95894524e+08, runTimePercent: 8},
		{name: "wifi", state: 2, priority: 23, runTimeCounter: 2.141062145e+09, runTimePercent: 94},
		{name: "uart", state: 2, priority: 15, runTimeCounter: 1.228550683e+09, runTimePercent: 54},
		{name: "sys_evt", state: 2, priority: 20, runTimeCounter: 1809, runTimePercent: 0},
	}, out)
}
