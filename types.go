package pac194x5x

// ProductID represents the product ID.
type ProductID uint8

const (
	PAC1941   ProductID = 104 // PAC1941 power monitor.
	PAC1942_1 ProductID = 105 // PAC1942_1 power monitor.
	PAC1943   ProductID = 106 // PAC1943 power monitor.
	PAC1944   ProductID = 107 // PAC1944 power monitor.
	PAC1941_2 ProductID = 108 // PAC1941_2 power monitor.
	PAC1942_2 ProductID = 109 // PAC1942_2 power monitor.
	PAC1951   ProductID = 120 // PAC1951 power monitor
	PAC1952_1 ProductID = 121 // PAC1952_1 power monitor
	PAC1953   ProductID = 122 // PAC1953 power monitor
	PAC1954   ProductID = 123 // PAC1954 power monitor
	PAC1951_2 ProductID = 124 // PAC1951_2 power monitor
	PAC1952_2 ProductID = 125 // PAC1952_2 power monitor
)

// UnitType represent the unit type.
type UnitType int

const (
	Unknown UnitType = iota // Unknown unit type.
	Volts                   // Volts unit type.
	Watts                   // Watts unit type.
)

// SampleMode represents the sample mode.
type SampleMode uint16

const (
	SampleMode1024Adaptive SampleMode = 0  // SampleMode1024Adaptive - 1024 SPS adaptive accumulation.
	SampleMode256Adaptive  SampleMode = 1  // SampleMode256Adaptive - 256 SPS adaptive accumulation.
	SampleMode64Adaptive   SampleMode = 2  // SampleMode64Adaptive - 64 SPS adaptive accumulation.
	SampleMode8Adaptive    SampleMode = 3  // SampleMode8Adaptive - 8 SPS adaptive accumulation.
	SampleMode1024         SampleMode = 4  // SampleMode1024 - 1024 SPS.
	SampleMode256          SampleMode = 5  // SampleMode256 - 256 SPS.
	SampleMode64           SampleMode = 6  // SampleMode64 - 64 SPS.
	SampleMode8            SampleMode = 7  // SampleMode8 - 8 SPS.
	SampleModeSingleShot   SampleMode = 8  // SampleModeSingleShot - Single-Shot mode.
	SampleModeSingleShot8x SampleMode = 9  // SampleModeSingleShot8x - Single-Shot 8x.
	SampleModeFast         SampleMode = 10 // SampleModeFast - Fast mode.
	SampleModeBurst        SampleMode = 11 // SampleModeBurst - Burst mode.
	SampleModeSleep        SampleMode = 15 //  SampleModeSleep - Sleep.
)
