// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build js

package events

import tge "github.com/thommil/tge"

// -------------------------------------------------------------------- //
// Implementation
// -------------------------------------------------------------------- //

type plugin struct {
	isInit  bool
	runtime tge.Runtime
}

func (p *plugin) subscribeProxies() {
	//NOP no support
}

func (p *plugin) unsubscribeProxies() {
	//NOP no support
}
