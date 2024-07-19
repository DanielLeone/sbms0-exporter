package main

import "C"
import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"
)

var (
	respBytes = prometheus.NewCounter(prometheus.CounterOpts{Namespace: "sbms", Subsystem: "exporter", Name: "resp_bytes", Help: "Number of bytes received"})
	reqsCount = prometheus.NewCounter(prometheus.CounterOpts{Namespace: "sbms", Subsystem: "exporter", Name: "requests", Help: "Number of requests made"})
	Soc       = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sbms",
		Name:      "battery_soc",
		Help:      "State of Charge",
	})
	Cell1Voltage   = cellVoltageGauge(1)
	Cell2Voltage   = cellVoltageGauge(2)
	Cell3Voltage   = cellVoltageGauge(3)
	Cell4Voltage   = cellVoltageGauge(4)
	Cell5Voltage   = cellVoltageGauge(5)
	Cell6Voltage   = cellVoltageGauge(6)
	Cell7Voltage   = cellVoltageGauge(7)
	Cell8Voltage   = cellVoltageGauge(8)
	Cell1Balancing = cellBalancingGauge(1)
	Cell2Balancing = cellBalancingGauge(2)
	Cell3Balancing = cellBalancingGauge(3)
	Cell4Balancing = cellBalancingGauge(4)
	Cell5Balancing = cellBalancingGauge(5)
	Cell6Balancing = cellBalancingGauge(6)
	Cell7Balancing = cellBalancingGauge(7)
	Cell8Balancing = cellBalancingGauge(8)
	TempInt        = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sbms",
		Name:      "internal_temp",
	})
	TempExt = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sbms",
		Name:      "external_temp",
	})
	BatteryPower   = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "battery_power"})
	BatteryVoltage = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "battery_voltage"})
	BatteryCurrent = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "battery_current"})
	CurrentPv1     = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "pv1_current"})
	CurrentPv2     = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "pv2_current"})
	CurrentExtLoad = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "ext_current"})
	Ad2            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "ad2"})
	Ad3            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "ad3"})
	Ad4            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "ad4"})
	Heat1          = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "heat1"})
	Heat2          = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "heat2"})
	Ov             = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "ov", Help: "Over Voltage"})
	Ovlk           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "ovlk", Help: "Over Voltage Lock"})
	Uv             = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "uv", Help: "Under Voltage"})
	Uvlk           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "uvlk", Help: "Under Voltage Lock"})
	Iot            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "iot"})
	Coc            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "coc", Help: "Charge Over Current"})
	Doc            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "doc", Help: "Discharge Over Current"})
	Dsc            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "dsc", Help: "Discharge Short Circuit"})
	Celf           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "celf", Help: "Cell Fail"})
	Open           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "open", Help: "OpenCellWire Cell Wire"})
	Lvc            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "lvc", Help: "Low Voltage Cell"})
	Eccf           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "eccf"})
	Cfet           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "cfet", Help: "Charge FET"})
	Eoc            = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "eoc", Help: "End Of Charge"})
	Dfet           = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "flag", Name: "dfet", Help: "Discharge FET"})

	batteryEnergyWh = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "battery_wh"})
	batteryEnergyAh = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "battery_ah"})
	pV1EnergyWh     = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "pv1_wh"})
	pV1EnergyAh     = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "pv1_ah"})
	pV2EnergyWh     = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "pv2_wh"})
	pV2EnergyAh     = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "pv2_ah"})
	dmpptEnergyWh   = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "dmppt_wh"})
	dmpptEnergyAh   = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "dmppt_ah"})
	loadEnergyWh    = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "load_wh"})
	loadEnergyAh    = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "load_ah"})
	extLoadEnergyWh = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "ext_load_wh"})
	extLoadEnergyAh = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "energy", Name: "ext_load_ah"})
	cellType        = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "type"})
	capacity        = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "capacity"})
	status          = promauto.NewGauge(prometheus.GaugeOpts{Namespace: "sbms", Name: "status"})

	systemTaskPriority       = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "system", Name: "task_priority"}, []string{"task"})
	systemTaskRunTime        = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "system", Name: "task_run_time"}, []string{"task"})
	systemTaskRunTimePercent = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "system", Name: "task_run_time_percent"}, []string{"task"})
	systemTaskState          = promauto.NewGaugeVec(prometheus.GaugeOpts{Namespace: "sbms", Subsystem: "system", Name: "task_state"}, []string{"task"})
)

func cellVoltageGauge(idx int) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sbms",
		Name:      "cell_voltage",
		Help:      "Cell Voltage",
		ConstLabels: map[string]string{
			"cell": strconv.Itoa(idx),
		},
	})
}
func cellBalancingGauge(idx int) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sbms",
		Name:      "cell_balancing",
		Help:      "Cell Balancing",
		ConstLabels: map[string]string{
			"cell": strconv.Itoa(idx),
		},
	})
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

type Cell struct {
	mV          int
	isBalancing bool
}

type Flags struct {
	OverVoltage             bool
	OverVoltageLock         bool
	UnderVoltage            bool
	UnderVoltageLock        bool
	InternalOverTemperature bool
	ChargeOverCurrent       bool
	DischargeOverCurrent    bool
	DischargeShortCircuit   bool
	CellFail                bool
	OpenCellWire            bool
	LowVoltageCell          bool
	EEPROMFail              bool
	ChargeFETActive         bool
	EndOfCharge             bool
	DischargeFETActive      bool
}

type SystemTaskInfo struct {
	name           string
	state          float64
	priority       float64
	runTimeCounter float64
	runTimePercent float64
}

type SBMSData struct {
	ts                  string
	soc                 float64
	batteryVoltage      float64
	batteryPower        float64
	cells               []Cell
	minMV               int
	maxMV               int
	internalTemperature float64
	externalTemperature float64
	batteryCurrent      float64
	pv1Current          float64
	pv2Current          float64
	externalCurrent     float64
	adc2                int
	adc3                int
	adc4                int
	heat1               int
	heat2               int
	flags               Flags
	batteryEnergyWh     float64
	batteryEnergyAh     float64
	pV1EnergyWh         float64
	pV1EnergyAh         float64
	pV2EnergyWh         float64
	pV2EnergyAh         float64
	loadEnergyWh        float64
	loadEnergyAh        float64
	extLoadEnergyWh     float64
	extLoadEnergyAh     float64
	dmpptEnergyWh       float64
	dmpptEnergyAh       float64
	cellType            float64
	capacity            float64
	status              float64
}

func dcmp(offset, count int, runes []uint16) float64 {
	var sum float64 = 0
	var l = ""
	for z := 0; z < count; z++ {
		i := offset + count - 1 - z
		r := runes[i]
		n1 := r - 35
		n2 := math.Pow(91, float64(z))
		v := float64(n1) * n2
		sum = sum + v
		l += fmt.Sprintf("rune %d is %d is %s is %f; ", i, r, string(utf16.Decode([]uint16{r})), v)
	}
	log.Printf("offset=%d count=%d sum=%f debug=%s\n", offset, count, sum, l)
	return sum
}

func binToBool(bin rune) bool {
	if bin == '1' {
		return true
	} else {
		return false
	}
}

func decodeDebugResponse(b []byte) []SystemTaskInfo {
	var tasks []SystemTaskInfo
	value := string(b)
	for _, row := range strings.Split(value, "\n") {
		if len(strings.TrimSpace(row)) <= 0 {
			continue
		}
		task := new(SystemTaskInfo)
		realIdx := 0
		for _, column := range strings.Split(row, "\t") {
			// account for "double tabs" or empty columns
			trimmedValue := strings.TrimSpace(column)
			if len(trimmedValue) <= 0 {
				continue
			}
			switch realIdx {
			case 0:
				task.name = trimmedValue
				break
			case 1:
				task.state = taskStateToValue(trimmedValue)
				break
			case 2:
				task.priority, _ = strconv.ParseFloat(trimmedValue, 64)
				break
			case 3:
				task.runTimeCounter, _ = strconv.ParseFloat(trimmedValue, 64)
				break
			case 4:
				task.runTimePercent, _ = strconv.ParseFloat(strings.Trim(trimmedValue, "%"), 64)
				break
			}
			realIdx += 1
		}
		tasks = append(tasks, *task)
	}
	return tasks
}

func decodeResponse(b []byte) *SBMSData {
	output := new(SBMSData)
	data, err := parseRawData(b)
	if err != nil {
		return nil
	}

	sbms := data.sbms
	s2 := data.s2
	xsbms := data.xsbms
	eW := data.eW
	eA := data.eA

	Y := dcmp(0, 1, sbms)
	M := dcmp(1, 1, sbms)
	D := dcmp(2, 1, sbms)
	H := dcmp(3, 1, sbms)
	m := dcmp(4, 1, sbms)
	S := dcmp(5, 1, sbms)

	output.ts = fmt.Sprintf("20%d-%02d-%02dT%02d:%02d:%02d", int(Y), int(M), int(D), int(H), int(m), int(S))

	output.soc = dcmp(6, 2, sbms)
	output.cells = []Cell{}
	var batteryVoltage float64 = 0
	for i := 0; i < 8; i++ {
		cellMv := dcmp((i*2)+8, 2, sbms)
		batteryVoltage += cellMv
		output.cells = append(output.cells, Cell{
			mV:          int(cellMv),
			isBalancing: s2[i] == 1,
		})
	}
	output.batteryVoltage = batteryVoltage

	output.internalTemperature = (dcmp(24, 2, sbms) - 450) / 10
	output.externalTemperature = (dcmp(26, 2, sbms) - 450) / 10

	n := sbms[28]
	var scalar float64 = 1
	// 45 is negative sign in utf16
	if n == 45 {
		scalar = -1
	}

	output.batteryCurrent = dcmp(29, 3, sbms) * scalar
	output.pv1Current = dcmp(32, 3, sbms)
	output.pv2Current = dcmp(35, 3, sbms)
	output.externalCurrent = dcmp(38, 3, sbms)

	output.batteryPower = batteryVoltage / 1000 * output.batteryCurrent / 1000

	output.adc2 = int(dcmp(41, 3, sbms))
	output.adc3 = int(dcmp(44, 3, sbms))
	output.adc4 = int(dcmp(47, 3, sbms))
	output.heat1 = int(dcmp(50, 3, sbms))
	output.heat2 = int(dcmp(53, 3, sbms))

	// There are 15 Flags (see monitoring page3 on SBMS)
	// not all flags represent an error some are representing normal operation
	// and if you consider this a 15bit binary number
	// each bit representing one flag with “1” highlighted and “0” not highlighted
	// and you start from the top left with OV flag and continue to right and then down to the last flag DFET
	// and consider the OV as the least significant bit and DFET most significant bit
	// then you get this binary number and convert to decimal
	// that will be the priority contained in this ERR Status
	// and you can also see this number on the SBMS LCD just under the battery SOC priority and also displayed as Status: in the html page.

	// For some reason I think this is the order, but I can't remember where I found it:
	//  Ov Ovlk Uv Uvlk InternalOverTemperature ChargeOverCurrent DischargeOverCurrent DischargeShortCircuit CellFail OpenCellWire LowVoltageCell EEPROMFail ChargeFETActive EndOfCharge DischargeFETActive

	errorFlags := int64(dcmp(56, 3, sbms))
	errorBits := strconv.FormatInt(errorFlags, 2)

	// in some cases I found this wasn't 15 bits, I'm guessing we prefix with zeros?
	errorBitsPadded := fmt.Sprintf("%015s", errorBits)

	errorRunes := []rune(errorBitsPadded)
	output.flags = Flags{
		DischargeFETActive:      binToBool(errorRunes[0]),
		EndOfCharge:             binToBool(errorRunes[1]),
		ChargeFETActive:         binToBool(errorRunes[2]),
		EEPROMFail:              binToBool(errorRunes[3]),
		LowVoltageCell:          binToBool(errorRunes[4]),
		OpenCellWire:            binToBool(errorRunes[5]),
		CellFail:                binToBool(errorRunes[6]),
		DischargeShortCircuit:   binToBool(errorRunes[7]),
		DischargeOverCurrent:    binToBool(errorRunes[8]),
		ChargeOverCurrent:       binToBool(errorRunes[9]),
		InternalOverTemperature: binToBool(errorRunes[10]),
		UnderVoltageLock:        binToBool(errorRunes[11]),
		UnderVoltage:            binToBool(errorRunes[12]),
		OverVoltageLock:         binToBool(errorRunes[13]),
		OverVoltage:             binToBool(errorRunes[14]),
	}

	output.minMV = int(dcmp(5, 2, xsbms))
	output.maxMV = int(dcmp(3, 2, xsbms))

	// todo sometimes we get zero from these counters for some reason, we should potentially exclude those as invalid

	//Batt
	output.batteryEnergyWh = dcmp(0*6, 6, eW) / 10
	output.batteryEnergyAh = dcmp(0*6, 6, eA) / 1000

	//PV1
	output.pV1EnergyWh = dcmp(1*6, 6, eW) / 10
	output.pV1EnergyAh = dcmp(1*6, 6, eA) / 1000

	//PV2
	output.pV2EnergyWh = dcmp(2*6, 6, eW) / 10
	output.pV2EnergyAh = dcmp(2*6, 6, eA) / 1000

	//DMPPT
	output.dmpptEnergyWh = dcmp(3*6, 6, eW)
	output.dmpptEnergyAh = dcmp(3*6, 6, eA)

	//Load
	output.loadEnergyWh = dcmp(5*6, 6, eW) / 10
	output.loadEnergyAh = dcmp(5*6, 6, eA) / 1000

	//ExtLd
	output.extLoadEnergyWh = dcmp(6*6, 6, eW) / 10
	output.extLoadEnergyAh = dcmp(6*6, 6, eA) / 1000

	output.cellType = dcmp(7, 1, xsbms)
	output.capacity = dcmp(8, 3, xsbms)
	output.status = dcmp(56, 3, sbms)

	return output
}

type SBMSRawData struct {
	//Btn string
	//Btp string

	//ELd string
	//Ld string

	//PV1 string
	//PV2 string

	//dmppt string
	//gsbms string
	//s1    string

	eA    []uint16
	eW    []uint16
	s2    []int64
	sbms  []uint16
	xsbms []uint16
}

var re = regexp.MustCompile(`var\s+(\w+)\s*=\s*(.*);`)

func extractJSVariableLiteralValues(content []byte) map[string]string {
	output := map[string]string{}
	for _, line := range strings.Split(string(content), "\n") {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, pair := range matches {
			n := pair[1]
			v := pair[2]
			output[n] = v
		}
	}
	return output
}

func parseRawData(content []byte) (*SBMSRawData, error) {
	output := new(SBMSRawData)
	extracted := extractJSVariableLiteralValues(content)

	output.eW = extractStrLiteral(extracted["eW"])
	output.eA = extractStrLiteral(extracted["eA"])
	output.sbms = extractStrLiteral(extracted["sbms"])
	output.xsbms = extractStrLiteral(extracted["xsbms"])

	s2, err := extractIntArrayLiteral(extracted["s2"])
	if err != nil {
		return nil, err
	}
	output.s2 = s2

	return output, nil
}

func extractIntArrayLiteral(value string) ([]int64, error) {
	var out []int64
	escaped := strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
	split := strings.Split(escaped, ",")
	for _, v := range split {
		parseInt, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, parseInt)
	}

	log.Printf("priority=%s", value)
	log.Printf("escaped=%s", escaped)
	log.Printf("split=%#+v", split)
	log.Printf("out=%#+v", out)

	return out, nil
}

func extractStrLiteral(value string) []uint16 {
	stripped := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(value), "\""), "\"")
	// replace double backslash with a single backslash
	escape1 := strings.ReplaceAll(stripped, "\\\\", "\\")
	// replace backslash followed by a double quote with just a double quote
	escape2 := strings.ReplaceAll(escape1, "\\\"", "\"")
	runes := utf16.Encode([]rune(escape2))

	log.Printf("priority=%s", value)
	log.Printf("stripped=%s", stripped)
	log.Printf("escaped=%s", escape2)
	log.Printf("runes=%v", runes)

	return runes
}

type SBMS0Collector struct {
	url string
}
type SBMS0SystemCollector struct {
	url string
}

func (cc SBMS0Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

func (cc SBMS0SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect collects the rawData endpoint directly from the SBMS0 device
//
// Note that Collect could be called concurrently, so we depend on
// the /rawData endpoint to be concurrency-safe.
func (cc SBMS0Collector) Collect(ch chan<- prometheus.Metric) {
	reqsCount.Inc()
	resp, err := http.Get(cc.url)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	respBytes.Add(float64(len(b)))

	log.Printf("resp is\n%s\n", string(b))
	response := decodeResponse(b)

	ch <- respBytes
	ch <- reqsCount

	setAndExport(ch, Cell1Voltage, float64(response.cells[0].mV))
	setAndExport(ch, Cell2Voltage, float64(response.cells[1].mV))
	setAndExport(ch, Cell3Voltage, float64(response.cells[2].mV))
	setAndExport(ch, Cell4Voltage, float64(response.cells[3].mV))
	setAndExport(ch, Cell5Voltage, float64(response.cells[4].mV))
	setAndExport(ch, Cell6Voltage, float64(response.cells[5].mV))
	setAndExport(ch, Cell7Voltage, float64(response.cells[6].mV))
	setAndExport(ch, Cell8Voltage, float64(response.cells[7].mV))

	setAndExport(ch, Cell1Balancing, boolToFloat(response.cells[0].isBalancing))
	setAndExport(ch, Cell2Balancing, boolToFloat(response.cells[1].isBalancing))
	setAndExport(ch, Cell3Balancing, boolToFloat(response.cells[2].isBalancing))
	setAndExport(ch, Cell4Balancing, boolToFloat(response.cells[3].isBalancing))
	setAndExport(ch, Cell5Balancing, boolToFloat(response.cells[4].isBalancing))
	setAndExport(ch, Cell6Balancing, boolToFloat(response.cells[5].isBalancing))
	setAndExport(ch, Cell7Balancing, boolToFloat(response.cells[6].isBalancing))
	setAndExport(ch, Cell8Balancing, boolToFloat(response.cells[7].isBalancing))

	setAndExport(ch, TempInt, response.internalTemperature)
	setAndExport(ch, TempExt, response.externalTemperature)

	setAndExport(ch, Soc, response.soc)

	setAndExport(ch, batteryEnergyWh, response.batteryEnergyWh)
	setAndExport(ch, batteryEnergyAh, response.batteryEnergyAh)
	setAndExport(ch, pV1EnergyWh, response.pV1EnergyWh)
	setAndExport(ch, pV1EnergyAh, response.pV1EnergyAh)
	setAndExport(ch, pV2EnergyWh, response.pV2EnergyWh)
	setAndExport(ch, pV2EnergyAh, response.pV2EnergyAh)
	setAndExport(ch, loadEnergyWh, response.loadEnergyWh)
	setAndExport(ch, loadEnergyAh, response.loadEnergyAh)
	setAndExport(ch, extLoadEnergyWh, response.extLoadEnergyWh)
	setAndExport(ch, extLoadEnergyAh, response.extLoadEnergyAh)
	setAndExport(ch, dmpptEnergyWh, response.dmpptEnergyWh)
	setAndExport(ch, dmpptEnergyAh, response.dmpptEnergyAh)

	setAndExport(ch, cellType, response.cellType)
	setAndExport(ch, capacity, response.capacity)
	setAndExport(ch, status, response.status)

	setAndExport(ch, BatteryCurrent, response.batteryCurrent)
	setAndExport(ch, BatteryPower, response.batteryPower)
	setAndExport(ch, BatteryVoltage, response.batteryVoltage)
	setAndExport(ch, CurrentPv1, response.pv1Current)
	setAndExport(ch, CurrentPv2, response.pv2Current)
	setAndExport(ch, CurrentExtLoad, response.externalCurrent)

	setAndExport(ch, Ad2, float64(response.adc2))
	setAndExport(ch, Ad3, float64(response.adc3))
	setAndExport(ch, Ad4, float64(response.adc4))
	setAndExport(ch, Heat1, float64(response.heat1))
	setAndExport(ch, Heat2, float64(response.heat2))

	setAndExport(ch, Dfet, boolToFloat(response.flags.DischargeFETActive))
	setAndExport(ch, Eoc, boolToFloat(response.flags.EndOfCharge))
	setAndExport(ch, Cfet, boolToFloat(response.flags.ChargeFETActive))
	setAndExport(ch, Eccf, boolToFloat(response.flags.EEPROMFail))
	setAndExport(ch, Lvc, boolToFloat(response.flags.LowVoltageCell))
	setAndExport(ch, Open, boolToFloat(response.flags.OpenCellWire))
	setAndExport(ch, Celf, boolToFloat(response.flags.CellFail))
	setAndExport(ch, Dsc, boolToFloat(response.flags.DischargeShortCircuit))
	setAndExport(ch, Doc, boolToFloat(response.flags.DischargeOverCurrent))
	setAndExport(ch, Coc, boolToFloat(response.flags.ChargeOverCurrent))
	setAndExport(ch, Iot, boolToFloat(response.flags.InternalOverTemperature))
	setAndExport(ch, Uvlk, boolToFloat(response.flags.UnderVoltageLock))
	setAndExport(ch, Uv, boolToFloat(response.flags.UnderVoltage))
	setAndExport(ch, Ovlk, boolToFloat(response.flags.OverVoltageLock))
	setAndExport(ch, Ov, boolToFloat(response.flags.OverVoltage))
}

func (cc SBMS0SystemCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get(cc.url)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("debug resp is\n%s\n", string(b))
	data := decodeDebugResponse(b)

	// reset all the labels from last time
	systemTaskPriority.Reset()
	systemTaskRunTime.Reset()
	systemTaskRunTimePercent.Reset()
	systemTaskState.Reset()

	for _, d := range data {
		setAndExport(ch, systemTaskPriority.WithLabelValues(d.name), d.priority)
		setAndExport(ch, systemTaskRunTime.WithLabelValues(d.name), d.runTimeCounter)
		setAndExport(ch, systemTaskRunTimePercent.WithLabelValues(d.name), d.runTimePercent)
		setAndExport(ch, systemTaskState.WithLabelValues(d.name), d.state)
	}
}

func taskStateToValue(state string) float64 {
	// https://github.com/armageddon421/electrodacus-esp32/blob/1d8f5ec6a86613d09cdd9cef91cffe9b9a56a0bd/src/main.cpp#L697C1-L698C1
	states := map[string]float64{
		"RUN": 0,
		"RDY": 1,
		"BLK": 2,
		"SUS": 3,
		"DEL": 4,
	}

	value, ok := states[state]
	if !ok {
		return -1
	}
	return value
}

func setAndExport(ch chan<- prometheus.Metric, gauge prometheus.Gauge, value float64) {
	gauge.Set(value)
	ch <- gauge
}

func parseRawURL(raw string) (*url.URL, error) {
	u, err := url.ParseRequestURI(raw)
	if err != nil || u.Host == "" {
		u, repErr := url.ParseRequestURI("http://" + raw)
		if repErr != nil {
			fmt.Printf("Could not parse raw url: %s, error: %v", raw, err)
			return nil, err
		}
		return u, nil
	}
	return u, nil
}

func setDefaultPath(src, defaultPath string) (*url.URL, error) {
	u, err := parseRawURL(src)
	if err != nil {
		return nil, err
	}

	if u.Path == "" {
		u.Path = defaultPath
	}

	return u, nil
}

func getURL(src string) (string, error) {
	p, err := setDefaultPath(src, "/rawData")
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

func getDebugURL(src string) (string, error) {
	p, err := setDefaultPath(src, "/debug")
	if err != nil {
		return "", err
	}
	return p.String(), nil
}

func shouldEnableDefaultCollectors() bool {
	v := os.Getenv("ENABLE_DEFAULT_COLLECTORS")
	if slices.Contains([]string{"t", "1", "true", "yes"}, strings.ToLower(strings.TrimSpace(v))) {
		return true
	}
	return false
}

func main() {
	log.Println("starting")

	u, err := getURL(os.Getenv("URL"))
	debugURL, err := getDebugURL(os.Getenv("URL"))
	if err != nil {
		log.Fatal(err)
	}

	reg := prometheus.NewPedanticRegistry()
	systemMetricsReg := prometheus.NewPedanticRegistry()

	if shouldEnableDefaultCollectors() {
		reg.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)
	}

	reg.MustRegister(SBMS0Collector{url: u})
	systemMetricsReg.MustRegister(SBMS0SystemCollector{url: debugURL})

	handler := promhttp.InstrumentMetricHandler(reg, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	systemMetricsHandler := promhttp.InstrumentMetricHandler(systemMetricsReg, promhttp.HandlerFor(systemMetricsReg, promhttp.HandlerOpts{}))

	http.Handle("/metrics", handler)
	http.Handle("/metrics_system", systemMetricsHandler)
	log.Fatal(http.ListenAndServe(":9000", nil))
}
