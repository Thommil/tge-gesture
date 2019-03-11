// Copyright (c) 2019 Thomas MILLET. All rights reserved.

package events

import (
	fmt "fmt"
	"math"
	"time"

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

// LongPressEventEnabled enabled long press events recevier on App
const LongPressEventEnabled = 0x10

// LongPressEvent for touch gesture long press event
type LongPressEvent struct {
	X, Y int32
}

// Channel for long press gesture/events
func (e LongPressEvent) Channel() string {
	return "longpress"
}

// PinchEventEnabled enabled pinch events recevier on App
const PinchEventEnabled = 0x20

// PinchEvent for touch gesture pinch event
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

// SwipeEventEnabled enabled swipe events recevier on App
const SwipeEventEnabled = 0x40

// SwipeEvent for touch gesture swipe event
type SwipeEvent struct {
	X, Y int32
	// Angle of swipe in radians
	Angle float64
	// Velocity of swipe in pixels per frame
	Velocity int32
}

// Channel for swipe gesture/events
func (e SwipeEvent) Channel() string {
	return "swipe"
}

// SetLongPressTreshold sets the delay before triggering longpress (default: 1s)
func SetLongPressTreshold(delay time.Duration) {
	longPressThreshold = int32(delay.Seconds() * 500)
}

// SetSwipeTreshold sets pixels variation average triggering swipe (default: 40px)
func SetSwipeTreshold(value int32) {
	swipeThreshold = value
}

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

var longPressThreshold = int32(time.Duration(50*time.Millisecond).Seconds() * 500) // empiric
var swipeThreshold = int32(40)

type plugin struct {
	isInit                bool
	runtime               tge.Runtime
	previousEvents        []*tge.MouseEvent
	longpressCounter      int32
	swipeCounter          int32
	swipeStartEvent       *tge.MouseEvent
	pinchMode             bool
	previousPinchDistance int32
}

func (p *plugin) subscribeProxies() {
	settings := p.runtime.GetSettings()
	p.previousEvents = make([]*tge.MouseEvent, 2)

	if settings.EventMask&LongPressEventEnabled != 0 {
		p.longpressCounter = -1
		p.runtime.Subscribe(tge.MouseEvent{}.Channel(), p.longpressProxy)
	}
	if settings.EventMask&SwipeEventEnabled != 0 {
		p.swipeCounter = 0
		p.swipeStartEvent = nil
		p.runtime.Subscribe(tge.MouseEvent{}.Channel(), p.swipeProxy)
	}
	if settings.EventMask&PinchEventEnabled != 0 {
		p.pinchMode = false
		p.previousPinchDistance = 0
		p.runtime.Subscribe(tge.MouseEvent{}.Channel(), p.pinchProxy)
	}
}

func (p *plugin) unsubscribeProxies() {
	settings := p.runtime.GetSettings()

	if settings.EventMask&PinchEventEnabled != 0 {
		p.runtime.Unsubscribe(tge.MouseEvent{}.Channel(), p.longpressProxy)
	}
	if settings.EventMask&SwipeEventEnabled != 0 {
		p.runtime.Unsubscribe(tge.MouseEvent{}.Channel(), p.swipeProxy)
	}
	if settings.EventMask&LongPressEventEnabled != 0 {
		p.runtime.Unsubscribe(tge.MouseEvent{}.Channel(), p.longpressProxy)
	}
}

func (p *plugin) longpressProxy(event tge.Event) bool {
	mouseEvent := event.(tge.MouseEvent)

	if mouseEvent.Button == tge.TouchFirst {
		switch mouseEvent.Type {
		case tge.TypeDown:
			p.longpressCounter = 0
			p.previousEvents[0] = &mouseEvent
		case tge.TypeMove:
			if p.longpressCounter > longPressThreshold {
				return true
			}

			if (p.previousEvents[0] != nil) &&
				(p.longpressCounter >= 0) &&
				(p.previousEvents[0].X == mouseEvent.X) &&
				(p.previousEvents[0].Y == mouseEvent.Y) {
				p.longpressCounter++
			} else {
				p.longpressCounter = -1
			}

			p.previousEvents[0] = &mouseEvent

			if p.longpressCounter > longPressThreshold {
				p.longpressCounter = -1
				p.runtime.Publish(LongPressEvent{
					X: mouseEvent.X,
					Y: mouseEvent.Y,
				})
				return true
			}
		case tge.TypeUp:
			p.longpressCounter = -1
		}
	} else {
		p.longpressCounter = -1
	}
	return false
}

func (p *plugin) swipeProxy(event tge.Event) bool {
	mouseEvent := event.(tge.MouseEvent)

	if mouseEvent.Button == tge.TouchFirst && !p.pinchMode {
		switch mouseEvent.Type {
		case tge.TypeDown:
			p.swipeCounter = 0
			p.swipeStartEvent = &mouseEvent
		case tge.TypeMove:
			p.swipeCounter++
		case tge.TypeUp:
			if p.swipeStartEvent != nil {
				xOffset := float64(mouseEvent.X - p.swipeStartEvent.X)
				yOffset := float64(mouseEvent.Y - p.swipeStartEvent.Y)
				velocity := int32(math.Sqrt(math.Pow(xOffset, 2)+math.Pow(yOffset, 2))) / p.swipeCounter
				if velocity > swipeThreshold {
					p.runtime.Publish(SwipeEvent{
						X:        p.swipeStartEvent.X,
						Y:        p.swipeStartEvent.Y,
						Velocity: velocity,
						Angle:    math.Atan2(-yOffset, xOffset),
					})
					return true
				}
				p.swipeCounter = 0
				p.swipeStartEvent = nil
			}
		}
	} else {
		p.swipeCounter = 0
		p.swipeStartEvent = nil
	}

	return false
}

func (p *plugin) pinchProxy(event tge.Event) bool {
	mouseEvent := event.(tge.MouseEvent)

	// In any case, just touch 1 and 2 are handled
	if mouseEvent.Button == tge.TouchFirst || mouseEvent.Button == tge.TouchSecond {
		switch mouseEvent.Type {
		case tge.TypeDown:
			if mouseEvent.Button == tge.TouchSecond {
				p.pinchMode = true
			}
			switch mouseEvent.Button {
			case tge.TouchFirst:
				p.previousEvents[0] = &mouseEvent
			case tge.TouchSecond:
				p.previousEvents[1] = &mouseEvent
			}
		case tge.TypeUp:
			if p.pinchMode {
				p.pinchMode = false
				p.previousPinchDistance = 0
				p.previousEvents[0] = nil
				p.previousEvents[1] = nil
			}
		case tge.TypeMove:
			// Pinchmode
			if p.pinchMode {
				var xOffset, yOffset float64
				if mouseEvent.Button == tge.TouchFirst && p.previousEvents[1] != nil {
					xOffset = float64(mouseEvent.X - p.previousEvents[1].X)
					yOffset = float64(mouseEvent.Y - p.previousEvents[1].Y)
				} else if p.previousEvents[0] != nil {
					xOffset = float64(mouseEvent.X - p.previousEvents[0].X)
					yOffset = float64(mouseEvent.Y - p.previousEvents[0].Y)
				}
				distance := int32(math.Sqrt(math.Pow(xOffset, 2) + math.Pow(yOffset, 2)))
				if p.previousPinchDistance != 0 {
					delta := distance - p.previousPinchDistance
					ratio := float32(math.Abs(xOffset / yOffset))
					switch mouseEvent.Button {
					case tge.TouchFirst:
						p.previousEvents[0] = &mouseEvent
					case tge.TouchSecond:
						p.previousEvents[1] = &mouseEvent
					}
					if delta != 0 {
						p.previousPinchDistance = distance
						p.runtime.Publish(PinchEvent{
							Delta: delta,
							Ratio: ratio,
						})
					}
				} else {
					p.previousPinchDistance = distance
				}
				return true
			} else {
				switch mouseEvent.Button {
				case tge.TouchFirst:
					p.previousEvents[0] = &mouseEvent
				case tge.TouchSecond:
					p.previousEvents[1] = &mouseEvent
				}
			}
		}
	}
	return false
}
