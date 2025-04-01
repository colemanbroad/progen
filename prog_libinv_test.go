package main

import "testing"

func Test_libinv(t *testing.T) {
	set_lib := NewSetFromMapKeys(fn_library)
	tassert(set_lib.Size() > 0, func() { t.Error("zero size set_lib") })

	// libinv := buildLibraryInverse()
	// set_libinv_p := NewSetFromMapKeys(libinv.provides)
	// set_libinv_r := NewSetFromMapKeys(libinv.requires)
	// tassert(set_lib.Size() > 0, "zero size set_lib")
}
