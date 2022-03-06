package state

import . "github.com/v0/api"

type luaValue = interface{}

// go type to lua type
func typeOf(v luaValue) LuaType {
	switch v.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	default:
		panic("unknown data")
	}
}
