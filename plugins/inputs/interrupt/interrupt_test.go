package entropy

import (
	"bytes"
	"testing"
	//	"time"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
)

// TODO: use telegraf.testutil.Accumulator?
/*
type MockAccum struct {
	fields map[string]interface{}
	tags map[string]string
}

func (m *MockAccum) Add(measurement string, value interface{}, t ...time.Time) {}
func (m *MockAccum) AddFields(measurement string, tags map[string]string, t ...time.Time) {}
func (m *MockAccum) Debug() bool { return false}
func (m *MockAccum) SetDebug(enabled bool) {}
*/

func TestCountInterrupt(t *testing.T) {
	line := "8:         57         75   IO-APIC   8-edge      rtc0"
	label, counts := countInterrupt(line, 2)
	assert.Len(t, counts, 2)
	assert.Equal(t, uint64(57), counts[0])
	assert.Equal(t, uint64(75), counts[1])
	assert.Equal(t, "IO-APIC.8-edge.rtc0.8", label)

}

func TestFirstGatherReturnsAllZeros(t *testing.T) {
	invocation := 0
	testData := []string{testData0, testData1}
	i := &Interrupt{
		readProcFile: func(s string) (*bytes.Buffer, error) {
			data := bytes.NewBuffer([]byte(testData[invocation]))
			invocation++
			return data, nil
		},
	}

	acc := &testutil.Accumulator{}
	i.Gather(acc)
	acc.AssertContainsTaggedFields(t, inputName, expectedTaggedFields, map[string]string{"cpu": "cpu0"})
	acc.AssertContainsTaggedFields(t, inputName, expectedTaggedFields, map[string]string{"cpu": "cpu1"})
	acc.AssertContainsFields(t, inputName, expectedUntaggedFields)

}

var expectedTaggedFields = map[string]interface{}{
	"APIC.ICR.read.retries.RTR":                        uint64(0),
	"Deferred.Error.APIC.interrupts.DFR":               uint64(0),
	"Hypervisor.callback.interrupts.HYP":               uint64(0),
	"IO-APIC.1-edge.i8042.1":                           uint64(0),
	"IO-APIC.12-edge.i8042.12":                         uint64(0),
	"IO-APIC.14-edge.ata_piix.14":                      uint64(0),
	"IO-APIC.15-edge.ata_piix.15":                      uint64(0),
	"IO-APIC.19-fasteoi.enp0s3.19":                     uint64(0),
	"IO-APIC.2-edge.timer.0":                           uint64(0),
	"IO-APIC.20-fasteoi.vboxguest.20":                  uint64(0),
	"IO-APIC.21-fasteoi.0000:00:0d.0_.snd_intel8x0.21": uint64(0),
	"IO-APIC.22-fasteoi.ohci_hcd:usb1.22":              uint64(0),
	"IO-APIC.8-edge.rtc0.8":                            uint64(0),
	"IRQ.work.interrupts.IWI":                          uint64(0),
	"Machine.check.exceptions.MCE":                     uint64(0),
	"Non-maskable.interrupts.NMI":                      uint64(0),
	"Performance.monitoring.interrupts.PMI":            uint64(0),
	"Posted-interrupt.notification.event.PIN":          uint64(0),
	"Posted-interrupt.wakeup.event.PIW":                uint64(0),
	"Spurious.interrupts.SPU":                          uint64(0),
	"Thermal.event.interrupts.TRM":                     uint64(0),
	"Threshold.APIC.interrupts.THR":                    uint64(0),
	"Rescheduling.interrupts.RES":                      uint64(0),
	"Machine.check.polls.MCP":                          uint64(0),
	"IO-APIC.9-fasteoi.acpi.9":                         uint64(0),
	"Function.call.interrupts.CAL":                     uint64(0),
	"Local.timer.interrupts.LOC":                       uint64(0),
	"TLB.shootdowns.TLB":                               uint64(0),
}

var expectedUntaggedFields = map[string]interface{}{
	"ERR": uint64(0),
	"MIS": uint64(0),
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
