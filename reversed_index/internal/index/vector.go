package index

type BitVector struct {
	RoaringBitmap
}

func (bv BitVector) And(vec BitVector) {
	bv.RoaringBitmap.And(vec.RoaringBitmap)
}

func (bv BitVector) AndNot(vec BitVector) {
	bv.RoaringBitmap.AndNot(vec.RoaringBitmap)
}

func (bv BitVector) Or(vec BitVector) {
	bv.RoaringBitmap.Or(vec.RoaringBitmap)
}

func (bv BitVector) Xor(vec BitVector) {
	bv.RoaringBitmap.Xor(vec.RoaringBitmap)
}
