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

type plugin struct {
	isInit                bool
	runtime               tge.Runtime
	gestureMode           bool
	previousPinchDistance int32
	previousPinchEvents   [3]*tge.MouseEvent
}

func (p *plugin) subscribeProxies() {
	settings := p.runtime.GetSettings()
	if settings.EventMask&PinchEventEnabled != 0 {
		p.runtime.Subscribe(tge.MouseEvent{}.Channel(), p.mouseEventProxy)
	}
}

func (p *plugin) unsubscribeProxies() {
	settings := p.runtime.GetSettings()
	if settings.EventMask&PinchEventEnabled != 0 {
		p.runtime.Unsubscribe(tge.MouseEvent{}.Channel(), p.mouseEventProxy)
	}
}

func (p *plugin) mouseEventProxy(event tge.Event) bool {
	mouseEvent := event.(tge.MouseEvent)

	if mouseEvent.Button < 3 {
		switch mouseEvent.Type {
		case tge.TypeDown:
			if mouseEvent.Button == 2 {
				p.gestureMode = true
			}
			p.previousPinchEvents[mouseEvent.Button] = &mouseEvent
		case tge.TypeUp:
			p.gestureMode = false
			p.previousPinchDistance = 0
		case tge.TypeMove:
			if p.gestureMode {
				var xOffset, yOffset float64
				if mouseEvent.Button == 1 {
					xOffset = float64(mouseEvent.X - p.previousPinchEvents[2].X)
					yOffset = float64(mouseEvent.Y - p.previousPinchEvents[2].Y)
				} else {
					xOffset = float64(mouseEvent.X - p.previousPinchEvents[1].X)
					yOffset = float64(mouseEvent.Y - p.previousPinchEvents[1].Y)
				}
				distance := int32(math.Sqrt(math.Pow(xOffset, 2) + math.Pow(yOffset, 2)))
				if p.previousPinchDistance != 0 {
					delta := distance - p.previousPinchDistance
					ratio := float32(math.Abs(xOffset / yOffset))
					p.previousPinchEvents[mouseEvent.Button] = &mouseEvent
					p.previousPinchDistance = distance
					p.runtime.Publish(PinchEvent{
						Delta: delta,
						Ratio: ratio,
					})
				} else {
					p.previousPinchDistance = distance
				}
				return true
			}
		}
	}
	return false
}
