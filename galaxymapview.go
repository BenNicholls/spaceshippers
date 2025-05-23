package main

import (
	"slices"

	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/gfx/col"
	"github.com/bennicholls/tyumi/gfx/ui"
	"github.com/bennicholls/tyumi/input"
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

	cursor       vec.Coord
	OnCursorMove func(new_pos vec.Coord) //callback called when cursor is moved
	highlight    gfx.PulseAnimation      //animation object for cursor

	galacticObjects []MapDrawable //objects to draw on the galactic map
	localObjects    []MapDrawable //objects to draw on the local map
	markerLayer     ui.Element

	zoom             zoomLevel
	OnZoomTypeChange func(zoom zoomLevel) // callback called when map is zoomed to/from galaxy/local mode

	//for local map drawing
	localZoom   int
	localFocus  Locatable
	systemFocus *StarSystem
}

func (gmv *GalaxyMapView) Init(size vec.Dims, pos vec.Coord, depth int, galaxy *Galaxy) {
	gmv.Element.Init(size, pos, depth)
	gmv.TreeNode.Init(gmv)
	gmv.galaxy = galaxy

	gmv.markerLayer.Init(size, vec.ZERO_COORD, 10)
	gmv.markerLayer.SetDefaultVisuals(gfx.Visuals{
		Mode:    gfx.DRAW_NONE,
		Colours: col.Pair{col.WHITE, col.BLACK},
	})
	gmv.markerLayer.OnRender = gmv.DrawObjectMarkers
	gmv.AddChild(&gmv.markerLayer)

	gmv.cursor = vec.Coord{gmv.Size().W / 2, gmv.Size().H / 2}
	gmv.highlight = gfx.NewPulseAnimation(vec.Dims{1, 1}.Bounds(), 0, 30, col.Pair{col.NONE, col.WHITE})
	gmv.highlight.Repeat = true
	gmv.markerLayer.AddAnimation(&gmv.highlight)
	gmv.ToggleHighlight()
}

func (gmv *GalaxyMapView) ZoomIn() {
	switch gmv.zoom {
	case zoom_GALAXY:
		gmv.zoom = zoom_LOCAL
		gmv.ToggleHighlight()
		if gmv.OnZoomTypeChange != nil {
			gmv.OnZoomTypeChange(gmv.zoom)
		}
	case zoom_LOCAL:
		if gmv.localZoom == 7 {
			return
		}

		gmv.localZoom += 1
	}

	gmv.Updated = true
	gmv.markerLayer.Updated = true
}

func (gmv *GalaxyMapView) ZoomOut() {
	switch gmv.zoom {
	case zoom_GALAXY:
		return
	case zoom_LOCAL:
		if gmv.localZoom == 0 {
			gmv.zoom = zoom_GALAXY
			gmv.ToggleHighlight()
			if gmv.OnZoomTypeChange != nil {
				gmv.OnZoomTypeChange(gmv.zoom)
			}
		} else {
			gmv.localZoom -= 1
		}
	}

	gmv.Updated = true
	gmv.markerLayer.Updated = true
}

func (gmv *GalaxyMapView) Render() {
	if gmv.galaxy == nil {
		return
	}

	switch gmv.zoom {
	case zoom_GALAXY:
		gmv.ClearAtDepth(1, gmv.DrawableArea())
		// TODO: offset galaxy drawing to be centered inside mapview
		for cursor := range vec.EachCoordInArea(gmv.galaxy.Dims()) {
			s := gmv.galaxy.GetSector(cursor)
			bright := util.Lerp[uint8](0, 255, s.Density, 100)
			g := gfx.GLYPH_FILL_SPARSE
			if bright == 0 {
				g = gfx.GLYPH_NONE
			}
			gmv.DrawVisuals(cursor, 0, gfx.NewGlyphVisuals(g, col.Pair{col.MakeOpaque(bright, bright, bright), col.NONE}))
		}

		if earth := gmv.galaxy.GetEarth(); earth.IsKnown() {
			gmv.DrawMapObject(earth.(MapDrawable))
		}

	case zoom_LOCAL:
		gmv.ClearAtDepth(1, gmv.DrawableArea())
		smc := gmv.calcLocalMapCoord(gmv.systemFocus.Star.GetCoords())

		//draw system things!
		orbitVisuals := gfx.NewGlyphVisuals(gfx.GLYPH_PERIOD, col.Pair{0xFF114411, col.BLACK})
		//draw orbits. TODO: some way of culling orbit drawing. currently drawing all of them
		for _, p := range gmv.systemFocus.Planets {
			gmv.DrawCircle(smc, 0, util.RoundFloatToInt(p.oDistance/gmv.localZoomFactor()), orbitVisuals, false)
			gmv.DrawMapObject(p)
		}

		gmv.DrawMapObject(gmv.systemFocus.Star)
	}
}

// TODO: cache this computation for gfactor and the camera, it doesn't change often.
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

// Draws an object to the map (NOT the marker layer!). Use for static bodies.
func (gmv *GalaxyMapView) DrawMapObject(object MapDrawable) {
	var c vec.Coord

	switch gmv.zoom {
	case zoom_GALAXY:
		c := object.GetCoords().Sector
		gmv.DrawObject(c, 1, object)
	case zoom_LOCAL:
		c = gmv.calcLocalMapCoord(object.GetCoords())
		gmv.DrawObject(c, 1, object)
	}
}

// Adds an object to be drawn in both galactic and local modes.
func (gmv *GalaxyMapView) AddMapObjectMarker(object MapDrawable) {
	gmv.AddGalacticMapObjectMarker(object)
	gmv.AddLocalMapObjectMarker(object)
}

// Adds an object to be drawn on the galactic version of the map.
func (gmv *GalaxyMapView) AddGalacticMapObjectMarker(object MapDrawable) {
	if !slices.Contains(gmv.galacticObjects, object) {
		gmv.galacticObjects = append(gmv.galacticObjects, object)
	}
}

// Adds an object to be drawn on the local version of the map.
func (gmv *GalaxyMapView) AddLocalMapObjectMarker(object MapDrawable) {
	if !slices.Contains(gmv.localObjects, object) {
		gmv.localObjects = append(gmv.localObjects, object)
	}
}

func (gmv *GalaxyMapView) DrawObjectMarkers() {
	gmv.markerLayer.Clear()

	switch gmv.zoom {
	case zoom_GALAXY:
		for _, object := range gmv.galacticObjects {
			c := object.GetCoords().Sector
			gmv.markerLayer.DrawObject(c, 0, object)
		}
	case zoom_LOCAL:
		for _, object := range gmv.localObjects {
			c := gmv.calcLocalMapCoord(object.GetCoords())
			gmv.markerLayer.DrawObject(c, 0, object)
		}
	}
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
		if dir := key_event.Direction(); dir != vec.DIR_NONE {
			gmv.MoveCursor(dir)
			event_handled = true
		}
	}

	return
}

func (gmv *GalaxyMapView) SetCursor(pos vec.Coord) {
	if !pos.IsInside(gmv.Size()) || gmv.cursor == pos {
		return
	}

	gmv.cursor = pos
	gmv.highlight.MoveTo(gmv.cursor)
	if gmv.OnCursorMove != nil {
		gmv.OnCursorMove(gmv.cursor)
	}
}

func (gmv *GalaxyMapView) MoveCursor(dir vec.Direction) {
	gmv.SetCursor(gmv.cursor.Step(dir))
}
