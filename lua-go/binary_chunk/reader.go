package binary_chunk

import (
	"encoding/binary"
	"math"
)

type Reader struct {
	data []byte
}

func New(data []byte) *Reader {
	return &Reader{data}
}

func (r *Reader) readByte() byte {
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *Reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(r.data)
	r.data = r.data[4:]
	return i
}

func (r *Reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(r.data)
	r.data = r.data[8:]
	return i
}

func (r *Reader) readLuaInteger() int64 {
	i := r.readUint64()
	return int64(i)
}

func (r *Reader) readLuaNumber() float64 {
	return math.Float64frombits(r.readUint64())
}

func (r *Reader) readBytes(n uint) []byte {
	bs := r.data[0:n]
	r.data = r.data[n:]
	return bs
}

func (r *Reader) readString() string {
	size := uint(r.readByte())

	if size == 0 {
		return ""
	}

	if size == 0xff {
		size = uint(r.readUint64())
	}

	bytes := r.readBytes(size - 1)
	return string(bytes)
}

func (r *Reader) checkHeader() {
	if string(r.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk")
	}
	if r.readByte() != LUAC_VERSION {
		panic("version not match")
	}
	if r.readByte() != LUAC_FORMAT {
		panic("format mismatch")
	}
	if string(r.readBytes(6)) != LUAC_DATA {
		panic("corrupted")
	}
	if r.readByte() != CINT_SIZE {
		panic("int size mismatch")
	}
	if r.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch")
	}
	if r.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch")
	}
	if r.readByte() != LUA_INTEGER_SIZE {
		panic("lua_integer size mismatch")
	}
	if r.readByte() != LUA_NUMBER_SIZE {
		panic("lua_number size mismatch")
	}
	if r.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch")
	}
	if r.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch")
	}
}

func (r *Reader) readProto(parentSource string) *Prototype {
	source := r.readString()
	if source == "" {
		source = parentSource
	}

	return &Prototype{
		Source:         source,
		LineDefine:     r.readUint32(),
		LastLineDefine: r.readUint32(),
		NumParams:      r.readByte(),
		IsVararg:       r.readByte(),
		MaxStackSize:   r.readByte(),
		Code:           r.readCode(),
		Constants:      r.readConstants(),
		Upvalues:       r.readUpvalues(),
		Protos:         r.readProtos(source),
		LineInfo:       r.readLineInfo(),
		LocVars:        r.readLocVars(),
		UpvalueNames:   r.readUpvalueNames(),
	}
}

/*
读取指令列表：四字节指令数量
每个指令长度为 4byte
*/
func (r *Reader) readCode() []uint32 {
	code := make([]uint32, r.readUint32())
	for i := range code {
		code[i] = r.readUint32()
	}
	return code
}

func (r *Reader) readConstant() interface{} {
	switch r.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return r.readByte() != 0
	case TAG_INTEGER:
		return r.readLuaInteger()
	case TAG_NUMBER:
		return r.readLuaNumber()
	case TAG_SHORT_STR:
		return r.readString()
	case TAG_LONG_STR:
		return r.readString()
	default:
		panic("corrupted")
	}
}

func (r *Reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, r.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: r.readByte(),
			Idx:     r.readByte(),
		}
	}
	return upvalues
}

func (r *Reader) readConstants() []interface{} {
	constants := make([]interface{}, r.readUint32())
	for i := range constants {
		constants[i] = r.readConstant()
	}
	return constants
}

func (r *Reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, r.readUint32())
	for i := range protos {
		protos[i] = r.readProto(parentSource)
	}
	return protos
}

func (r *Reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, r.readUint32())
	for i := range lineInfo {
		lineInfo[i] = r.readUint32()
	}
	return lineInfo
}

func (r *Reader) readLocVars() []LocVar {
	locVars := make([]LocVar, r.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: r.readString(),
			StartPC: r.readUint32(),
			EndPC:   r.readUint32(),
		}
	}
	return locVars
}

func (r *Reader) readUpvalueNames() []string {
	names := make([]string, r.readUint32())
	for i := range names {
		names[i] = r.readString()
	}
	return names
}
