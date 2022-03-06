package state

type LuaStack struct {
	slots []luaValue
	top   int
}

func NewLuaStack(size int) *LuaStack {
	return &LuaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

func (l *LuaStack) check(n int) {
	free := len(l.slots) - l.top
	for i := free; i < n; i++ {
		l.slots = append(l.slots, nil)
	}
}

func (l *LuaStack) push(v luaValue) {
	if l.top == len(l.slots) {
		panic("stack overflow")
	}

	l.slots[l.top] = v
	l.top++
}

func (l *LuaStack) pop() luaValue {
	if l.top < 1 {
		panic("stack overflow")
	}

	l.top--
	v := l.slots[l.top]
	l.slots[l.top] = nil
	return v
}

func (l *LuaStack) absIndex(index int) int {
	if index < 0 {
		return index + 1 + l.top
	}
	return index
}

func (l *LuaStack) isValid(index int) bool {
	index = l.absIndex(index)
	return index <= l.top && index > 0
}

func (l *LuaStack) get(n int) luaValue {
	absIdx := l.absIndex(n)
	if absIdx > 0 && absIdx <= l.top {
		return l.slots[absIdx-1]
	}
	return nil
}

func (l *LuaStack) set(idx int, val luaValue) {
	absIdx := l.absIndex(idx)
	if absIdx > 0 && absIdx <= l.top {
		l.slots[absIdx-1] = val
		return
	}
	panic("invalid index")
}

func (l *LuaStack) reverse(from, to int) {
	slots := l.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
