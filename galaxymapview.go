package main

import (
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type zoomLevel int

const (
	zoom_GALAXY zoomLevel = iota
	zoom_LOCAL
)

type MapDrawable interface {
	Locatable
	gfx.VisualObject
}

type GalaxyMapView struct {
	ui.Element

	galaxy *Galaxy

	cursor    vec.Coord
	highlight gfx.PulseAnimation

	zoom zoomLevel

	//for local map drawing
	localZoom   int
	localFocus  Locatable
	systemFocus *StarSystem
}

func (gmv *GalaxyMapView) Init(size vec.Dims, pos vec.Coord, depth int, galaxy *Galaxy) {
	gmv.Element.Init(size, pos, depth)
	gmv.TreeNode.Init(gmv)
	gmv.galaxy = galaxy

	gmv.highlight = gfx.NewPulseAnimation(vec.Rect{vec.Coord{0, 0}, vec.Dims{1, 1}}, 0, 100, col.Pair{col.WHITE, col.WHITE})
	gmv.highlight.Repeat = true
	gmv.AddAnimation(&gmv.highlight)
}

func (gmv *GalaxyMapView) ZoomIn() {
	if gmv.zoom == zoom_GALAXY {
		gmv.zoom = zoom_LOCAL
		gmv.DrawLocalMap()
		gmv.ToggleHighlight()
	} else if gmv.zoom == zoom_LOCAL {
		if gmv.localZoom < 7 {
			gmv.localZoom += 1
			gmv.DrawLocalMap()
		}
	}
}

func (gmv *GalaxyMapView) ZoomOut() {
	if gmv.zoom == zoom_LOCAL {
		if gmv.localZoom == 0 {
			gmv.zoom = zoom_GALAXY
			gmv.DrawGalaxyMap()
			gmv.ToggleHighlight()
		} else {
			gmv.localZoom -= 1
			gmv.DrawLocalMap()
		}
	}
}

func (gmv *GalaxyMapView) DrawMap() {
	switch gmv.zoom {
	case zoom_GALAXY:
		gmv.DrawGalaxyMap()
	case zoom_LOCAL:
		gmv.DrawLocalMap()
	default:
		log.Error("No draw function for selected zoom level.")
	}
}

func (gmv *GalaxyMapView) DrawGalaxyMap() {
	gmv.Updated = true
}

func (gmv *GalaxyMapView) Render() {
	if gmv.galaxy == nil {
		return
	}
	// TODO: offset galaxy drawing to be centered inside mapview
	for cursor := range vec.EachCoordInArea(vec.Rect{vec.ZERO_COORD, gmv.galaxy.Dims()}) {
		s := gmv.galaxy.GetSector(cursor)
		bright := util.Lerp(uint8(0), uint8(255), s.Density, 100)
		g := gfx.GLYPH_FILL_SPARSE
		if bright == 0 {
			g = gfx.GLYPH_NONE
		}
		gmv.DrawVisuals(cursor, 0, gfx.NewGlyphVisuals(g, col.Pair{col.MakeOpaque(bright, bright, bright), col.BLACK}))
	}
}

func (gmv *GalaxyMapView) calcLocalMapCoord(c Coordinates) (mc vec.Coord) {
	gFactor := gmv.localZoomFactor()

	//camera computation
	cX := gmv.localFocus.GetCoords().Local.X - (gFactor * float64(gmv.Size().W) / 2)
	cY := gmv.localFocus.GetCoords().Local.Y - (gFactor * float64(gmv.Size().H) / 2)

	mc.X = int((c.Local.X - cX) / gFactor)
	mc.Y = int((c.Local.Y - cY) / gFactor)

	return
}

func (gmv *GalaxyMapView) localZoomFactor() float64 {
	return coord_LOCAL_MAX / float64(gmv.Size().W) / float64(util.Pow(2, gmv.localZoom))
}

func (gmv *GalaxyMapView) DrawLocalMap() {
	gmv.Clear()
	smc := gmv.calcLocalMapCoord(gmv.systemFocus.Star.GetCoords())

	//draw system things!
	orbitVisuals := gfx.NewGlyphVisuals(gfx.GLYPH_PERIOD, col.Pair{0xFF114411, col.BLACK})
	for _, p := range gmv.systemFocus.Planets {
		//draw orbits. TODO: some way of culling orbit drawing. currently drawing all of them
		gmv.DrawCircle(smc, 0, util.RoundFloatToInt(p.oDistance/gmv.localZoomFactor()), orbitVisuals, false)
		gmv.DrawMapObject(p)
	}

	gmv.DrawMapObject(gmv.systemFocus.Star)
}

func (gmv *GalaxyMapView) DrawMapObject(object MapDrawable) {
	var c vec.Coord

	switch gmv.zoom {
	case zoom_GALAXY:
		c = object.GetCoords().Sector
	case zoom_LOCAL:
		c = gmv.calcLocalMapCoord(object.GetCoords())
	}

	gmv.DrawObject(c, 0, object)
}

// Draws a custom visual marker v on the map at map coord (not physical coord) c. Good for waypoints, etc.
func (gmv *GalaxyMapView) DrawMapMarker(pos vec.Coord, marker gfx.Visuals) {
	gmv.DrawVisuals(pos, 0, marker)
}

func (gmv *GalaxyMapView) ToggleHighlight() {
	gmv.highlight.MoveTo(gmv.cursor)
	if gmv.highlight.IsPlaying() {
		gmv.highlight.Stop()
	} else {
		gmv.highlight.Start()
	}
}

func (gmv *GalaxyMapView) HandleKeypress(key_event *input.KeyboardEvent) (event_handled bool) {
	if key_event.Handled() || key_event.PressType != input.KEY_PRESSED {
		return
	}

	//generic
	switch key_event.Key {
	case input.K_PAGEUP:
		gmv.ZoomIn()
		event_handled = true
	case input.K_PAGEDOWN:
		gmv.ZoomOut()
		event_handled = true
	}

	//zoom-level specific
	switch gmv.zoom {
	case zoom_GALAXY:
		switch key_event.Key {
		case input.K_UP:
			gmv.MoveCursor(0, -1)
			event_handled = true
		case input.K_DOWN:
			gmv.MoveCursor(0, 1)
			event_handled = true
		case input.K_LEFT:
			gmv.MoveCursor(-1, 0)
			event_handled = true
		case input.K_RIGHT:
			gmv.MoveCursor(1, 0)
			event_handled = true
		}
	}

	return
}

func (gmv *GalaxyMapView) MoveCursor(dx, dy int) {
	new_pos := gmv.cursor.Add(vec.Coord{dx, dy})
	if new_pos.IsInside(vec.Rect{vec.ZERO_COORD, gmv.Size()}) {
		gmv.cursor.Move(dx, dy)
		gmv.highlight.MoveTo(gmv.cursor)
	}
}
