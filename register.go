package pac194x5x

// Register defines a PAC194x5x register.
type Register[T any] struct {
	Address uint8    // Address of register.
	Length  int      // Length in bytes.
	Codec   Codec[T] // Codec used to marshal/unmarshal values.
}

var (
	RefreshRegister           = Register[Void]{Address: 0x00, Length: 0, Codec: VoidCodec}     // RefreshRegister - REFRESH register.
	CtrlRegister              = Register[uint16]{Address: 0x01, Length: 2, Codec: Uint16Codec} // CtrlRegister - CTRL register.
	AccCountRegister          = Register[uint32]{Address: 0x02, Length: 4, Codec: Uint32Codec} // AccCountRegister - ACC_COUNT register.
	VAcc1Register             = Register[uint64]{Address: 0x03, Length: 7, Codec: Uint64Codec} // VAcc1Register - VACC1 register.
	VAcc2Register             = Register[uint64]{Address: 0x04, Length: 7, Codec: Uint64Codec} // VAcc2Register - VACC2 register.
	VAcc3Register             = Register[uint64]{Address: 0x05, Length: 7, Codec: Uint64Codec} // VAcc3Register - VACC3 register.
	VAcc4Register             = Register[uint64]{Address: 0x06, Length: 7, Codec: Uint64Codec} // VAcc4Register - VACC4 register.
	VBus1Register             = Register[uint16]{Address: 0x07, Length: 2, Codec: Uint16Codec} // VBus1Register - VBUS1 register.
	VBus2Register             = Register[uint16]{Address: 0x08, Length: 2, Codec: Uint16Codec} // VBus2Register - VBUS2 register.
	VBus3Register             = Register[uint16]{Address: 0x09, Length: 2, Codec: Uint16Codec} // VBus3Register - VBUS3 register.
	VBus4Register             = Register[uint16]{Address: 0x0a, Length: 2, Codec: Uint16Codec} // VBus4Register - VBUS4 register.
	VSense1Register           = Register[uint16]{Address: 0x0b, Length: 2, Codec: Uint16Codec} // VSense1Register - VSENSE1 register.
	VSense2Register           = Register[uint16]{Address: 0x0c, Length: 2, Codec: Uint16Codec} // VSense2Register - VSENSE2 register.
	VSense3Register           = Register[uint16]{Address: 0x0d, Length: 2, Codec: Uint16Codec} // VSense3Register - VSENSE3 register.
	VSense4Register           = Register[uint16]{Address: 0x0e, Length: 2, Codec: Uint16Codec} // VSense4Register - VSENSE4 register.
	VBus1AvgRegister          = Register[uint16]{Address: 0x0f, Length: 2, Codec: Uint16Codec} // VBus1AvgRegister - VBUS1_AVG register.
	VBus2AvgRegister          = Register[uint16]{Address: 0x10, Length: 2, Codec: Uint16Codec} // VBus2AvgRegister - VBUS2_AVG register.
	VBus3AvgRegister          = Register[uint16]{Address: 0x11, Length: 2, Codec: Uint16Codec} // VBus3AvgRegister - VBUS3_AVG register.
	VBus4AvgRegister          = Register[uint16]{Address: 0x12, Length: 2, Codec: Uint16Codec} // VBus4AvgRegister - VBUS4_AVG register.
	VSense1AvgRegister        = Register[uint16]{Address: 0x13, Length: 2, Codec: Uint16Codec} // VSense1AvgRegister - VSENSE1_AVG register.
	VSense2AvgRegister        = Register[uint16]{Address: 0x14, Length: 2, Codec: Uint16Codec} // VSense2AvgRegister - VSENSE2_AVG register.
	VSense3AvgRegister        = Register[uint16]{Address: 0x15, Length: 2, Codec: Uint16Codec} // VSense3AvgRegister - VSENSE3_AVG register.
	VSense4AvgRegister        = Register[uint16]{Address: 0x16, Length: 2, Codec: Uint16Codec} // VSense4AvgRegister - VSENSE4_AVG register.
	VPower1Register           = Register[uint32]{Address: 0x17, Length: 4, Codec: Uint32Codec} // VPower1Register - VPOWER1 register.
	VPower2Register           = Register[uint32]{Address: 0x18, Length: 4, Codec: Uint32Codec} // VPower2Register - VPOWER2 register.
	VPower3Register           = Register[uint32]{Address: 0x19, Length: 4, Codec: Uint32Codec} // VPower3Register - VPOWER3 register.
	VPower4Register           = Register[uint32]{Address: 0x1a, Length: 4, Codec: Uint32Codec} // VPower4Register - VPOWER4 register.
	SMBusRegister             = Register[uint8]{Address: 0x1c, Length: 1, Codec: Uint8Codec}   // SMBusRegister - SMBUS SETTINGS register.
	NegPwrFsrRegister         = Register[uint16]{Address: 0x1d, Length: 2, Codec: Uint16Codec} // NegPwrFsrRegister - NEG_PWR_FSR register.
	RefreshGRegister          = Register[Void]{Address: 0x1e, Length: 0, Codec: VoidCodec}     // RefreshGRegister - REFRESH_G register.
	RefreshVRegister          = Register[Void]{Address: 0x1f, Length: 0, Codec: VoidCodec}     // RefreshVRegister - REFRESH_V register.
	SlowRegister              = Register[any]{Address: 0x20, Length: 1}                        // SlowRegister - SLOW register.
	CtrlActRegister           = Register[uint16]{Address: 0x21, Length: 2, Codec: Uint16Codec} // CtrlActRegister - CTRL_ACT register.
	NegPwrFsrActRegister      = Register[uint16]{Address: 0x22, Length: 2, Codec: Uint16Codec} // NegPwrFsrActRegister - NEG_PWR_FSR_ACT register.
	CtrlLatRegister           = Register[uint16]{Address: 0x23, Length: 2, Codec: Uint16Codec} // CtrlLatRegister - CTRL_LAT register.
	NegPwrFsrLatRegister      = Register[uint16]{Address: 0x24, Length: 2, Codec: Uint16Codec} // NegPwrFsrLatRegister - NEG_PWR_FSR_LAT register.
	AccumConfigRegister       = Register[uint8]{Address: 0x25, Length: 1, Codec: Uint8Codec}   // AccumConfigRegister - ACCUM_CONFIG register.
	AlertStatusRegister       = Register[any]{Address: 0x26, Length: 3}
	SlowAlert1Register        = Register[any]{Address: 0x27, Length: 3}
	GPIOAlert2Register        = Register[any]{Address: 0x28, Length: 3}
	AccFullnessLimitsRegister = Register[any]{Address: 0x29, Length: 2}
	OCLimit1Register          = Register[any]{Address: 0x30, Length: 2}
	OCLimit2Register          = Register[any]{Address: 0x31, Length: 2}
	OCLimit3Register          = Register[any]{Address: 0x32, Length: 2}
	OCLimit4Register          = Register[any]{Address: 0x33, Length: 2}
	UCLimit1Register          = Register[any]{Address: 0x34, Length: 2}
	UCLimit2Register          = Register[any]{Address: 0x35, Length: 2}
	UCLimit3Register          = Register[any]{Address: 0x36, Length: 2}
	UCLimit4Register          = Register[any]{Address: 0x37, Length: 2}
	OPLimit1Register          = Register[any]{Address: 0x38, Length: 3}
	OPLimit2Register          = Register[any]{Address: 0x39, Length: 3}
	OPLimit3Register          = Register[any]{Address: 0x3a, Length: 3}
	OPLimit4Register          = Register[any]{Address: 0x3b, Length: 3}
	OVLimit1Register          = Register[any]{Address: 0x3c, Length: 2}
	OVLimit2Register          = Register[any]{Address: 0x3d, Length: 2}
	OVLimit3Register          = Register[any]{Address: 0x3e, Length: 2}
	OVLimit4Register          = Register[any]{Address: 0x3f, Length: 2}
	UVLimit1Register          = Register[any]{Address: 0x40, Length: 2}
	UVLimit2Register          = Register[any]{Address: 0x41, Length: 2}
	UVLimit3Register          = Register[any]{Address: 0x42, Length: 2}
	UVLimit4Register          = Register[any]{Address: 0x43, Length: 2}
	OCLimitNSamplesRegister   = Register[any]{Address: 0x44, Length: 1}
	UCLimitNSamplesRegister   = Register[any]{Address: 0x45, Length: 1}
	OPLimitNSamplesRegister   = Register[any]{Address: 0x46, Length: 1}
	OVLimitNSamplesRegister   = Register[any]{Address: 0x47, Length: 1}
	UVLimitNSamplesRegister   = Register[any]{Address: 0x48, Length: 1}
	AlertEnableRegister       = Register[any]{Address: 0x49, Length: 3}
	AccumConfigActRegister    = Register[uint8]{Address: 0x4a, Length: 1, Codec: Uint8Codec}
	AccumConfigLatRegister    = Register[uint8]{Address: 0x4b, Length: 1, Codec: Uint8Codec}
	ProductIDRegister         = Register[ProductID]{Address: 0xfd, Length: 1, Codec: ProductIDCodec}
	ManufacturerIDRegister    = Register[uint8]{Address: 0xfe, Length: 1, Codec: Uint8Codec}
	RevisionIDRegister        = Register[uint8]{Address: 0xff, Length: 1, Codec: Uint8Codec}
)
