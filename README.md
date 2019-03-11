<h1 align="center">TGE-GESTURE - Gestures plugin for TGE</h1>

 <p align="center">
    <a href="https://godoc.org/github.com/thommil/tge-gesture"><img src="https://godoc.org/github.com/thommil/tge-gesture?status.svg" alt="Godoc"></img></a>
    <a href="https://goreportcard.com/report/github.com/thommil/tge-gesture"><img src="https://goreportcard.com/badge/github.com/thommil/tge-gesture"  alt="Go Report Card"/></a>
</p>

Gestures support for TGE runtime - [TGE](https://github.com/thommil/tge)

This plugin is mainly used to adapt applications to touch inputs/events (mobile and mobile browser).

Supported gestures:
 * long press 
 * pinch
 * swipe

## Targets
 * Mobile
 * Browser (mobile)

## Dependencies
 * [TGE core](https://github.com/thommil/tge)

## Limitations
Currenlty no support for desktop and browser with no touch support.

## Implementation
See example at [GESTURE examples](https://github.com/Thommil/tge-examples/tree/master/plugins/tge-gesture)


```golang
package main

import (
    tge "github.com/thommil/tge"
    gesture "github.com/thommil/tge-gesture"
)

type GestureApp struct {
}

func (app *GestureApp) OnCreate(settings *tge.Settings) error {
    // Set all events or add only needed ones (ex: gesture.SwipeEventEnabled)
    settings.EventMask = tge.AllEventsEnabled
    return nil
}

func (app *GestureApp) OnStart(runtime tge.Runtime) error {
    runtime.Subscribe(gesture.LongPressEvent{}.Channel(), app.OnLongPress)
    runtime.Subscribe(gesture.PinchEvent{}.Channel(), app.OnPinch)
    runtime.Subscribe(gesture.SwipeEvent{}.Channel(), app.OnSwipe)
    return nil
}

func (app *GestureApp) OnLongPress(event tge.Event) bool {
    ...
    return false
}

func (app *GestureApp) OnPinch(event tge.Event) bool {
    ...
    return false
}

func (app *GestureApp) OnSwipe(event tge.Event) bool {
    ...
    return false
}

...

```

