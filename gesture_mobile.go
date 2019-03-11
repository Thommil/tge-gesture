// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build android ios

package events

import (
	"math"

	tge "github.com/thommil/tge"
)

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

const longPressThreshold = 30
const swipeThreshold = 40

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
	p.previousEvents = make([]*tge.MouseEvent, 10)

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

func (p *plugin) checkSkipMode(mouseEvent tge.MouseEvent) bool {
	return mouseEvent.Button == 1 && mouseEvent.Type == tge.TypeUp
}

func (p *plugin) longpressProxy(event tge.Event) bool {
	mouseEvent := event.(tge.MouseEvent)

	if mouseEvent.Button == 1 {
		switch mouseEvent.Type {
		case tge.TypeDown:
			p.longpressCounter = 0
			p.previousEvents[mouseEvent.Button] = &mouseEvent
		case tge.TypeMove:
			if p.longpressCounter > longPressThreshold {
				return true
			}

			if (p.previousEvents[mouseEvent.Button] != nil) &&
				(p.longpressCounter >= 0) &&
				(p.previousEvents[mouseEvent.Button].X == mouseEvent.X) &&
				(p.previousEvents[mouseEvent.Button].Y == mouseEvent.Y) {
				p.longpressCounter++
			} else {
				p.longpressCounter = -1
			}

			p.previousEvents[mouseEvent.Button] = &mouseEvent

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

	if mouseEvent.Button == 1 && !p.pinchMode {
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
				distance := int32(math.Sqrt(math.Pow(xOffset, 2)+math.Pow(yOffset, 2))) / p.swipeCounter
				p.swipeStartEvent = nil
				p.swipeCounter = 0
				if distance > swipeThreshold {
					p.runtime.Publish(SwipeEvent{
						Angle: 0,
					})
					return true
				}
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
	if mouseEvent.Button < 3 {
		switch mouseEvent.Type {
		case tge.TypeDown:
			if mouseEvent.Button == 2 {
				p.pinchMode = true
			}
			p.previousEvents[mouseEvent.Button] = &mouseEvent
		case tge.TypeUp:
			if p.pinchMode {
				p.pinchMode = false
				p.previousPinchDistance = 0
				p.previousEvents[mouseEvent.Button] = nil
			}
		case tge.TypeMove:
			// Pinchmode
			if p.pinchMode {
				var xOffset, yOffset float64
				if mouseEvent.Button == 1 && p.previousEvents[2] != nil {
					xOffset = float64(mouseEvent.X - p.previousEvents[2].X)
					yOffset = float64(mouseEvent.Y - p.previousEvents[2].Y)
				} else if p.previousEvents[1] != nil {
					xOffset = float64(mouseEvent.X - p.previousEvents[1].X)
					yOffset = float64(mouseEvent.Y - p.previousEvents[1].Y)
				}
				distance := int32(math.Sqrt(math.Pow(xOffset, 2) + math.Pow(yOffset, 2)))
				if p.previousPinchDistance != 0 {
					delta := distance - p.previousPinchDistance
					ratio := float32(math.Abs(xOffset / yOffset))
					p.previousEvents[mouseEvent.Button] = &mouseEvent
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
				p.previousEvents[mouseEvent.Button] = &mouseEvent
			}
		}
	}
	return false
}
