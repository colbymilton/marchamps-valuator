package controller

import (
	"fmt"
	"testing"
)

func TestTest(t *testing.T) {
	fmt.Println(parseLockingTrait("Play only if your identity has the [[guardian]] trait.\n<b>Hero Action</b>: Take 1 damage. Confuse the villain."))
}
