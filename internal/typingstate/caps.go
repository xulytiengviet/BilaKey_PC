package typingstate

import "strings"

const DefaultDoubleTapWindowMS int64 = 360

type WordCase struct {
	UpperFirst bool
	AllUpper   bool
}

type Snapshot struct {
	SentenceStart bool
	OneShotUpper  bool
	StickyUpper   bool
	LastShiftTap  int64
}

type Capitalizer struct {
	autoInitial     bool
	autoSentence    bool
	doubleShiftCaps bool
	doubleTapWindow int64

	sentenceStart bool
	oneShotUpper  bool
	stickyUpper   bool
	lastShiftTap  int64
	shiftDown     bool
	shiftUsed     bool
}

func New(autoInitial, autoSentence, doubleShiftCaps bool) *Capitalizer {
	c := &Capitalizer{doubleTapWindow: DefaultDoubleTapWindowMS}
	c.Configure(autoInitial, autoSentence, doubleShiftCaps)
	c.Reset()
	return c
}

func (c *Capitalizer) Configure(autoInitial, autoSentence, doubleShiftCaps bool) {
	c.autoInitial = autoInitial
	c.autoSentence = autoSentence
	c.doubleShiftCaps = doubleShiftCaps
	if !doubleShiftCaps {
		c.stickyUpper = false
		c.lastShiftTap = 0
	}
}

func (c *Capitalizer) Reset() {
	c.sentenceStart = c.autoInitial
	c.oneShotUpper = false
	c.stickyUpper = false
	c.lastShiftTap = 0
	c.shiftDown = false
	c.shiftUsed = false
}

func (c *Capitalizer) BeginWord() WordCase {
	wc := WordCase{UpperFirst: c.sentenceStart || c.oneShotUpper || c.stickyUpper, AllUpper: c.stickyUpper}
	c.sentenceStart = false
	c.oneShotUpper = false
	return wc
}

func (c *Capitalizer) ObserveDelimiter(delim string) {
	if !c.autoSentence {
		return
	}
	if strings.ContainsAny(delim, ".!?") || strings.ContainsAny(delim, "\r\n") {
		c.sentenceStart = true
	}
}

func (c *Capitalizer) ShiftDown() {
	if c.shiftDown {
		return
	}
	c.shiftDown = true
	c.shiftUsed = false
}

func (c *Capitalizer) MarkShiftUsed() {
	if c.shiftDown {
		c.shiftUsed = true
	}
}

func (c *Capitalizer) ShiftUp(nowMS int64) bool {
	if !c.shiftDown {
		return false
	}
	used := c.shiftUsed
	c.shiftDown = false
	c.shiftUsed = false
	if used {
		return false
	}
	if c.stickyUpper {
		c.stickyUpper = false
		c.oneShotUpper = false
		c.lastShiftTap = 0
		return true
	}
	if c.doubleShiftCaps && c.lastShiftTap > 0 && nowMS >= c.lastShiftTap && nowMS-c.lastShiftTap <= c.doubleTapWindow {
		c.stickyUpper = true
		c.oneShotUpper = false
		c.lastShiftTap = 0
		return true
	}
	c.oneShotUpper = true
	c.lastShiftTap = nowMS
	return true
}

func (c *Capitalizer) Snapshot() Snapshot {
	return Snapshot{SentenceStart: c.sentenceStart, OneShotUpper: c.oneShotUpper, StickyUpper: c.stickyUpper, LastShiftTap: c.lastShiftTap}
}

func (c *Capitalizer) Restore(s Snapshot) {
	c.sentenceStart = s.SentenceStart
	c.oneShotUpper = s.OneShotUpper
	c.stickyUpper = s.StickyUpper && c.doubleShiftCaps
	c.lastShiftTap = s.LastShiftTap
	c.shiftDown = false
	c.shiftUsed = false
}

func (c *Capitalizer) ModeLabel() string {
	switch {
	case c.stickyUpper:
		return "SHIFT×2 LOCK"
	case c.oneShotUpper:
		return "SHIFT×1 NEXT"
	case c.sentenceStart:
		return "AUTO-CAP READY"
	default:
		return "normal"
	}
}

func (c *Capitalizer) StickyUpper() bool { return c.stickyUpper }
