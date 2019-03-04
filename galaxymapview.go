package main

import (
	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type zoomLevel int

const (
	zoom_GALAXY zoomLevel = iota
	zoom_LOCAL
)

type MapDrawable interface {
	Locatable
	burl.Drawable
}

type GalaxyMapView struct {
	burl.TileView

	galaxy *Galaxy

	cursor    burl.Coord
	highlight *burl.PulseAnimation

	zoom zoomLevel

	//for local map drawing
	localZoom   int
	localFocus  Locatable
	systemFocus *StarSystem
}

func NewGalaxyMapView(w, x, y, z int, bord bool, g *Galaxy) (gmv *GalaxyMapView) {
	gmv = new(GalaxyMapView)
	gmv.TileView = *burl.NewTileView(w, w, x, y, z, bord)
	gmv.galaxy = g

	gmv.highlight = burl.NewPulseAnimation(1, 1, 0, 1, 1, 100, 10, true)
	gmv.AddAnimation(gmv.highlight)

	return
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
		burl.LogError("No draw function for selected zoom level.")
	}
}

func (gmv *GalaxyMapView) DrawGalaxyMap() {
	if gmv.galaxy != nil {
		w, h := gmv.galaxy.Dims()
		for i := 0; i < w*h; i++ {
			x, y := i%w, i/w
			s := gmv.galaxy.GetSector(x, y)
			bright := burl.Lerp(0, 255, s.Density, 100)
			g := burl.GLYPH_FILL_SPARSE
			if bright == 0 {
				g = burl.GLYPH_NONE
			}
			gmv.Draw(x, y, g, burl.MakeOpaqueColour(bright, bright, bright), burl.COL_BLACK)
		}
	}
}

func (gmv *GalaxyMapView) calcLocalMapCoord(c Coordinates) (mc burl.Coord) {
	w, h := gmv.Dims()
	gFactor := gmv.localZoomFactor()

	//camera computation
	cX := gmv.localFocus.GetCoords().Local.X - (gFactor * float64(w) / 2)
	cY := gmv.localFocus.GetCoords().Local.Y - (gFactor * float64(h) / 2)

	mc.X = int((c.Local.X - cX) / gFactor)
	mc.Y = int((c.Local.Y - cY) / gFactor)

	return
}

func (gmv *GalaxyMapView) localZoomFactor() float64 {
	w, _ := gmv.Dims()
	return coord_LOCAL_MAX / float64(w) / float64(burl.Pow(2, gmv.localZoom))
}

func (gmv *GalaxyMapView) DrawLocalMap() {
	gmv.Reset()
	smc := gmv.calcLocalMapCoord(gmv.systemFocus.Star.GetCoords())

	//draw system things!
	for _, p := range gmv.systemFocus.Planets {
		//draw orbits. TODO: some way of culling orbit drawing. currently drawing all of them
		gmv.DrawCircle(smc.X, smc.Y, burl.RoundFloatToInt(p.oDistance/gmv.localZoomFactor()), burl.GLYPH_PERIOD, 0xFF114411, burl.COL_BLACK)
		gmv.DrawMapObject(p)
	}

	gmv.DrawMapObject(gmv.systemFocus.Star)
}

func (gmv *GalaxyMapView) DrawMapObject(o MapDrawable) {
	var c burl.Coord

	switch gmv.zoom {
	case zoom_GALAXY:
		c = o.GetCoords().Sector
	case zoom_LOCAL:
		c = gmv.calcLocalMapCoord(o.GetCoords())
	}

	gmv.DrawObject(c.X, c.Y, o)
}

//Draws a custom visual marker v on the map at map coord (not physical coord) c. Good for waypoints, etc.
func (gmv *GalaxyMapView) DrawMapMarker(c burl.Coord, v burl.Visuals) {
	gmv.DrawObject(c.X, c.Y, v)
}

func (gmv *GalaxyMapView) ToggleHighlight() {
	gmv.highlight.MoveTo(gmv.cursor.X, gmv.cursor.Y)
	gmv.highlight.Toggle()
}

func (gmv *GalaxyMapView) HandleKeypress(key sdl.Keycode) {
	//generic
	switch key {
	case sdl.K_PAGEUP:
		gmv.ZoomIn()
	case sdl.K_PAGEDOWN:
		gmv.ZoomOut()
	}

	//zoom-level specific
	switch gmv.zoom {
	case zoom_GALAXY:
		switch key {
		case sdl.K_UP:
			gmv.MoveCursor(0, -1)
		case sdl.K_DOWN:
			gmv.MoveCursor(0, 1)
		case sdl.K_LEFT:
			gmv.MoveCursor(-1, 0)
		case sdl.K_RIGHT:
			gmv.MoveCursor(1, 0)
		}
	}
}

func (gmv *GalaxyMapView) MoveCursor(dx, dy int) {
	w, h := gmv.Dims()
	if burl.CheckBounds(gmv.cursor.X+dx, gmv.cursor.Y+dy, w, h) {
		gmv.cursor.Move(dx, dy)
		gmv.highlight.Move(dx, dy, 0)
	}
}
