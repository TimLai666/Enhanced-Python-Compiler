//go:build windows
// +build windows

package runtime

/*
#cgo CFLAGS: -I${SRCDIR}/../../python_includes
#cgo LDFLAGS: -L${SRCDIR}/../../python_libs/windows -lpython312 -static
#include <Python.h>

void executePythonCode(const char* code) {
    Py_Initialize();
    PyRun_SimpleString(code);
    Py_Finalize();
}
*/
import "C"

func ExecuteCPython(code string) (string, error) {
	C.executePythonCode(C.CString(code))
	return "Executed by CPython on Windows", nil
}
