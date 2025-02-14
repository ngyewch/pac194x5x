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

// Dev is a handle for a configured PAC194x5x device.
type Dev struct {
	i2cDev       *i2c.Dev
	voltageRatio []float64
	rSense       []float64
	productID    ProductID
	isPAC5x      bool
	channelCount int
}

// NewI2C initializes a power monitor through I2C connection.
func NewI2C(b i2c.Bus, addr uint16, voltageRatio []float64, rSense []float64) (*Dev, error) {
	dev := &Dev{
		i2cDev: &i2c.Dev{
			Addr: addr,
			Bus:  b,
		},
		voltageRatio: voltageRatio,
		rSense:       rSense,
	}

	productID, err := dev.GetProductID()
	if err != nil {
		return nil, err
	}

	dev.productID = productID

	switch productID {
	case PAC1941:
		dev.isPAC5x = false
		dev.channelCount = 1
	case PAC1942_1:
		dev.isPAC5x = false
		dev.channelCount = 2
	case PAC1943:
		dev.isPAC5x = false
		dev.channelCount = 3
	case PAC1944:
		dev.isPAC5x = false
		dev.channelCount = 4
	case PAC1941_2:
		dev.isPAC5x = false
		dev.channelCount = 1
	case PAC1942_2:
		dev.isPAC5x = false
		dev.channelCount = 2
	case PAC1951:
		dev.isPAC5x = true
		dev.channelCount = 1
	case PAC1952_1:
		dev.isPAC5x = true
		dev.channelCount = 2
	case PAC1953:
		dev.isPAC5x = true
		dev.channelCount = 3
	case PAC1954:
		dev.isPAC5x = true
		dev.channelCount = 4
	case PAC1951_2:
		dev.isPAC5x = true
		dev.channelCount = 1
	case PAC1952_2:
		dev.isPAC5x = true
		dev.channelCount = 2
	default:
		return nil, fmt.Errorf("unknown product id: %dev", productID)
	}

	return dev, nil
}

// Channels returns the number of available channels.
func (dev *Dev) Channels() int {
	return dev.channelCount
}

// GetCtrl returns the Ctrl register value.
func (dev *Dev) GetCtrl() (uint16, error) {
	return CtrlCacheRegister.Read(dev)
}

// SetCtrl sets the Ctrl register value.
func (dev *Dev) SetCtrl(v uint16) error {
	return CtrlCacheRegister.Write(dev, v)
}

// GetAccCount returns the Acc_Count register value.
func (dev *Dev) GetAccCount() (uint32, error) {
	return AccCountCacheRegister.Read(dev)
}

// GetVAcc returns the Vacc_N register real data converted to W or V.
func (dev *Dev) GetVAcc(channelNo int) (float64, UnitType, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, Unknown, err
	}

	accumConfig, err := dev.getChAccumConfig(channelNo)
	if err != nil {
		return 0, Unknown, err
	}

	bidir := false
	var unit float64
	var unitType UnitType
	switch accumConfig {
	case 0:
		bidirV, _, err := dev.getBidirFsrVLat(channelNo)
		if err != nil {
			return 0, Unknown, err
		}
		bidirI, _, err := dev.getBidirFsrILat(channelNo)
		if err != nil {
			return 0, Unknown, err
		}
		bidir = bidirV || bidirI
		unit, err = dev.getPowerUnit(channelNo)
		unitType = Watts
	case 1:
		bidir, _, err = dev.getBidirFsrILat(channelNo)
		unit, err = dev.getVSenseLSB(channelNo)
		unitType = Volts
	case 2:
		bidir, _, err = dev.getBidirFsrVLat(channelNo)
		unit, err = dev.getVBusLSB(channelNo)
		unitType = Volts
	case 3:
		bidir = false
		unit = 0
		unitType = Unknown
	}

	registers := []*CacheRegister[uint64]{VAcc1CacheRegister, VAcc2CacheRegister, VAcc3CacheRegister, VAcc4CacheRegister}
	v, err := registers[channelNo].Read(dev)
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

	return raw * unit / dev.voltageRatio[channelNo], unitType, nil
}

// GetVBus returns the Vbus_N register real data converted to V.
func (dev *Dev) GetVBus(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := dev.getVBusLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := dev.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VBus1CacheRegister, VBus2CacheRegister, VBus3CacheRegister, VBus4CacheRegister}
	v, err := registers[channelNo].Read(dev)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int16(v))
	}

	return raw * lsb / dev.voltageRatio[channelNo], nil
}

// GetVSense returns the Vsense_N register real data converted to mV.
func (dev *Dev) GetVSense(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := dev.getVSenseLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := dev.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VSense1CacheRegister, VSense2CacheRegister, VSense3CacheRegister, VSense4CacheRegister}
	v, err := registers[channelNo].Read(dev)
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
func (dev *Dev) GetCurrent(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	v, err := dev.GetVSense(channelNo)
	if err != nil {
		return 0, err
	}

	return v / dev.rSense[channelNo], nil
}

// GetVBusAvg returns the Vbus_Avg_N register real data converted to V.
func (dev *Dev) GetVBusAvg(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := dev.getVBusLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := dev.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VBus1AvgCacheRegister, VBus2AvgCacheRegister, VBus3AvgCacheRegister, VBus4AvgCacheRegister}
	v, err := registers[channelNo].Read(dev)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int16(v))
	}

	return raw * lsb / dev.voltageRatio[channelNo], nil
}

// GetVSenseAvg returns the Vsense_Avg_N register real data converted to mV.
func (dev *Dev) GetVSenseAvg(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := dev.getVSenseLSB(channelNo)
	if err != nil {
		return 0, err
	}

	bidir, _, err := dev.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	registers := []*CacheRegister[uint16]{VSense1AvgCacheRegister, VSense2AvgCacheRegister, VSense3AvgCacheRegister, VSense4AvgCacheRegister}
	v, err := registers[channelNo].Read(dev)
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
func (dev *Dev) GetCurrentAvg(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	v, err := dev.GetVSenseAvg(channelNo)
	if err != nil {
		return 0, err
	}

	return v / dev.rSense[channelNo], nil
}

// GetVPower gets the Vpower_N register real data converted to W.
func (dev *Dev) GetVPower(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	lsb, err := dev.getPowerUnit(channelNo)
	if err != nil {
		return 0, err
	}

	bidirV, _, err := dev.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}

	bidirI, _, err := dev.getBidirFsrILat(channelNo)
	if err != nil {
		return 0, err
	}

	bidir := bidirV || bidirI

	registers := []*CacheRegister[uint32]{VPower1CacheRegister, VPower2CacheRegister, VPower3CacheRegister, VPower4CacheRegister}
	v, err := registers[channelNo].Read(dev)
	if err != nil {
		return 0, err
	}

	raw := float64(v)
	if bidir {
		raw = float64(int32(v))
	}
	raw /= 4

	return (raw * lsb) / dev.voltageRatio[channelNo], nil
}

// GetEnergy calculates the Energy_N value (µWh) using the Vacc_N register real value.
func (dev *Dev) GetEnergy(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	accumConfig, err := dev.getChAccumConfig(channelNo)
	if err != nil {
		return 0, err
	}

	if accumConfig != 0 {
		return math.NaN(), nil
	}

	sampleMode, err := dev.getSampleMode()
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

	v, _, err := dev.GetVAcc(channelNo)
	if err != nil {
		return 0, err
	}

	sampleFrequency, err := dev.getSampleFrequency()
	if err != nil {
		return 0, err
	}

	return (v / sampleFrequency) * (1000000 / 3600), nil
}

// GetNegPwrFsr returns the Neg_Pwr_Fsr register value.
func (dev *Dev) GetNegPwrFsr() (uint16, error) {
	return NegPwrFsrCacheRegister.Read(dev)
}

// SetNegPwrFsr sets the Neg_Pwr_Fsr register value.
func (dev *Dev) SetNegPwrFsr(v uint16) error {
	return NegPwrFsrCacheRegister.Write(dev, v)
}

// GetCtrlAct returns the Ctrl_Act register value.
func (dev *Dev) GetCtrlAct() (uint16, error) {
	return CtrlActCacheRegister.Read(dev)
}

// GetNegPwrFsrAct returns the Neg_Pwr_Fsr_Act register value.
func (dev *Dev) GetNegPwrFsrAct() (uint16, error) {
	return NegPwrFsrActCacheRegister.Read(dev)
}

// GetCtrlLat returns the Ctrl_Lat register value.
func (dev *Dev) GetCtrlLat() (uint16, error) {
	return CtrlLatCacheRegister.Read(dev)
}

// GetNegPwrFsrLat returns the Neg_Pwr_Fsr_Lat register value.
func (dev *Dev) GetNegPwrFsrLat() (uint16, error) {
	return NegPwrFsrLatCacheRegister.Read(dev)
}

// GetAccumConfig returns the Accum_Config register value.
func (dev *Dev) GetAccumConfig() (uint8, error) {
	return AccumConfigCacheRegister.Read(dev)
}

// SetAccumConfig sets the Accum_Config register value.
func (dev *Dev) SetAccumConfig(v uint8) error {
	return AccumConfigCacheRegister.Write(dev, v)
}

// Refresh sends a simple Refresh command to the device.
func (dev *Dev) Refresh(delay time.Duration) error {
	for _, cached := range cacheRegisters {
		cached.Invalidate()
	}

	err := dev.WriteRegister(RefreshRegister.Address, nil)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return nil
}

// RefreshG sends a Refresh_G command to the device.
func (dev *Dev) RefreshG(delay time.Duration) error {
	for _, cached := range cacheRegisters {
		cached.Invalidate()
	}

	err := dev.WriteRegister(RefreshGRegister.Address, nil)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return nil
}

// RefreshV sends a Refresh_V command to the device.
func (dev *Dev) RefreshV(delay time.Duration) error {
	for _, cached := range cacheRegisters {
		cached.Invalidate()
	}

	err := dev.WriteRegister(RefreshVRegister.Address, nil)
	if err != nil {
		return err
	}
	time.Sleep(delay)
	return nil
}

// GetAccumConfigAct returns the Accum_Config_Act register value.
func (dev *Dev) GetAccumConfigAct() (uint8, error) {
	return AccumConfigActCacheRegister.Read(dev)
}

// GetAccumConfigLat returns the Accum_Config_Lat register value.
func (dev *Dev) GetAccumConfigLat() (uint8, error) {
	return AccumConfigLatCacheRegister.Read(dev)
}

// GetProductID returns the product ID.
func (dev *Dev) GetProductID() (ProductID, error) {
	return ProductIDCacheRegister.Read(dev)
}

// GetManufacturerID returns the manufacturer ID.
func (dev *Dev) GetManufacturerID() (uint8, error) {
	return ManufacturerIDCacheRegister.Read(dev)
}

// GetRevisionID returns the revision ID.
func (dev *Dev) GetRevisionID() (uint8, error) {
	return RevisionIDCacheRegister.Read(dev)
}

// ReadRegister reads the register value.
func (dev *Dev) ReadRegister(address uint8, len int) ([]byte, error) {
	readBytes := make([]byte, len)
	err := dev.i2cDev.Tx([]byte{address}, readBytes)
	if err != nil {
		return nil, err
	}
	return readBytes, nil
}

// WriteRegister writes the value to the register.
func (dev *Dev) WriteRegister(address uint8, data []byte) error {
	var writeBytes = []byte{address}
	writeBytes = append(writeBytes, data...)
	return dev.i2cDev.Tx([]byte{address}, nil)
}

func (dev *Dev) checkChannelNo(channelNo int) error {
	if (channelNo < 0) || (channelNo > dev.channelCount) {
		return fmt.Errorf("invalid channel no: %dev", channelNo)
	}
	return nil
}

func (dev *Dev) getVBusLSB(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	var vBusScaleLat float64 = 9 // pac194x
	if dev.isPAC5x {
		vBusScaleLat = 32
	}

	bidir, fsr, err := dev.getBidirFsrVLat(channelNo)
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

func (dev *Dev) getVSenseLSB(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	vSenseScaleLat := 100.0

	bidir, fsr, err := dev.getBidirFsrILat(channelNo)
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

func (dev *Dev) getPowerUnit(channelNo int) (float64, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	powerScaleLat := float64(100*9) / (dev.rSense[channelNo] * 1000)
	if dev.isPAC5x {
		powerScaleLat = float64(100*32) / (dev.rSense[channelNo] * 1000)
	}

	bidirV, fsrV, err := dev.getBidirFsrVLat(channelNo)
	if err != nil {
		return 0, err
	}
	bidirI, fsrI, err := dev.getBidirFsrILat(channelNo)
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

func (dev *Dev) getBidirFsrVLat(channelNo int) (bool, bool, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return false, false, err
	}

	negPwrFsrLat, err := dev.GetNegPwrFsrLat()
	if err != nil {
		return false, false, err
	}

	bitsValue := getBitsValue(negPwrFsrLat, 2, 6-(channelNo*2))
	bidir := (bitsValue == 1) || (bitsValue == 2)
	fsr := bitsValue != 2

	return bidir, fsr, nil
}

func (dev *Dev) getBidirFsrILat(channelNo int) (bool, bool, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return false, false, err
	}

	negPwrFsrLat, err := dev.GetNegPwrFsrLat()
	if err != nil {
		return false, false, err
	}

	bitsValue := getBitsValue(negPwrFsrLat, 2, 14-(channelNo*2))
	bidir := (bitsValue == 1) || (bitsValue == 2)
	fsr := bitsValue != 2

	return bidir, fsr, nil
}

func (dev *Dev) getChAccumConfig(channelNo int) (uint8, error) {
	err := dev.checkChannelNo(channelNo)
	if err != nil {
		return 0, err
	}

	accumConfig, err := dev.GetAccumConfig()
	if err != nil {
		return 0, err
	}

	return getBitsValue(accumConfig, 2, 6-(channelNo*2)), nil
}

func (dev *Dev) getSampleMode() (SampleMode, error) {
	v, err := dev.GetCtrlLat()
	if err != nil {
		return 0, err
	}
	sampleMode := getBitsValue(v, 4, 12)
	return SampleMode(sampleMode), nil
}

func (dev *Dev) getSampleFrequency() (float64, error) {
	sampleMode, err := dev.getSampleMode()
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
		v, err := dev.GetCtrlLat()
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
		if activeChannels > dev.channelCount {
			activeChannels = dev.channelCount
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
