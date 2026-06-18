package ui

/*
#include <stdlib.h>
*/
import "C"

//export cliperPopupItems
func cliperPopupItems() *C.char {
	if activeApp == nil {
		return C.CString("[]")
	}
	return C.CString(activeApp.popupItemsJSON())
}

//export cliperPopupItemClicked
func cliperPopupItemClicked(index C.int) {
	if activeApp == nil {
		return
	}
	activeApp.pasteHistoryItem(int(index))
}
