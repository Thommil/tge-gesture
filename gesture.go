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
const Name = "gesture"

var _pluginInstance = &plugin{}

func init() {
	tge.Register(_pluginInstance)
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
// Gesture events
// -------------------------------------------------------------------- //

// PinchEventEnabled enabled pinch events recevier on App
const PinchEventEnabled = 0x20

// PinchEvent for mobile gesture pinch event
type PinchEvent struct {
	// Delta indicates the variation of distance between touches
	Delta int32
	// Ratio between X and Y pinch
	Ratio float32
}

// Channel for pinch gesture/events
func (e PinchEvent) Channel() string {
	return "pinch"
}
