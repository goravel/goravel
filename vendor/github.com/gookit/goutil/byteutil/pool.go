package byteutil

// ChanPool struct
//
// Usage:
//
//	bp := strutil.NewByteChanPool(500, 1024, 1024)
//	buf:=bp.Get()
//	defer bp.Put(buf)
//	// use buf do something ...
//
// refer https://www.flysnow.org/2020/08/21/golang-chan-byte-pool.html
// from https://github.com/minio/minio/blob/master/internal/bpool/bpool.go
type ChanPool struct {
	c    chan []byte
	w    int // init byte width
	wcap int // set byte cap
}

// NewChanPool instance
func NewChanPool(chSize int, width int, capWidth int) *ChanPool {
	return &ChanPool{
		c:    make(chan []byte, chSize),
		w:    width,
		wcap: capWidth,
	}
}

// Get gets a []byte from the BytePool, or creates a new one if none are
// available in the pool.
func (bp *ChanPool) Get() (b []byte) {
	select {
	case b = <-bp.c: // reuse existing buffer
	default:
		// create new buffer
		if bp.wcap > 0 {
			b = make([]byte, bp.w, bp.wcap)
		} else {
			b = make([]byte, bp.w)
		}
	}
	return
}

// Put returns the given Buffer to the BytePool.
func (bp *ChanPool) Put(b []byte) {
	select {
	case bp.c <- b:
		// buffer went back into pool
	default:
		// buffer didn't go back into pool, just discard
	}
}

// Width returns the width of the byte arrays in this pool.
func (bp *ChanPool) Width() (n int) {
	return bp.w
}

// WidthCap returns the cap width of the byte arrays in this pool.
func (bp *ChanPool) WidthCap() (n int) {
	return bp.wcap
}
