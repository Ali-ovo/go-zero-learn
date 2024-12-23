package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	if size == 0 {
		size = 250
	}

	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8,
	}
}

func (b *Bitmap) Set(id string) {
	// id index
	idx := hash(id) % b.size

	byteIdx := idx / 8

	bitIdx := idx % 8

	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	// id index
	idx := hash(id) % b.size

	byteIdx := idx / 8

	bitIdx := idx % 8

	return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}

func (b *Bitmap) Export() []byte {
	return b.bits
}

func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}

	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

func hash(id string) int {
	send := 131313
	hash := 0
	for _, c := range id {
		hash = hash*send + int(c)
	}

	return hash & 0x7FFFFFFF
}
