package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Drawable defines the interface for components that can draw themselves.
type Drawable interface {
	Draw(screen *ebiten.Image)
}

// Updatable defines the interface for components that update every frame.
type Updatable interface {
	Update()
}

// Game is the core game state container.
// It holds an audio context, screen dimensions, and a slice of components.
type Game struct {
	AudioContext *audio.Context
	screenWidth  int
	screenHeight int
	Components   []interface{}
}

// New creates and returns a new Game instance with the given screen dimensions.
// It initializes the audio context at a sample rate of 44100 Hz.
func New(screenWidth, screenHeight int) *Game {
	return &Game{
		AudioContext: audio.NewContext(44100),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// AddComponent adds a new component to the game.
// Components can be any type; later we check if they implement Drawable and/or Updatable.
func (g *Game) AddComponent(c interface{}) {
	g.Components = append(g.Components, c)
}

// RemoveComponent marks the specified component for removal.
// Instead of immediately modifying the slice, we set the element to nil.
// The actual cleanup occurs later in cleanupComponents().
func (g *Game) RemoveComponent(c interface{}) {
	for i, comp := range g.Components {
		if comp == c {
			g.Components[i] = nil
			break
		}
	}
}

// cleanupComponents rebuilds the Components slice by filtering out nil entries.
// This avoids issues that can arise when modifying the slice during iteration.
func (g *Game) cleanupComponents() {
	var activeComponents []interface{}
	for _, comp := range g.Components {
		if comp != nil {
			activeComponents = append(activeComponents, comp)
		}
	}
	g.Components = activeComponents
}

// Draw fills the screen with black and then iterates over all components.
// If a component implements the Drawable interface, its Draw method is called.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for _, comp := range g.Components {
		if comp == nil {
			continue
		}
		// Use a clear variable name "drawable" for the type assertion.
		if drawable, ok := comp.(Drawable); ok {
			drawable.Draw(screen)
		}
	}
}

// Layout returns the logical screen dimensions used by Ebiten.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// Update iterates over all components and calls their Update method if they implement Updatable.
// Afterwards, it cleans up any components marked as nil.
func (g *Game) Update() error {
	for _, comp := range g.Components {
		if comp == nil {
			continue
		}
		if updatable, ok := comp.(Updatable); ok {
			updatable.Update()
		}
	}
	g.cleanupComponents()
	return nil
}
