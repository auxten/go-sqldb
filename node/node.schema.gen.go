package node

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

type Header struct {
	IsInternal bool
	IsRoot     bool
	Parent     uint32
}

func (d *Header) Size() (s uint64) {

	s += 6
	return
}
func (d *Header) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		if d.IsInternal {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
	}
	{
		if d.IsRoot {
			buf[1] = 1
		} else {
			buf[1] = 0
		}
	}
	{

		buf[0+2] = byte(d.Parent >> 0)

		buf[1+2] = byte(d.Parent >> 8)

		buf[2+2] = byte(d.Parent >> 16)

		buf[3+2] = byte(d.Parent >> 24)

	}
	return buf[:i+6], nil
}

func (d *Header) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		d.IsInternal = buf[0] == 1
	}
	{
		d.IsRoot = buf[1] == 1
	}
	{

		d.Parent = 0 | (uint32(buf[0+2]) << 0) | (uint32(buf[1+2]) << 8) | (uint32(buf[2+2]) << 16) | (uint32(buf[3+2]) << 24)

	}
	return i + 6, nil
}

type InternalNodeHeader struct {
	KeysNum    uint32
	RightChild uint32
}

func (d *InternalNodeHeader) Size() (s uint64) {

	s += 8
	return
}
func (d *InternalNodeHeader) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.KeysNum >> 0)

		buf[1+0] = byte(d.KeysNum >> 8)

		buf[2+0] = byte(d.KeysNum >> 16)

		buf[3+0] = byte(d.KeysNum >> 24)

	}
	{

		buf[0+4] = byte(d.RightChild >> 0)

		buf[1+4] = byte(d.RightChild >> 8)

		buf[2+4] = byte(d.RightChild >> 16)

		buf[3+4] = byte(d.RightChild >> 24)

	}
	return buf[:i+8], nil
}

func (d *InternalNodeHeader) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.KeysNum = 0 | (uint32(buf[0+0]) << 0) | (uint32(buf[1+0]) << 8) | (uint32(buf[2+0]) << 16) | (uint32(buf[3+0]) << 24)

	}
	{

		d.RightChild = 0 | (uint32(buf[0+4]) << 0) | (uint32(buf[1+4]) << 8) | (uint32(buf[2+4]) << 16) | (uint32(buf[3+4]) << 24)

	}
	return i + 8, nil
}

type LeafNodeHeader struct {
	Cells    uint32
	NextLeaf uint32
}

func (d *LeafNodeHeader) Size() (s uint64) {

	s += 8
	return
}
func (d *LeafNodeHeader) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Cells >> 0)

		buf[1+0] = byte(d.Cells >> 8)

		buf[2+0] = byte(d.Cells >> 16)

		buf[3+0] = byte(d.Cells >> 24)

	}
	{

		buf[0+4] = byte(d.NextLeaf >> 0)

		buf[1+4] = byte(d.NextLeaf >> 8)

		buf[2+4] = byte(d.NextLeaf >> 16)

		buf[3+4] = byte(d.NextLeaf >> 24)

	}
	return buf[:i+8], nil
}

func (d *LeafNodeHeader) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Cells = 0 | (uint32(buf[0+0]) << 0) | (uint32(buf[1+0]) << 8) | (uint32(buf[2+0]) << 16) | (uint32(buf[3+0]) << 24)

	}
	{

		d.NextLeaf = 0 | (uint32(buf[0+4]) << 0) | (uint32(buf[1+4]) << 8) | (uint32(buf[2+4]) << 16) | (uint32(buf[3+4]) << 24)

	}
	return i + 8, nil
}

type ICell struct {
	Key   uint32
	Child uint32
}

func (d *ICell) Size() (s uint64) {

	s += 8
	return
}
func (d *ICell) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Key >> 0)

		buf[1+0] = byte(d.Key >> 8)

		buf[2+0] = byte(d.Key >> 16)

		buf[3+0] = byte(d.Key >> 24)

	}
	{

		buf[0+4] = byte(d.Child >> 0)

		buf[1+4] = byte(d.Child >> 8)

		buf[2+4] = byte(d.Child >> 16)

		buf[3+4] = byte(d.Child >> 24)

	}
	return buf[:i+8], nil
}

func (d *ICell) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Key = 0 | (uint32(buf[0+0]) << 0) | (uint32(buf[1+0]) << 8) | (uint32(buf[2+0]) << 16) | (uint32(buf[3+0]) << 24)

	}
	{

		d.Child = 0 | (uint32(buf[0+4]) << 0) | (uint32(buf[1+4]) << 8) | (uint32(buf[2+4]) << 16) | (uint32(buf[3+4]) << 24)

	}
	return i + 8, nil
}

type InternalNode struct {
	CommonHeader Header
	Header       InternalNodeHeader
	ICells       [3]ICell
}

func (d *InternalNode) Size() (s uint64) {

	{
		s += d.CommonHeader.Size()
	}
	{
		s += d.Header.Size()
	}
	{
		for k := range d.ICells {
			_ = k // make compiler happy in case k is unused

			{
				s += d.ICells[k].Size()
			}

		}
	}
	return
}
func (d *InternalNode) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		nbuf, err := d.CommonHeader.Marshal(buf[0:])
		if err != nil {
			return nil, err
		}
		i += uint64(len(nbuf))
	}
	{
		nbuf, err := d.Header.Marshal(buf[i+0:])
		if err != nil {
			return nil, err
		}
		i += uint64(len(nbuf))
	}
	{
		for k := range d.ICells {

			{
				nbuf, err := d.ICells[k].Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		}
	}
	return buf[:i+0], nil
}

func (d *InternalNode) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		ni, err := d.CommonHeader.Unmarshal(buf[i+0:])
		if err != nil {
			return 0, err
		}
		i += ni
	}
	{
		ni, err := d.Header.Unmarshal(buf[i+0:])
		if err != nil {
			return 0, err
		}
		i += ni
	}
	{
		for k := range d.ICells {

			{
				ni, err := d.ICells[k].Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

		}
	}
	return i + 0, nil
}

type Cell struct {
	Key   uint32
	Value [296]byte
}

func (d *Cell) Size() (s uint64) {

	{
		s += 296
	}
	s += 4
	return
}
func (d *Cell) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Key >> 0)

		buf[1+0] = byte(d.Key >> 8)

		buf[2+0] = byte(d.Key >> 16)

		buf[3+0] = byte(d.Key >> 24)

	}
	{
		copy(buf[i+4:], d.Value[:])
		i += 296
	}
	return buf[:i+4], nil
}

func (d *Cell) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Key = 0 | (uint32(buf[i+0+0]) << 0) | (uint32(buf[i+1+0]) << 8) | (uint32(buf[i+2+0]) << 16) | (uint32(buf[i+3+0]) << 24)

	}
	{
		copy(d.Value[:], buf[i+4:])
		i += 296
	}
	return i + 4, nil
}

type LeafNode struct {
	CommonHeader Header
	Header       LeafNodeHeader
	Cells        [13]Cell
}

func (d *LeafNode) Size() (s uint64) {

	{
		s += d.CommonHeader.Size()
	}
	{
		s += d.Header.Size()
	}
	{
		for k := range d.Cells {
			_ = k // make compiler happy in case k is unused

			{
				s += d.Cells[k].Size()
			}

		}
	}
	return
}
func (d *LeafNode) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{
		nbuf, err := d.CommonHeader.Marshal(buf[0:])
		if err != nil {
			return nil, err
		}
		i += uint64(len(nbuf))
	}
	{
		nbuf, err := d.Header.Marshal(buf[i+0:])
		if err != nil {
			return nil, err
		}
		i += uint64(len(nbuf))
	}
	{
		for k := range d.Cells {

			{
				nbuf, err := d.Cells[k].Marshal(buf[i+0:])
				if err != nil {
					return nil, err
				}
				i += uint64(len(nbuf))
			}

		}
	}
	return buf[:i+0], nil
}

func (d *LeafNode) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{
		ni, err := d.CommonHeader.Unmarshal(buf[i+0:])
		if err != nil {
			return 0, err
		}
		i += ni
	}
	{
		ni, err := d.Header.Unmarshal(buf[i+0:])
		if err != nil {
			return 0, err
		}
		i += ni
	}
	{
		for k := range d.Cells {

			{
				ni, err := d.Cells[k].Unmarshal(buf[i+0:])
				if err != nil {
					return 0, err
				}
				i += ni
			}

		}
	}
	return i + 0, nil
}

type Row struct {
	Id       uint32
	Username [32]byte
	Email    [256]byte
}

func (d *Row) Size() (s uint64) {

	{
		s += 32
	}
	{
		s += 256
	}
	s += 4
	return
}
func (d *Row) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		buf[0+0] = byte(d.Id >> 0)

		buf[1+0] = byte(d.Id >> 8)

		buf[2+0] = byte(d.Id >> 16)

		buf[3+0] = byte(d.Id >> 24)

	}
	{
		copy(buf[i+4:], d.Username[:])
		i += 32
	}
	{
		copy(buf[i+4:], d.Email[:])
		i += 256
	}
	return buf[:i+4], nil
}

func (d *Row) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.Id = 0 | (uint32(buf[i+0+0]) << 0) | (uint32(buf[i+1+0]) << 8) | (uint32(buf[i+2+0]) << 16) | (uint32(buf[i+3+0]) << 24)

	}
	{
		copy(d.Username[:], buf[i+4:])
		i += 32
	}
	{
		copy(d.Email[:], buf[i+4:])
		i += 256
	}
	return i + 4, nil
}
