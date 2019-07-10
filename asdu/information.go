package asdu

// about information object 应用服务数据单元 - 信息对象

// InfoObjAddr is the information object address.
// The width is controlled by Params.InfoObjAddrSize.
// See companion standard 101, subclause 7.2.5.
// - width 1
// <0>: 无关的信息对象地址
// <1..255>: 信息对象地址
// - width 2
// <0>: 无关的信息对象地址
// <1..65535>: 信息对象地址
// - width 3
// <0>: 无关的信息对象地址
// <1..16777215>: 信息对象地址
type InfoObjAddr uint

// InfoObjIrrelevantAddr Zero means that the information object address is irrelevant.
const InfoObjIrrelevantAddr InfoObjAddr = 0

// SinglePoint is a measured value of a switch.
// See companion standard 101, subclause 7.2.6.1.
type SinglePoint byte

// 单点信息
const (
	SPIOff SinglePoint = iota
	SPIOn
)

// DoublePoint is a measured value of a determination aware switch.
// See companion standard 101, subclause 7.2.6.2.
type DoublePoint byte

// 双点信息
const (
	DPIIndeterminateOrIntermediate DoublePoint = iota // 不确定或中间状态
	DPIDeterminedOff                                  // 确定状态开
	DPIDeterminedOn                                   // 确定状态关
	DPIIndeterminate                                  // 不确定或中间状态
)

func (this DoublePoint) Value() byte {
	return byte(this & 0x03)
}

// Quality descriptor flags attribute measured values.
// See companion standard 101, subclause 7.2.6.3.
const (
	// QDSOverflow marks whether the value is beyond a predefined range.
	QDSOverflow = 1 << iota

	_ // reserve
	_ // reserve

	// QDSTimeInvalid flags that the elapsed time was incorrectly acquired.
	// This attribute is only valid for events of protection equipment.
	// See companion standard 101, subclause 7.2.6.4.
	QDSTimeInvalid

	// QDSBlocked flags that the value is blocked for transmission; the
	// value remains in the state that was acquired before it was blocked.
	QDSBlocked

	// QDSSubstituted flags that the value was provided by the input of
	// an operator (dispatcher) instead of an automatic source.
	QDSSubstituted

	// QDSNotTopical flags that the most recent update was unsuccessful.
	QDSNotTopical

	// QDSInvalid flags that the value was incorrectly acquired.
	QDSInvalid

	// QDSOK means no flags, no problems.
	QDSOK = 0
)

// StepPos is a measured value with transient state indication.
// 带瞬变状态指示的测量值，用于变压器步位置或其它步位置的值
// See companion standard 101, subclause 7.2.6.5.
type StepPos int

// NewStepPos returns a new step position.
// Values range<-64..63>
// bit[0-6]: <-64..63>
// NOTE: bit6 为符号位
// bit7: 0: 设备未在瞬变状态 1： 设备处于瞬变状态
func NewStepPos(value int, hasTransient bool) StepPos {
	p := StepPos(value & 0x7f)
	if hasTransient {
		p |= 0x80
	}
	return p
}

// ToPos 返回 value in [-64, 63] 和 hasTransient 是否瞬变状态.
func (this StepPos) ToPos() (value int, hasTransient bool) {
	u := uint(this)
	if u&0x40 == 0 {
		value = int(u & 0x3f)
	} else {
		value = int(u) | (-1 &^ 0x3f)
	}
	hasTransient = (u & 0x80) != 0
	return
}

// Normalize is a 16-bit normalized value.
// 规一化值
// See companion standard 101, subclause 7.2.6.6.
type Normalize int16

// Float64 returns the value in [-1, 1 − 2⁻¹⁵].
func (this Normalize) Float64() float64 {
	return float64(this) / 32768
}

// Qualifier Of Parameter Of Measured Values
// 测量值参数限定词
// See companion standard 101, subclause 7.2.6.24.
const (
	_             = iota // 0: not used
	QPMThreashold        // 1: threshold value
	QPMSmoothing         // 2: smoothing factor (filter time constant)
	QPMLowLimit          // 3: low limit for transmission of measured values
	QPMHighLimit         // 4: high limit for transmission of measured values

	// 5‥31: reserved for standard definitions of this companion standard (compatible range)
	// 32‥63: reserved for special use (private range)

	QPMChangeFlag      = 0x40 // bit6 marks local parameter change  当地参数改变
	QPMInOperationFlag = 0x80 // bit7 marks parameter operation 参数在运行
)

// CmdQualifier is a qualifier of qual.
// See companion standard 101, subclause 7.2.6.26.
// <0>: 未用
//  the qualifier of command.
//	0: no additional definition
//	1: short pulse duration (circuit-breaker), duration determined by a system parameter in the outstation
//	2: long pulse duration, duration determined by a system parameter in the outstation
//	3: persistent output
//	4‥8: reserved for standard definitions of this companion standard
//	9‥15: reserved for the selection of other predefined functions
//	16‥31: reserved for special use (private range)
type CmdQualifier byte

// QualifierOfCmd is a  qualifier of command.
// 命令限定词
type QualifierOfCmd struct {
	CmdQ CmdQualifier
	// See section 5, subclause 6.8.
	// executes(false) (or selects(true)).
	InExec bool
}

func DecodeQualifierOfCmd(b byte) QualifierOfCmd {
	return QualifierOfCmd{
		CmdQ:   CmdQualifier((b >> 2) & 0x1f),
		InExec: b&0x80 == 0,
	}
}

func (this QualifierOfCmd) Value() byte {
	v := (byte(this.CmdQ) & 0x1f) << 2
	if !this.InExec {
		v |= 0x80
	}
	return v
}

// CmdSetPoint is the qualifier of a set-point command qual.
// See companion standard 101, subclause 7.2.6.39.
//	0: default
//	0‥63: reserved for standard definitions of this companion standard (compatible range)
//	64‥127: reserved for special use (private range)
type CmdSetPoint uint

// QualifierOfCmd is a  qualifier of command.
type QualifierOfSetpointCmd struct {
	CmdS CmdSetPoint
	// See section 5, subclause 6.8.
	// executes(false) (or selects(true)).
	InExec bool
}

func DecodeQualifierOfSetpointCmd(b byte) QualifierOfSetpointCmd {
	return QualifierOfSetpointCmd{
		CmdS:   CmdSetPoint(b & 0x7f),
		InExec: b&0x80 == 0,
	}
}

func (this QualifierOfSetpointCmd) Value() byte {
	v := byte(this.CmdS) & 0x7f
	if !this.InExec {
		v |= 0x80
	}
	return v
}