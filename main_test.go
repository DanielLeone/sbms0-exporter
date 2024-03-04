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

func TestRawData(t *testing.T) {
	content := readFileContent(t, "./__source__/rawData")
	out := decodeResponse(content)
	log.Printf("%+#v", out)

	assert.Equal(t, &SBMSData{
		ts:             "2024-03-03T06:44:23",
		soc:            100,
		batteryVoltage: 27020,
		batteryPower:   20.18394,
		cells: []Cell{
			{mV: 3375, isBalancing: false},
			{mV: 3376, isBalancing: false},
			{mV: 3382, isBalancing: false},
			{mV: 3378, isBalancing: false},
			{mV: 3397, isBalancing: false},
			{mV: 3335, isBalancing: true},
			{mV: 3402, isBalancing: false},
			{mV: 3375, isBalancing: false},
		},
		minMV:               2500,
		maxMV:               3750,
		internalTemperature: 30.3,
		externalTemperature: -45,
		batteryCurrent:      747,
		pv1Current:          5512,
		pv2Current:          0,
		externalCurrent:     0,
		adc2:                0,
		adc3:                0,
		adc4:                0,
		heat1:               0,
		heat2:               11090,
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
		pV1EnergyWh:     774464.4,
		pV1EnergyAh:     28684.307,
		pV2EnergyWh:     0,
		pV2EnergyAh:     0,
		loadEnergyWh:    276130.1,
		loadEnergyAh:    10285.048,
		extLoadEnergyWh: 421845.2,
		extLoadEnergyAh: 15842.915,
		dmpptEnergyWh:   0,
		dmpptEnergyAh:   0,
		cellType:        1,
		capacity:        280,
		status:          20480,
	}, out)
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
