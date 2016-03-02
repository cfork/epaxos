package state

import (
	"encoding/binary"
	"io"
	"sync"
)

func (t *Command) BinarySize() (nbytes int, sizeKnown bool) {
	return 17, true
}

type CommandCache struct {
	mu    sync.Mutex
	cache []*Command
}

func NewCommandCache() *CommandCache {
	c := &CommandCache{}
	c.cache = make([]*Command, 0)
	return c
}

func (p *CommandCache) Get() *Command {
	var t *Command
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Command{}
	}
	return t
}
func (p *CommandCache) Put(t *Command) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Command) Marshal(wire io.Writer) {
	var b [17]byte
	var bs []byte
	bs = b[:17]
	bs[0] = byte(t.Op)
	binary.LittleEndian.PutUint64(bs[1:9], uint64(t.K))
	binary.LittleEndian.PutUint64(bs[9:17], uint64(t.V))
	wire.Write(bs)
}

func (t *Command) Unmarshal(wire io.Reader) error {
	var b [17]byte
	var bs []byte
	bs = b[:17]
	if _, err := io.ReadAtLeast(wire, bs, 17); err != nil {
		return err
	}
	t.Op = Operation(bs[0])
	t.K = Key(binary.LittleEndian.Uint64(bs[1:9]))
	t.V = Value(binary.LittleEndian.Uint64(bs[9:17]))
	return nil
}

func (t *Key) Marshal(w io.Writer) {
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, uint64(*t))
	w.Write(bs)
}

func (t *Value) Marshal(w io.Writer) {
	var b [8]byte
	bs := b[:8]
	binary.LittleEndian.PutUint64(bs, uint64(*t))
	w.Write(bs)
}

func (t *Key) Unmarshal(r io.Reader) error {
	var b [8]byte
	bs := b[:8]
	if _, err := io.ReadFull(r, bs); err != nil {
		return err
	}
	*t = Key(binary.LittleEndian.Uint64(bs))
	return nil
}

func (t *Value) Unmarshal(r io.Reader) error {
	var b [8]byte
	bs := b[:8]
	if _, err := io.ReadFull(r, bs); err != nil {
		return err
	}
	*t = Value(binary.LittleEndian.Uint64(bs))
	return nil
}
