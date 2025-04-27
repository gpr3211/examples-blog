package main

import (
	"github.com/ebitengine/debugui"
)

type DebugUI struct {
	Main                debugui.DebugUI
	inputCapturingState debugui.InputCapturingState
	HiRes               bool
	Lang                bool
}
