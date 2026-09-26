package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thesimpledev/GlyphEngine/internal/game"
	"github.com/thesimpledev/GlyphEngine/internal/intro"
	"github.com/thesimpledev/GlyphEngine/internal/level"
	"github.com/thesimpledev/GlyphEngine/internal/player"
)

const (
	title        = "The Social Contract"
	screenWidth  = 1280
	screenHeight = 720
)

func main() {
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := game.New(screenWidth, screenHeight)

	p := player.New()

	i := intro.New()

	l := level.New()
	l.Player = p

	p.Level = l

	i.Game = g
	i.Level = l

	g.AddComponent(i)
	g.AddComponent(l)
	g.AddComponent(p)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
