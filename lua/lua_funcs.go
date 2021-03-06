package lua

/*
#cgo darwin CFLAGS: -DLUA_USE_MACOSX
#cgo darwin LDFLAGS: -lreadline

#cgo linux  CFLAGS: -DLUA_USE_LINUX
#cgo linux  LDFLAGS: -Wl,-E -ldl -lm -lreadline

#include	<stdlib.h>
#include    "lua_header.h"

static	int	X_LuaF_RegLen(void* p){
	int	len	= 0;
	if(p != 0){
		luaL_Reg* r	= (luaL_Reg*)p;
		while(r->name && r->func){
			len++;
			r++;
		}
	}
	return	len;
}

static	const char*	X_LuaF_Ident(){
	return	lua_ident;
}

*/
import "C"

import (
	"unsafe"
)

//
//	extra types.
//

type Lua_CInt int32
type Lua_CUint uint32
type Lua_CFile uintptr
type Lua_CStatePtr *C.lua_State

type Lua_Ref int
type Lua_FuncPtr *_swig_fnptr

type Lua_State unsafe.Pointer
type Lua_CFunction Lua_FuncPtr

type Lua_KContext uintptr
type Lua_KFunction Lua_FuncPtr

type Lua_Number C.lua_Number
type Lua_Integer C.lua_Integer
type Lua_Unsigned C.lua_Unsigned

var Lua_ident string = C.GoString(C.X_LuaF_Ident())

func Lua_isnumber(L Lua_State, idx int) bool {
	return C.lua_isnumber(Lua_CStatePtr(L), C.int(idx)) != 0
}

func Lua_isstring(L Lua_State, idx int) bool {
	return C.lua_isstring(Lua_CStatePtr(L), C.int(idx)) != 0
}

func Lua_iscfunction(L Lua_State, idx int) bool {
	return C.lua_iscfunction(Lua_CStatePtr(L), C.int(idx)) != 0
}

func Lua_isinteger(L Lua_State, idx int) bool {
	return C.lua_isinteger(Lua_CStatePtr(L), C.int(idx)) != 0
}

func Lua_isuserdata(L Lua_State, idx int) bool {
	return C.lua_isuserdata(Lua_CStatePtr(L), C.int(idx)) != 0
}

func Lua_isfunction(L Lua_State, n int) bool {
	return X_Lua_isfunction(L, Lua_CInt(n)) != 0
}

func Lua_istable(L Lua_State, n int) bool {
	return X_Lua_istable(L, Lua_CInt(n)) != 0
}

func Lua_islightuserdata(L Lua_State, n int) bool {
	return X_Lua_islightuserdata(L, Lua_CInt(n)) != 0
}

func Lua_isnil(L Lua_State, n int) bool {
	return X_Lua_isnil(L, Lua_CInt(n)) != 0
}

func Lua_isboolean(L Lua_State, n int) bool {
	return X_Lua_isboolean(L, Lua_CInt(n)) != 0
}

func Lua_isthread(L Lua_State, n int) bool {
	return X_Lua_isthread(L, Lua_CInt(n)) != 0
}

func Lua_isnone(L Lua_State, n int) bool {
	return X_Lua_isnone(L, Lua_CInt(n)) != 0
}

func Lua_isnoneornil(L Lua_State, n int) bool {
	return X_Lua_isnoneornil(L, Lua_CInt(n)) != 0
}

func Lua_toboolean(L Lua_State, idx int) bool {
	return C.lua_toboolean(Lua_CStatePtr(L), C.int(idx)) != 0
}

func Lua_pushboolean(L Lua_State, b bool) {
	var v C.int = 0
	if b {
		v = 1
	}
	C.lua_pushboolean(Lua_CStatePtr(L), v)
}

//
//	luaL_? impls.
//
func LuaL_checkoption(L Lua_State, arg int, def string, lst []string) int {
	l := make([](*C.char), 0, len(lst)+1)
	for i := 0; i < len(lst); i++ {
		s := C.CString(lst[i])
		defer C.free(unsafe.Pointer(s))
		l = append(l, s)
	}
	l = append(l, nil)

	s := C.CString(def)
	defer C.free(unsafe.Pointer(s))

	R := C.luaL_checkoption(Lua_CStatePtr(L), C.int(arg), s, &l[0])
	return int(R)
}

func LuaL_loadbufferx(L Lua_State, buff unsafe.Pointer, sz uint, name string, mode string) int {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	m := C.CString(mode)
	defer C.free(unsafe.Pointer(m))

	r := C.luaL_loadbufferx(Lua_CStatePtr(L), (*C.char)(buff), C.size_t(sz), n, m)
	return int(r)
}

func LuaL_setfuncs(L Lua_State, l unsafe.Pointer, nup int) {
	C.luaL_setfuncs(Lua_CStatePtr(L), (*C.luaL_Reg)(l), C.int(nup))
}

func LuaL_newlibtable(L Lua_State, l unsafe.Pointer) {
	C.lua_createtable(Lua_CStatePtr(L), 0, LuaF_RegLen(l))
}

func LuaL_newlib(L Lua_State, l unsafe.Pointer) {
	LuaL_checkversion(L)
	LuaL_newlibtable(L, l)
	LuaL_setfuncs(L, l, 0)
}

//
// lua function params' utils.
//

func LuaF_RegLen(p unsafe.Pointer) C.int {
	return C.X_LuaF_RegLen(p)
}

//
//  Lua callback types.
//

//
// convert unsafe.Pointer to Lua_Alloc
//
//extern    void*   myGoAlloc(void *ud, void *ptr, size_t osize, size_t nsize);
//func              myGoAlloc(ud unsafe.Pointer, ptr unsafe.Pointer, osize uintptr, nsize uintptr) unsafe.Pointer
func LuaF_Alloc(fp unsafe.Pointer) Lua_FuncPtr {
	return Lua_FuncPtr(fp)
}

//
// convert unsafe.Pointer to Lua_CFunction
//
//extern    int	    myGoCFunc(void* l);
//tfunc             myGoCFunc(L unsafe.Pointer) int32
func LuaF_CFunction(fp unsafe.Pointer) Lua_FuncPtr {
	return Lua_FuncPtr(fp)
}

//
// convert unsafe.Pointer to Lua_CFunction
//
//extern    int	    myGoCFunc(void* l);
//tfunc             myGoCFunc(L unsafe.Pointer) int32
func LuaF_KFunction(fp unsafe.Pointer) Lua_FuncPtr {
	return Lua_FuncPtr(fp)
}

//
// convert unsafe.Pointer to Lua_Reader
//
//extern    const char* myGoReader(void *L, void *ud, size_t *sz);
//func                  myGoReader(L unsafe.Pointer, ud unsafe.Pointer, sz uintptr)unsafe.Pointer
func LuaF_Reader(fp unsafe.Pointer) Lua_FuncPtr {
	return Lua_FuncPtr(fp)
}

//
// convert unsafe.Pointer to Lua_Writer
//
//extern    int     myGoWriter(void *L, void* p, size_t sz, void* ud);
//func              myGoWriter(L unsafe.Pointer, p unsafe.Pointer, sz uintptr, ud unsafe.Pointer)int32
func LuaF_Writer(fp unsafe.Pointer) Lua_FuncPtr {
	return Lua_FuncPtr(fp)
}

//
// convert unsafe.Pointer to Lua_Hook
//
//extern    void    myGoHook(void *L, void *ar);
//func              myGoHook(L unsafe.Pointer, ar unsafe.Pointer)
func LuaF_Hook(fp unsafe.Pointer) Lua_FuncPtr {
	return Lua_FuncPtr(fp)
}
