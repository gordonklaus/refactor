// Copyright 2014 Gordon Klaus. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"code.google.com/p/gordon-go/go/types"
	. "code.google.com/p/gordon-go/gui"
)

type sliceNode struct {
	*nodeBase
	x, low, high, max, y *port
}

func newSliceNode() *sliceNode {
	n := &sliceNode{}
	n.nodeBase = newNodeBase(n)
	n.x = n.newInput(nil)
	n.x.connsChanged = n.connsChanged
	n.low = n.newInput(newVar("low", nil))
	n.high = n.newInput(newVar("high", nil))
	n.y = n.newOutput(nil)
	n.text.SetText("[:]")
	n.text.SetTextColor(color(&types.Func{}, true, false))
	n.connsChanged()
	return n
}

func (n *sliceNode) connectable(t types.Type, dst *port) bool {
	if dst == n.x {
		ok := false
		switch t := underlying(t).(type) {
		case *types.Basic:
			ok = t.Info&types.IsString != 0 && n.max == nil
		case *types.Slice:
			ok = true
		case *types.Pointer:
			_, ok = underlying(t.Elem).(*types.Array)
		}
		return ok
	}
	return assignable(t, dst.obj.Type)
}

func (n *sliceNode) connsChanged() {
	t := untypedToTyped(inputType(n.x))
	var i types.Type
	if t != nil {
		i = types.Typ[types.Int]
	}
	n.x.setType(t)
	n.low.setType(i)
	if n.high != nil {
		n.high.setType(i)
	}
	if n.max != nil {
		n.max.setType(i)
	}
	n.y.setType(t)
}

func (n *sliceNode) KeyPress(event KeyEvent) {
	if event.Text == "," {
		_, str := underlying(inputType(n.x)).(*types.Basic)
		switch {
		case n.high == nil:
			n.high = n.newInput(newVar("high", nil))
		case n.max == nil && !str:
			n.max = n.newInput(newVar("max", nil))
		default:
			n.removePortBase(n.high)
			n.high = nil
			if !str {
				n.removePortBase(n.max)
				n.max = nil
			}
		}
		n.connsChanged()
		SetKeyFocus(n)
	} else {
		n.nodeBase.KeyPress(event)
	}
}
