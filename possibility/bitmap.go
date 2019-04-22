package possibility

type BitMap interface {
	NumBits() uint
	And(other BitMap) BitMap
	Or(other BitMap) BitMap
	SetBit(bit uint)
	Set(val uint64)
	Contains(val uint) bool
	IsEmpty() bool
	Get() uint64
}

type PossBitMap struct {
	bitset  uint64
	numBits uint
}

func (bm *PossBitMap) NumBits() uint {
	return bm.numBits
}

func (bm *PossBitMap) And(other BitMap) BitMap {
	return GetBitMap(bm.Get()&other.Get(), bm.NumBits())
}

func (bm *PossBitMap) Or(other BitMap) BitMap {
	return GetBitMap(bm.Get()|other.Get(), bm.NumBits())
}

func (bm *PossBitMap) SetBit(bit uint) {
	if bit < 64 {
		currValue := bm.Get()
		currValue &^= 1 << bit
		bm.Set(currValue)
	}
}

func (bm *PossBitMap) Set(val uint64) {
	bm.bitset = val
}

func (bm *PossBitMap) Contains(num uint) bool {
	if num >= 64 {
		return false
	}
	return bm.Get()&(1<<num) != 0
}

func (bm *PossBitMap) IsEmpty() bool {
	return bm.Get() == 0
}

func (bm *PossBitMap) Get() uint64 {
	return bm.bitset
}

func GetBitMap(value uint64, numBits uint) BitMap {
	return &PossBitMap{
		bitset:  value,
		numBits: numBits,
	}
}
