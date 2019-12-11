package canal


type LeakyBuf struct {
	bufSize  int // size of each buffer
	freeList chan []byte
}

// NewLeakyBuf creates a leaky buffer which can hold at most n buffer, each
// with bufSize bytes.
func NewLeakyBuf(n, bufSize int) *LeakyBuf {
	return &LeakyBuf{
		bufSize:  bufSize,
		freeList: make(chan []byte, n),
	}
}

// Get returns a buffer from the leaky buffer or create a new buffer.
func (lb *LeakyBuf) Get() (b []byte) {
	select {
	case b = <-lb.freeList:
	default:
		b = make([]byte, lb.bufSize)
	}
	return
}

// Put add the buffer into the free buffer pool for reuse. Panic if the buffer
// size is not the same with the leaky buffer's. This is intended to expose
// error usage of leaky buffer.
func (lb *LeakyBuf) Put(b []byte) {
	if len(b) != lb.bufSize {
		panic("invalid buffer size that's put into leaky buffer")
	}
	select {
	case lb.freeList <- b[:lb.bufSize]:
	default:
	}
	return
}

type LeakyBufCommand struct {
	freeList chan *Command
}

func NewLeakyBufCommand(n int) *LeakyBufCommand {
	return &LeakyBufCommand{
		freeList: make(chan *Command, n),
	}
}

// Get returns a buffer from the leaky buffer or create a new buffer.
func (lb *LeakyBufCommand) Get() (b *Command) {
	select {
	case b = <-lb.freeList:
	default:
		b = &Command{}
	}
	return
}
func (lb *LeakyBufCommand) Put(b *Command) {
	select {
	case lb.freeList <- b:
	default:
	}
	return
}