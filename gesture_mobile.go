// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build android ios

package events

import (
	fmt "fmt"

	tge "github.com/thommil/tge"
)

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

type plugin struct {
	isInit         bool
	runtime        tge.Runtime
	firstTouchEvt  *tge.MouseEvent
	secondTouchEvt *tge.MouseEvent
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
			if mouseEvent.Button == 1 {
				p.firstTouchEvt = &mouseEvent
			} else {
				p.secondTouchEvt = &mouseEvent
			}
			return false
		case tge.TypeUp:
			if mouseEvent.Button == 1 {
				p.firstTouchEvt = nil
			} else {
				p.secondTouchEvt = nil
			}
			return false
		case tge.TypeMove:
			xOffset := int32(0)
			yOffset := int32(0)
			if mouseEvent.Button == 1 {
				fmt.Println("1")
				if p.firstTouchEvt != nil {
					fmt.Println("11")
					xOffset = mouseEvent.X - p.firstTouchEvt.X
					yOffset = mouseEvent.Y - p.firstTouchEvt.Y
				}
				if p.secondTouchEvt != nil {
					p.firstTouchEvt = &mouseEvent
				}
			} else {
				fmt.Println("2")
				if p.secondTouchEvt != nil {
					fmt.Println("21")
					xOffset = mouseEvent.X - p.secondTouchEvt.X
					yOffset = mouseEvent.Y - p.secondTouchEvt.Y
				}
				if p.firstTouchEvt != nil {
					p.secondTouchEvt = &mouseEvent
				}
			}

			if p.firstTouchEvt != nil && p.firstTouchEvt != nil {
				p.runtime.Publish(PinchEvent{
					XOffset: xOffset,
					YOffset: yOffset,
				})
				return true
			}
			return false
		}
	}

	return false
}
