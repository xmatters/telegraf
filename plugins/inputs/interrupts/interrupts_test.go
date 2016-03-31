package interrupts

import (
	"bytes"
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"math"
	"time"
)

func TestCountInterrupt(t *testing.T) {
	line := "8:         57         75   IO-APIC   8-edge      rtc0"
	label, counts := countInterrupt(line, 2)
	assert.Len(t, counts, 2)
	assert.Equal(t, uint64(57), counts[0])
	assert.Equal(t, uint64(75), counts[1])
	assert.Equal(t, int8, label)
}

func TestFirstGatherReturnsAllZeros(t *testing.T) {
	saved := readProcFile
	defer func() { readProcFile = saved }()
	readProcFile = newReadProcFileStub()

	i := &Interrupt{}
	acc := &testutil.Accumulator{}
	i.Gather(acc)
	acc.AssertContainsTaggedFields(t, inputName, expectedTaggedFields, map[string]string{"cpu": "cpu0"})
	acc.AssertContainsTaggedFields(t, inputName, expectedTaggedFields, map[string]string{"cpu": "cpu1"})
	acc.AssertContainsFields(t, inputName, expectedUntaggedFields)
}

func TestTenSecondInterval(t *testing.T) {
	savedNow := now
	savedRead := readProcFile
	defer func() {
		now = savedNow
		readProcFile = savedRead
	}()

	readProcFile = newReadProcFileStub()
	expectedCpu0 := make(map[string]interface{})
	expectedCpu1 := make(map[string]interface{})

	for k, counts := range expectedTenSecondIntervalByCpu {
		expectedCpu0[k] = counts[0]
		expectedCpu1[k] = counts[1]
	}

	i := &Interrupt{}
	acc := &testutil.Accumulator{}
	curTime := time.Now()

	now = func() time.Time { return curTime.Add(-10 * time.Second)}
	i.Gather(acc)

	acc = &testutil.Accumulator{}
	now = func() time.Time { return curTime }
	i.Gather(acc)

	acc.AssertContainsTaggedFields(t, inputName, expectedCpu0, map[string]string{"cpu": "cpu0"})
	acc.AssertContainsTaggedFields(t, inputName, expectedCpu1, map[string]string{"cpu": "cpu1"})
	acc.AssertContainsFields(t, inputName, map[string]interface{}{intERR: uint64(1), intMIS: uint64(1)})
}

func TestDerivative(t *testing.T) {
	type m struct {
		count     uint64
		prevCount uint64
		interval  time.Duration
		expected  uint64
	}

	data := []m{
		m{uint64(10), uint64(0), time.Second, uint64(10)},
		m{uint64(160), uint64(100), 60 * time.Second, uint64(1)},
		m{uint64(9462), uint64(0), 50 * time.Minute, uint64(3)},

		// Test rollovers
		m{uint64(1), math.MaxUint64 - uint64(10), time.Second, uint64(12)},
		m{uint64(0), math.MaxUint64, time.Second, uint64(1)},
	}

	now := time.Now()
	i := &Interrupt{
		prevMetrics: map[string][]uint64{
			"foo": []uint64{uint64(0)},
		},
		prevTime: now,
	}
	for _, dp := range data {
		i.prevMetrics["foo"][0] = dp.prevCount
		actual := i.derivative("foo", 0, dp.count, now.Add(dp.interval))
		assert.Equal(t, dp.expected, actual)
	}
}

func newReadProcFileStub() func(string) (*bytes.Buffer, error) {
	invocation := 0
	var testData = []string{testData0, testData1}

	return func(s string) (*bytes.Buffer, error) {
		data := bytes.NewBuffer([]byte(testData[invocation]))
		invocation++
		return data, nil
	}
}

const (
	int0   = "IO-APIC.2-edge.timer.0"
	int1   = "IO-APIC.1-edge.i8042.1"
	int8   = "IO-APIC.8-edge.rtc0.8"
	int9   = "IO-APIC.9-fasteoi.acpi.9"
	int12  = "IO-APIC.12-edge.i8042.12"
	int14  = "IO-APIC.14-edge.ata_piix.14"
	int15  = "IO-APIC.15-edge.ata_piix.15"
	int19  = "IO-APIC.19-fasteoi.enp0s3.19"
	int20  = "IO-APIC.20-fasteoi.vboxguest.20"
	int21  = "IO-APIC.21-fasteoi.0000:00:0d.0_.snd_intel8x0.21"
	int22  = "IO-APIC.22-fasteoi.ohci_hcd:usb1.22"
	intNMI = "Non-maskable.interrupts.NMI"
	intLOC = "Local.timer.interrupts.LOC"
	intSPU = "Spurious.interrupts.SPU"
	intPMI = "Performance.monitoring.interrupts.PMI"
	intIWI = "IRQ.work.interrupts.IWI"
	intRTR = "APIC.ICR.read.retries.RTR"
	intRES = "Rescheduling.interrupts.RES"
	intCAL = "Function.call.interrupts.CAL"
	intTLB = "TLB.shootdowns.TLB"
	intTRM = "Thermal.event.interrupts.TRM"
	intTHR = "Threshold.APIC.interrupts.THR"
	intDFR = "Deferred.Error.APIC.interrupts.DFR"
	intMCE = "Machine.check.exceptions.MCE"
	intMCP = "Machine.check.polls.MCP"
	intHYP = "Hypervisor.callback.interrupts.HYP"
	intPIN = "Posted-interrupt.notification.event.PIN"
	intPIW = "Posted-interrupt.wakeup.event.PIW"
	intERR = "ERR"
	intMIS = "MIS"
)

var expectedTaggedFields = map[string]interface{}{
	int0:   uint64(0),
	int1:   uint64(0),
	int8:   uint64(0),
	int9:   uint64(0),
	int12:  uint64(0),
	int14:  uint64(0),
	int15:  uint64(0),
	int19:  uint64(0),
	int20:  uint64(0),
	int21:  uint64(0),
	int22:  uint64(0),
	intNMI: uint64(0),
	intLOC: uint64(0),
	intSPU: uint64(0),
	intPMI: uint64(0),
	intIWI: uint64(0),
	intRTR: uint64(0),
	intRES: uint64(0),
	intCAL: uint64(0),
	intTLB: uint64(0),
	intTRM: uint64(0),
	intTHR: uint64(0),
	intDFR: uint64(0),
	intMCE: uint64(0),
	intMCP: uint64(0),
	intHYP: uint64(0),
	intPIN: uint64(0),
	intPIW: uint64(0),
}

var expectedTenSecondIntervalByCpu = map[string][2]interface{}{
int0:   [2]interface{}{uint64(10), uint64(1)},
int1:   [2]interface{}{uint64(10), uint64(10)},
int8:   [2]interface{}{uint64(1), uint64(1)},
int9:   [2]interface{}{uint64(2), uint64(1)},
int12:  [2]interface{}{uint64(3), uint64(3)},
int14:  [2]interface{}{uint64(1), uint64(1)},
int15:  [2]interface{}{uint64(1), uint64(1)},
int19:  [2]interface{}{uint64(10), uint64(1)},
int20:  [2]interface{}{uint64(10), uint64(10)},
int21:  [2]interface{}{uint64(10), uint64(20)},
int22:  [2]interface{}{uint64(1), uint64(2)},
intNMI: [2]interface{}{uint64(0), uint64(1)},
intLOC: [2]interface{}{uint64(1), uint64(1)},
intSPU: [2]interface{}{uint64(6), uint64(1)},
intPMI: [2]interface{}{uint64(1), uint64(1)},
intIWI: [2]interface{}{uint64(2), uint64(1)},
intRTR: [2]interface{}{uint64(100), uint64(100)},
intRES: [2]interface{}{uint64(100), uint64(100)},
intCAL: [2]interface{}{uint64(100), uint64(100)},
intTLB: [2]interface{}{uint64(1), uint64(10)},
intTRM: [2]interface{}{uint64(2), uint64(10)},
intTHR: [2]interface{}{uint64(1), uint64(1)},
intDFR: [2]interface{}{uint64(1), uint64(1)},
intMCE: [2]interface{}{uint64(1), uint64(802)},
intMCP: [2]interface{}{uint64(1), uint64(5)},
intHYP: [2]interface{}{uint64(1), uint64(1)},
intPIN: [2]interface{}{uint64(1), uint64(1)},
intPIW: [2]interface{}{uint64(1), uint64(1)},
}

var expectedUntaggedFields = map[string]interface{}{
	intERR: uint64(0),
	intMIS: uint64(0),
}

var testData0 = `          CPU0       CPU1
  0:         36         27   IO-APIC   2-edge      timer
  1:        123        321   IO-APIC   1-edge      i8042
  8:         57         75   IO-APIC   8-edge      rtc0
  9:          2          2   IO-APIC   9-fasteoi   acpi
 12:        204        402   IO-APIC  12-edge      i8042
 14:        71          17   IO-APIC  14-edge      ata_piix
 15:       1462       2641   IO-APIC  15-edge      ata_piix
 19:       2064       4602   IO-APIC  19-fasteoi   enp0s3
 20:      19828      82891   IO-APIC  20-fasteoi   vboxguest
 21:      23592      29532   IO-APIC  21-fasteoi   0000:00:0d.0, snd_intel8x0
 22:         29         92   IO-APIC  22-fasteoi   ohci_hcd:usb1
NMI:         13         31   Non-maskable interrupts
LOC:     125636     144093   Local timer interrupts
SPU:          7          7   Spurious interrupts
PMI:         67         76   Performance monitoring interrupts
IWI:        103        301   IRQ work interrupts
RTR:        987        789   APIC ICR read retries
RES:      34382      40965   Rescheduling interrupts
CAL:        270      12097   Function call interrupts
TLB:       1645       2104   TLB shootdowns
TRM:        123        321   Thermal event interrupts
THR:         45         54   Threshold APIC interrupts
DFR:          6          7   Deferred Error APIC interrupts
MCE:       8910       1980   Machine check exceptions
MCP:          5          5   Machine check polls
HYP:         24         42   Hypervisor callback interrupts
ERR:          3
MIS:          5
PIN:          7         11   Posted-interrupt notification event
PIW:          13        17   Posted-interrupt wakeup event
`

var testData1 = `          CPU0       CPU1
  0:        136         37   IO-APIC   2-edge      timer
  1:        223        421   IO-APIC   1-edge      i8042
  8:         67         85   IO-APIC   8-edge      rtc0
  9:         22         12   IO-APIC   9-fasteoi   acpi
 12:        234        432   IO-APIC  12-edge      i8042
 14:         81         27   IO-APIC  14-edge      ata_piix
 15:       1472       2651   IO-APIC  15-edge      ata_piix
 19:       2164       4612   IO-APIC  19-fasteoi   enp0s3
 20:      19928      82991   IO-APIC  20-fasteoi   vboxguest
 21:      23692      29732   IO-APIC  21-fasteoi   0000:00:0d.0, snd_intel8x0
 22:         38        112   IO-APIC  22-fasteoi   ohci_hcd:usb1
NMI:         14         41   Non-maskable interrupts
LOC:     125646     144103   Local timer interrupts
SPU:         71         17   Spurious interrupts
PMI:         77         86   Performance monitoring interrupts
IWI:        123        311   IRQ work interrupts
RTR:       1987       1789   APIC ICR read retries
RES:      35382      41965   Rescheduling interrupts
CAL:       1270      13097   Function call interrupts
TLB:       1655       2204   TLB shootdowns
TRM:        143        421   Thermal event interrupts
THR:         55         64   Threshold APIC interrupts
DFR:         16         17   Deferred Error APIC interrupts
MCE:       8920      10000   Machine check exceptions
MCP:         15         51   Machine check polls
HYP:         34         52   Hypervisor callback interrupts
ERR:         13
MIS:         15
PIN:         17         21   Posted-interrupt notification event
PIW:         23         27   Posted-interrupt wakeup event
`
