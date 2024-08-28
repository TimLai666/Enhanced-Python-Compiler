package runtime

/*
#cgo CFLAGS: -I${SRCDIR}/../../python_includes
#cgo LDFLAGS: -L${SRCDIR}/../../python_libs -lpython3.12
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
	return "Executed by CPython", nil
}
