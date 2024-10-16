package pac194x5x

var (
	AccCountCacheRegister       = NewCacheRegister[uint32](AccCountRegister, true)
	CtrlCacheRegister           = NewCacheRegister[uint16](CtrlRegister, true)
	VAcc1CacheRegister          = NewCacheRegister[uint64](VAcc1Register, true)
	VAcc2CacheRegister          = NewCacheRegister[uint64](VAcc2Register, true)
	VAcc3CacheRegister          = NewCacheRegister[uint64](VAcc3Register, true)
	VAcc4CacheRegister          = NewCacheRegister[uint64](VAcc4Register, true)
	VBus1CacheRegister          = NewCacheRegister[uint16](VBus1Register, true)
	VBus2CacheRegister          = NewCacheRegister[uint16](VBus2Register, true)
	VBus3CacheRegister          = NewCacheRegister[uint16](VBus3Register, true)
	VBus4CacheRegister          = NewCacheRegister[uint16](VBus4Register, true)
	VSense1CacheRegister        = NewCacheRegister[uint16](VSense1Register, true)
	VSense2CacheRegister        = NewCacheRegister[uint16](VSense2Register, true)
	VSense3CacheRegister        = NewCacheRegister[uint16](VSense3Register, true)
	VSense4CacheRegister        = NewCacheRegister[uint16](VSense4Register, true)
	VBus1AvgCacheRegister       = NewCacheRegister[uint16](VBus1AvgRegister, true)
	VBus2AvgCacheRegister       = NewCacheRegister[uint16](VBus2AvgRegister, true)
	VBus3AvgCacheRegister       = NewCacheRegister[uint16](VBus3AvgRegister, true)
	VBus4AvgCacheRegister       = NewCacheRegister[uint16](VBus4AvgRegister, true)
	VSense1AvgCacheRegister     = NewCacheRegister[uint16](VSense1AvgRegister, true)
	VSense2AvgCacheRegister     = NewCacheRegister[uint16](VSense2AvgRegister, true)
	VSense3AvgCacheRegister     = NewCacheRegister[uint16](VSense3AvgRegister, true)
	VSense4AvgCacheRegister     = NewCacheRegister[uint16](VSense4AvgRegister, true)
	VPower1CacheRegister        = NewCacheRegister[uint32](VPower1Register, true)
	VPower2CacheRegister        = NewCacheRegister[uint32](VPower2Register, true)
	VPower3CacheRegister        = NewCacheRegister[uint32](VPower3Register, true)
	VPower4CacheRegister        = NewCacheRegister[uint32](VPower4Register, true)
	SMBusCacheRegister          = NewCacheRegister[uint8](SMBusRegister, false)
	NegPwrFsrCacheRegister      = NewCacheRegister[uint16](NegPwrFsrRegister, true)
	CtrlActCacheRegister        = NewCacheRegister[uint16](CtrlActRegister, true)
	NegPwrFsrActCacheRegister   = NewCacheRegister[uint16](NegPwrFsrActRegister, true)
	CtrlLatCacheRegister        = NewCacheRegister[uint16](CtrlLatRegister, true)
	NegPwrFsrLatCacheRegister   = NewCacheRegister[uint16](NegPwrFsrLatRegister, true)
	AccumConfigCacheRegister    = NewCacheRegister[uint8](AccumConfigRegister, true)
	AccumConfigActCacheRegister = NewCacheRegister[uint8](AccumConfigActRegister, true)
	AccumConfigLatCacheRegister = NewCacheRegister[uint8](AccumConfigLatRegister, true)
	ProductIDCacheRegister      = NewCacheRegister[ProductID](ProductIDRegister, true)
	ManufacturerIDCacheRegister = NewCacheRegister[uint8](ManufacturerIDRegister, true)
	RevisionIDCacheRegister     = NewCacheRegister[uint8](RevisionIDRegister, true)

	cacheRegisters = []Cached{
		CtrlCacheRegister,
		AccCountCacheRegister,
		VAcc1CacheRegister,
		VAcc2CacheRegister,
		VAcc3CacheRegister,
		VAcc4CacheRegister,
		VBus1CacheRegister,
		VBus2CacheRegister,
		VBus3CacheRegister,
		VBus4CacheRegister,
		VSense1CacheRegister,
		VSense2CacheRegister,
		VSense3CacheRegister,
		VSense4CacheRegister,
		VBus1AvgCacheRegister,
		VBus2AvgCacheRegister,
		VBus3AvgCacheRegister,
		VBus4AvgCacheRegister,
		VSense1AvgCacheRegister,
		VSense2AvgCacheRegister,
		VSense3AvgCacheRegister,
		VSense4AvgCacheRegister,
		VPower1CacheRegister,
		VPower2CacheRegister,
		VPower3CacheRegister,
		VPower4CacheRegister,
		SMBusCacheRegister,
		NegPwrFsrCacheRegister,
		CtrlActCacheRegister,
		NegPwrFsrActCacheRegister,
		CtrlLatCacheRegister,
		NegPwrFsrLatCacheRegister,
		AccumConfigCacheRegister,
		AccumConfigActCacheRegister,
		AccumConfigLatCacheRegister,
		ProductIDCacheRegister,
		ManufacturerIDCacheRegister,
		RevisionIDCacheRegister,
	}
)

type Cached interface {
	IsValid() bool
	Invalidate()
}

type RegisterReader interface {
	ReadRegister(address uint8, len int) ([]byte, error)
}

type RegisterWriter interface {
	WriteRegister(address uint8, data []byte) error
}

type CacheRegister[T any] struct {
	register       Register[T]
	cachedRegister bool
	value          T
	valid          bool
}

func NewCacheRegister[T any](register Register[T], cachedRegister bool) *CacheRegister[T] {
	return &CacheRegister[T]{
		register:       register,
		cachedRegister: cachedRegister,
	}
}

func (cr *CacheRegister[T]) IsValid() bool {
	return cr.valid
}

func (cr *CacheRegister[T]) Invalidate() {
	cr.valid = false
}

func (cr *CacheRegister[T]) Read(reader RegisterReader) (T, error) {
	if cr.valid {
		return cr.value, nil
	}
	data, err := reader.ReadRegister(cr.register.Address, cr.register.Length)
	if err != nil {
		return *new(T), err
	}
	v, err := cr.register.Codec.Unmarshal(data)
	if err != nil {
		return *new(T), err
	}
	if cr.cachedRegister {
		cr.valid = true
		cr.value = v
	}
	return v, nil
}

func (cr *CacheRegister[T]) Write(writer RegisterWriter, v T) error {
	data, err := cr.register.Codec.Marshal(v)
	if err != nil {
		return err
	}
	err = writer.WriteRegister(cr.register.Address, data)
	if err != nil {
		return err
	}
	if cr.cachedRegister {
		cr.valid = true
		cr.value = v
	}
	return nil
}
