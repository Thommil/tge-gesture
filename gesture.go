// Copyright (c) 2019 Thomas MILLET. All rights reserved.

package events

import (
	fmt "fmt"

	tge "github.com/thommil/tge"
)

// -------------------------------------------------------------------- //
// Plugin definition
// -------------------------------------------------------------------- //

// Name name of the plugin
const Name = "events"

type plugin struct {
	isInit  bool
	runtime tge.Runtime
}

func (p *plugin) Init(runtime tge.Runtime) error {
	if !p.isInit {
		p.runtime = runtime
		p.subscribeProxies()
		return nil
	}
	return fmt.Errorf("Already initialized")
}

func (p *plugin) GetName() string {
	return Name
}

func (p *plugin) Dispose() {
	p.unsubscribeProxies()
	p.runtime = nil
	p.isInit = false
}

// -------------------------------------------------------------------- //
// Events
// -------------------------------------------------------------------- //

func (p *plugin) subscribeProxies() {
	settings := p.runtime.GetSettings()
	if settings.EventMask&PinchEventEnabled != 0 {

	}
}

func (p *plugin) unsubscribeProxies() {
	settings := p.runtime.GetSettings()
	if settings.EventMask&PinchEventEnabled != 0 {

	}
}

// -------------------------------------------------------------------- //
// Gesture events
// -------------------------------------------------------------------- //

// PinchEventEnabled enabled pinch events recevier on App
const PinchEventEnabled = 0x20

// PinchEvent for mobile gesture pinch event
type PinchEvent struct {
	// XOffset indicates zoom on X
	XOffset int32
	// YOffset indicates zoom on Y
	YOffset int32
}

func mouseEventProxy(event tge.Event) bool {
	return false
}
