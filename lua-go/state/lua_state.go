package state

import (
	"fmt"
	. "github.com/v0/api"
)

type luaState struct {
	stack *LuaStack
}

func (l luaState) GetTop() int {
	return l.stack.top
}

func (l luaState) AbsIndex(idx int) int {
	return l.stack.absIndex(idx)
}

func (l luaState) CheckStack(n int) bool {
	l.stack.check(n)
	return true
}

func (l luaState) Pop(n int) {
	l.SetTop(-n - 1)
}

func (l luaState) Copy(fromIdx, toIdx int) {
	val := l.stack.get(fromIdx)
	l.stack.set(toIdx, val)
}

func (l luaState) PushValue(idx int) {
	l.stack.push(l.stack.get(idx))
}

func (l luaState) Replace(idx int) {
	val := l.stack.pop()
	l.stack.set(idx, val)
}

// 弹出栈顶值，插入指定位置
func (l luaState) Insert(idx int) {
	l.Rotate(idx, 1)
}

func (l luaState) Remove(idx int) {
	l.Rotate(idx, -1)
	l.Pop(1)
}

func (l luaState) Rotate(idx, n int) {
	t := l.stack.top - 1
	p := l.stack.absIndex(idx) - 1
	var m int

	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}

	l.stack.reverse(p, m)
	l.stack.reverse(m+1, t)
	l.stack.reverse(p, t)
}

func (l luaState) SetTop(idx int) {
	top := l.AbsIndex(idx)
	if top < 0 {
		panic("stack underflow")
	}

	n := top - l.stack.top
	if n > 0 {
		l.CheckStack(n)
		return
	}

	if n < 0 {
		for i := 0; i < n; i++ {
			l.stack.pop()
		}
		return
	}
}

func (l *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

func (l luaState) Type(idx int) LuaType {
	if l.stack.isValid(idx) {
		return LUA_TNONE
	}
	return typeOf(l.stack.get(idx))
}

func (l *luaState) IsNone(idx int) bool {
	return l.Type(idx) == LUA_TNONE
}

func (l *luaState) IsNil(idx int) bool {
	return l.Type(idx) == LUA_TNIL
}

func (l *luaState) IsNoneOrNil(idx int) bool {
	return l.Type(idx) <= LUA_TNIL
}

func (l *luaState) IsBoolean(idx int) bool {
	return l.Type(idx) == LUA_TBOOLEAN
}

func (l *luaState) IsTable(idx int) bool {
	return l.Type(idx) == LUA_TTABLE
}

func (l *luaState) IsFunction(idx int) bool {
	return l.Type(idx) == LUA_TFUNCTION
}

func (l *luaState) IsThread(idx int) bool {
	return l.Type(idx) == LUA_TTHREAD
}

func (l *luaState) IsString(idx int) bool {
	t := l.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (l *luaState) IsNumber(idx int) bool {
	_, ok := l.ToNumberX(idx)
	return ok
}

func (l *luaState) IsInteger(idx int) bool {
	val := l.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

func (l luaState) ToBoolean(idx int) bool {
	val := l.stack.get(idx)
	return convertToBoolean(val)
}

func convertToBoolean[T luaValue](v T) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func (l luaState) ToInteger(idx int) int64 {
	n, _ := l.ToIntegerX(idx)
	return n
}

func (l luaState) ToIntegerX(idx int) (int64, bool) {
	// TODO add transform,
	v := l.stack.get(idx)
	i, ok := v.(int64)
	return i, ok
}

func (l luaState) ToNumber(idx int) float64 {
	n, _ := l.ToNumberX(idx)
	return n
}

func (l luaState) ToNumberX(idx int) (float64, bool) {
	v := l.stack.get(idx)
	switch x := v.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (l luaState) ToString(idx int) string {
	s, _ := l.ToStringX(idx)
	return s
}

func (l luaState) ToStringX(idx int) (string, bool) {
	v := l.stack.get(idx)
	switch x := v.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x)
		l.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}

func (l *luaState) PushNil() {
	l.stack.push(nil)
}

func (l *luaState) PushBoolean(b bool) {
	l.stack.push(b)
}

func (l *luaState) PushInteger(n int64) {
	l.stack.push(n)
}

func (l *luaState) PushNumber(n float64) {
	l.stack.push(n)
}

func (l *luaState) PushString(s string) {
	l.stack.push(s)
}

var _ LuaState = (*luaState)(nil)

func NewLuaState() *luaState {
	return &luaState{
		stack: NewLuaStack(20),
	}
}
