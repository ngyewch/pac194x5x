package pac194x5x

import (
	"encoding/binary"
	"fmt"
	"math"
	"periph.io/x/conn/v3/i2c"
	"time"
)

const (
	// DefaultDelay is the default delay after a refresh.
	DefaultDelay = 128 * time.Millisecond
)

// PAC194x5x power monitor.
type PAC194x5x struct {
	i2cDev       *i2c.Dev
	voltageRatio []float64
	rSense       []float64
	productID    ProductID
	isPAC5x      bool
	channelCount int
}

// NewI2C initializes a power monitor through I2C connection.
func NewI2C(b i2c.Bus, addr uint16, voltageRatio []float64, rSense []float64) (*PAC194x5x, error) {
	d := &PAC194x5x{
		i2cDev: &i2c.Dev{
			Addr: addr,
			Bus:  b,
		},
		voltageRatio: voltageRatio,
		rSense:       rSense,
	}

	productID, err := d.GetProductID()
	if err != nil {
		return nil, err
	}

	d.productID = productID

	switch productID {
	case PAC1941:
		d.isPAC5x = false
		d.channelCount = 1
	case PAC1942_1:
		d.isPAC5x = false
		d.channelCount = 2
	case PAC1943:
		d.isPAC5x = false
		d.channelCount = 3
	case PAC1944:
		d.isPAC5x = false
		d.channelCount = 4
	case PAC1941_2:
		d.isPAC5x = false
		d.channelCount = 1
	case PAC1942_2:
		d.isPAC5x = false
		d.channelCount = 2
	case PAC1951:
		d.isPAC5x = true
		d.channelCount = 1
	case PAC1952_1:
		d.isPAC5x = true
		d.channelCount = 2
	case PAC1953:
		d.isPAC5x = true
		d.channelCount = 3
	case PAC1954:
		d.isPAC5x = true
		d.channelCount = 4
	case PAC1951_2:
		d.isPAC5x = true
		d.channelCount = 1
	case PAC1952_2:
		d.isPAC5x = true
		d.channelCount = 2
	default:
		return nil, fmt.Errorf("unknown product id: %d", productID)
	}

	return d, nil
}

// Channels returns the number of available channels.
func (d *PAC194x5x) Channels() int {
	return d.channelCount
}

// GetCtrl returns the Ctrl register value.
func (d *PAC194x5x) GetCtrl() (uint16, error) {
	return CtrlCacheRegister.Read(d)
}

// SetCtrl sets the Ctrl register value.
func (d *PAC194x5x) SetCtrl(v uint16) error {
	return CtrlCacheRegister.Write(d, v)
}

// GetAccCount returns the Acc_Count register value.
func (d *PAC194x5x) GetAccCount() (uint32, error) {
	return AccCountCacheRegister.Read(d)
}

// GetVAcc returns the Vacc_N register real data converted to W or V.
func (d *PAC194x5x) GetVAcc(channelNo int) (float64, UnitType, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, Unknown, err
	}

	accumConfig, err := d.getChAccumConfig(channelNo)
	if err != nil {
		return 0, Unknown, err
	}

	bidir := false
	var unit float64
	var unitType UnitType
	switch accumConfig {
	case 0:
		bidirV, _, err := d.getBidirFsrVLat(channelNo)
		if err != nil {
			return 0, Unknown, err
		}
		bidirI, _, err := d.getBidirFsrILat(channelNo)
		if err != nil {
			return 0, Unknown, err
		}
		bidir = bidirV || bidirI
		unit, err = d.getPowerUnit(channelNo)
		unitType = Watts
	case 1:
		bidir, _, err = d.getBidirFsrILat(channelNo)
		unit, err = d.getVSenseLSB(channelNo)
		unitType = Volts
	case 2:
		bidir, _, err = d.getBidirFsrVLat(channelNo)
		unit, err = d.getVBusLSB(channelNo)
		unitType = Volts
	case 3:
		bidir = false
		unit = 0
		unitType = Unknown
	}

	registers := []*CacheRegister[uint64]{VAcc1CacheRegister, VAcc2CacheRegister, VAcc3CacheRegister, VAcc4CacheRegister}
	v, err := registers[channelNo].Read(d)
	if err != nil {
		return 0, Unknown, err
	}

	raw := float64(v)
	if bidir {
		b := binary.BigEndian.AppendUint64(nil, v)
		if (b[1] & 0x80) == 0x80 {
			b[0] = 0xff
		}
		v = binary.BigEndian.Uint64(b)
		raw = float64(int64(v))
	}

	return raw * unit / d.voltageRatio[channelNo], unitType, nil
}

// GetVBus returns the Vbus_N register real data converted to V.
func (d *PAC194x5x) GetVBus(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := d.getVBusLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := d.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VBus1CacheRegister, VBus2CacheRegister, VBus3CacheRegister, VBus4CacheRegister}
	v, err := registers[channelNo].Read(d)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int16(v))
	}

	return raw * lsb / d.voltageRatio[channelNo], nil
}

// GetVSense returns the Vsense_N register real data converted to mV.
func (d *PAC194x5x) GetVSense(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := d.getVSenseLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := d.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VSense1CacheRegister, VSense2CacheRegister, VSense3CacheRegister, VSense4CacheRegister}
	v, err := registers[channelNo].Read(d)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int16(v))
	}

	return raw * lsb, nil
}

// GetCurrent calculates the Current value using the Vsense_N register and the Rsense_N resistor value, reported in mA.
func (d *PAC194x5x) GetCurrent(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	v, err := d.GetVSense(channelNo)
	if err != nil {
		return 0, err
	}

	return v / d.rSense[channelNo], nil
}

// GetVBusAvg returns the Vbus_Avg_N register real data converted to V.
func (d *PAC194x5x) GetVBusAvg(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := d.getVBusLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := d.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VBus1CacheRegister, VBus2CacheRegister, VBus3CacheRegister, VBus4CacheRegister}
	v, err := registers[channelNo].Read(d)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int16(v))
	}

	return raw * lsb / d.voltageRatio[channelNo], nil
}

// GetVSenseAvg returns the Vsense_Avg_N register real data converted to mV.
func (d *PAC194x5x) GetVSenseAvg(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := d.getVSenseLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := d.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VSense1AvgCacheRegister, VSense2AvgCacheRegister, VSense3AvgCacheRegister, VSense4AvgCacheRegister}
	v, err := registers[channelNo].Read(d)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int16(v))
	}

	return raw * lsb, nil
}

// GetCurrentAvg calculates the Current_Avg value using the Vsense_Avg_N register and the Rsense_N resistor value, reported in mA.
func (d *PAC194x5x) GetCurrentAvg(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	v, err := d.GetVSenseAvg(channelNo)
	if err != nil {
		return 0, err
	}

	return v / d.rSense[channelNo], nil
}

// GetVPower gets the Vpower_N register real data converted to W.
func (d *PAC194x5x) GetVPower(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := d.getPowerUnit(channelNo)
	if err != nil {
		return 0, err
	}

	bidirV, _, err := d.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	bidirI, _, err := d.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	bidir := bidirV || bidirI

	registers := []*CacheRegister[uint32]{VPower1CacheRegister, VPower2CacheRegister, VPower3CacheRegister, VPower4CacheRegister}
	v, err := registers[channelNo].Read(d)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int32(v))
	}
	raw /= 4

	return (raw * lsb) / d.voltageRatio[channelNo], nil
}

// GetEnergy calculates the Energy_N value (µWh) using the Vacc_N register real value.
func (d *PAC194x5x) GetEnergy(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	accumConfig, err := d.getChAccumConfig(channelNo)
	if err != nil {
		return 0, err
	}

	if accumConfig != 0 {
		return math.NaN(), nil
	}

	sampleMode, err := d.getSampleMode()
	if err != nil {
		return 0, err
	}

	singleShot := false
	if (sampleMode == SampleModeSingleShot) || (sampleMode == SampleModeSingleShot8x) {
		singleShot = true
	}

	if singleShot {
		// TODO handle single shot
		return math.NaN(), nil
	}

	v, _, err := d.GetVAcc(channelNo)
	if err != nil {
		return 0, err
	}

	sampleFrequency, err := d.getSampleFrequency()
	if err != nil {
		return 0, err
	}

	return (v / sampleFrequency) * (1000000 / 3600), nil
}

// GetNegPwrFsr returns the Neg_Pwr_Fsr register value.
func (d *PAC194x5x) GetNegPwrFsr() (uint16, error) {
	return NegPwrFsrCacheRegister.Read(d)
}

// SetNegPwrFsr sets the Neg_Pwr_Fsr register value.
func (d *PAC194x5x) SetNegPwrFsr(v uint16) error {
	return NegPwrFsrCacheRegister.Write(d, v)
}

// GetCtrlAct returns the Ctrl_Act register value.
func (d *PAC194x5x) GetCtrlAct() (uint16, error) {
	return CtrlActCacheRegister.Read(d)
}

// GetNegPwrFsrAct returns the Neg_Pwr_Fsr_Act register value.
func (d *PAC194x5x) GetNegPwrFsrAct() (uint16, error) {
	return NegPwrFsrActCacheRegister.Read(d)
}

// GetCtrlLat returns the Ctrl_Lat register value.
func (d *PAC194x5x) GetCtrlLat() (uint16, error) {
	return CtrlLatCacheRegister.Read(d)
}

// GetNegPwrFsrLat returns the Neg_Pwr_Fsr_Lat register value.
func (d *PAC194x5x) GetNegPwrFsrLat() (uint16, error) {
	return NegPwrFsrLatCacheRegister.Read(d)
}

// GetAccumConfig returns the Accum_Config register value.
func (d *PAC194x5x) GetAccumConfig() (uint8, error) {
	return AccumConfigCacheRegister.Read(d)
}

// SetAccumConfig sets the Accum_Config register value.
func (d *PAC194x5x) SetAccumConfig(v uint8) error {
	return AccumConfigCacheRegister.Write(d, v)
}

// Refresh sends a simple Refresh command to the device.
func (d *PAC194x5x) Refresh(delay time.Duration) error {
	for _, cached := range cacheRegisters {
		cached.Invalidate()
	}

	err := d.WriteRegister(RefreshRegister.Address, nil)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return nil
}

// RefreshG sends a Refresh_G command to the device.
func (d *PAC194x5x) RefreshG(delay time.Duration) error {
	for _, cached := range cacheRegisters {
		cached.Invalidate()
	}

	err := d.WriteRegister(RefreshGRegister.Address, nil)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return nil
}

// RefreshV sends a Refresh_V command to the device.
func (d *PAC194x5x) RefreshV(delay time.Duration) error {
	for _, cached := range cacheRegisters {
		cached.Invalidate()
	}

	err := d.WriteRegister(RefreshVRegister.Address, nil)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return nil
}

// GetAccumConfigAct returns the Accum_Config_Act register value.
func (d *PAC194x5x) GetAccumConfigAct() (uint8, error) {
	return AccumConfigActCacheRegister.Read(d)
}

// GetAccumConfigLat returns the Accum_Config_Lat register value.
func (d *PAC194x5x) GetAccumConfigLat() (uint8, error) {
	return AccumConfigLatCacheRegister.Read(d)
}

// GetProductID returns the product ID.
func (d *PAC194x5x) GetProductID() (ProductID, error) {
	return ProductIDCacheRegister.Read(d)
}

// GetManufacturerID returns the manufacturer ID.
func (d *PAC194x5x) GetManufacturerID() (uint8, error) {
	return ManufacturerIDCacheRegister.Read(d)
}

// GetRevisionID returns the revision ID.
func (d *PAC194x5x) GetRevisionID() (uint8, error) {
	return RevisionIDCacheRegister.Read(d)
}

// ReadRegister reads the register value.
func (d *PAC194x5x) ReadRegister(address uint8, len int) ([]byte, error) {
	readBytes := make([]byte, len)
	err := d.i2cDev.Tx([]byte{address}, readBytes)
	if err != nil {
		return nil, err
	}
	return readBytes, nil
}

// WriteRegister writes the value to the register.
func (d *PAC194x5x) WriteRegister(address uint8, data []byte) error {
	var writeBytes = []byte{address}
	writeBytes = append(writeBytes, data...)
	return d.i2cDev.Tx([]byte{address}, nil)
}

func (d *PAC194x5x) checkChannelNo(channelNo int) error {
	if (channelNo < 0) || (channelNo > d.channelCount) {
		return fmt.Errorf("invalid channel no: %d", channelNo)
	}
	return nil
}

func (d *PAC194x5x) getVBusLSB(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	var vBusScaleLat float64 = 9 // pac194x
	if d.isPAC5x {
		vBusScaleLat = 32
	}

	bidir, fsr, err := d.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	if bidir {
		vBusScaleLat *= 2
	}
	if !fsr {
		vBusScaleLat /= 2
	}

	return vBusScaleLat / 65536.0, nil
}

func (d *PAC194x5x) getVSenseLSB(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	vSenseScaleLat := 100.0

	bidir, fsr, err := d.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	if bidir {
		vSenseScaleLat *= 2
	}
	if !fsr {
		vSenseScaleLat /= 2
	}

	return vSenseScaleLat / 65536.0, nil
}

func (d *PAC194x5x) getPowerUnit(channelNo int) (float64, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	powerScaleLat := float64(100*9) / (d.rSense[channelNo] * 1000)
	if d.isPAC5x {
		powerScaleLat = float64(100*32) / (d.rSense[channelNo] * 1000)
	}

	bidirV, fsrV, err := d.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}
	bidirI, fsrI, err := d.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	if bidirV || bidirI {
		powerScaleLat *= 2
	}
	if !fsrV || !fsrI {
		powerScaleLat /= 2
	}

	return powerScaleLat / 1073741824.0, nil
}

func (d *PAC194x5x) getBidirFsrVLat(channelNo int) (bool, bool, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return false, false, err
	}

	negPwrFsrLat, err := d.GetNegPwrFsrLat()
	if err != nil {
		return false, false, err
	}

	bitsValue := getBitsValue(negPwrFsrLat, 2, 6-(channelNo*2))
	bidir := (bitsValue == 1) || (bitsValue == 2)
	fsr := bitsValue != 2

	return bidir, fsr, nil
}

func (d *PAC194x5x) getBidirFsrILat(channelNo int) (bool, bool, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return false, false, err
	}

	negPwrFsrLat, err := d.GetNegPwrFsrLat()
	if err != nil {
		return false, false, err
	}

	bitsValue := getBitsValue(negPwrFsrLat, 2, 14-(channelNo*2))
	bidir := (bitsValue == 1) || (bitsValue == 2)
	fsr := bitsValue != 2

	return bidir, fsr, nil
}

func (d *PAC194x5x) getChAccumConfig(channelNo int) (uint8, error) {
	err := d.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	accumConfig, err := d.GetAccumConfig()
	if err != nil {
		return 0, err
	}

	return getBitsValue(accumConfig, 2, 6-(channelNo*2)), nil
}

func (d *PAC194x5x) getSampleMode() (SampleMode, error) {
	v, err := d.GetCtrlLat()
	if err != nil {
		return 0, err
	}
	sampleMode := getBitsValue(v, 4, 12)
	return SampleMode(sampleMode), nil
}

func (d *PAC194x5x) getSampleFrequency() (float64, error) {
	sampleMode, err := d.getSampleMode()
	if err != nil {
		return 0, err
	}

	switch sampleMode {
	case SampleMode1024Adaptive:
		return 1024, nil
	case SampleMode256Adaptive:
		return 1024, nil
	case SampleMode64Adaptive:
		return 1024, nil
	case SampleMode8Adaptive:
		return 1024, nil
	case SampleMode1024:
		return 1024, nil
	case SampleMode256:
		return 256, nil
	case SampleMode64:
		return 64, nil
	case SampleMode8:
		return 8, nil
	case SampleModeSingleShot:
		return math.NaN(), nil
	case SampleModeSingleShot8x:
		return math.NaN(), nil
	case SampleModeFast, SampleModeBurst:
		v, err := d.GetCtrlLat()
		if err != nil {
			return 0, err
		}

		channelOnBits := getBitsValue(v, 4, 4)
		activeChannels := 0
		if (channelOnBits & 0x01) == 0 {
			activeChannels++
		}
		if (channelOnBits & 0x02) == 0 {
			activeChannels++
		}
		if (channelOnBits & 0x04) == 0 {
			activeChannels++
		}
		if (channelOnBits & 0x08) == 0 {
			activeChannels++
		}
		if activeChannels > d.channelCount {
			activeChannels = d.channelCount
		}

		if activeChannels == 0 {
			return math.NaN(), nil
		}

		return (1024 * 5) / float64(activeChannels), nil
	case SampleModeSleep:
		return math.NaN(), nil
	default:
		return math.NaN(), nil
	}
}

func getBitsValue[T uint8 | uint16](v T, bits int, position int) T {
	bitMask := (2 ^ bits) - 1
	return (v >> T(position)) & T(bitMask)
}
