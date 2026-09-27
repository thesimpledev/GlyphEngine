package player

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/thesimpledev/GlyphEngine/internal/input"
)

const moveCooldown = 15

type Level interface {
	IsWalkable(x, y int) bool
	UpdateBoard(fromX, fromY, toX, toY int)
	UpdateCamera(x, y int)
}

type Player struct {
	x, y             int
	walk             []*audio.Player
	movementCooldown int
	Level            Level
	Keys             input.KeyReader
}

type PlayerMove struct {
	x int
	y int
}

type keyBinding struct {
	key  ebiten.Key
	move PlayerMove
}

var moveBindings = []keyBinding{
	{ebiten.KeyW, PlayerMove{0, -1}},
	{ebiten.KeyA, PlayerMove{-1, 0}},
	{ebiten.KeyS, PlayerMove{0, 1}},
	{ebiten.KeyD, PlayerMove{1, 0}},
}

func New() *Player {
	return &Player{Keys: input.Keyboard{}}
}

func (p *Player) Update() {
	if p.movementCooldown > 0 {
		p.movementCooldown--
		return
	}

	move, ok := p.readMove()
	if !ok {
		return
	}

	p.movementCooldown = moveCooldown
	p.move(PlayerMove{p.x + move.x, p.y + move.y})
}

func (p *Player) readMove() (PlayerMove, bool) {
	var move PlayerMove
	found := false
	for _, binding := range moveBindings {
		if p.keys().IsKeyPressed(binding.key) {
			move = binding.move
			found = true
		}
	}
	return move, found
}

func (p *Player) keys() input.KeyReader {
	if p.Keys == nil {
		return input.Keyboard{}
	}
	return p.Keys
}

func (p *Player) move(to PlayerMove) {
	if p.Level == nil || !p.Level.IsWalkable(to.x, to.y) {
		return
	}
	p.Level.UpdateBoard(p.x, p.y, to.x, to.y)
	p.x = to.x
	p.y = to.y
	p.PlayFootstep()
	p.Level.UpdateCamera(p.x, p.y)
}

func (p *Player) LoadFootsteps(ctx *audio.Context, wavs ...[]byte) error {
	if ctx == nil {
		return errors.New("player: audio context is nil")
	}
	players := make([]*audio.Player, 0, len(wavs))
	for i, data := range wavs {
		stream, err := wav.DecodeWithSampleRate(ctx.SampleRate(), bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("player: decode footstep %d: %w", i, err)
		}
		footstep, err := ctx.NewPlayer(stream)
		if err != nil {
			return fmt.Errorf("player: create footstep %d: %w", i, err)
		}
		players = append(players, footstep)
	}
	p.walk = append(p.walk, players...)
	return nil
}

func (p *Player) PlayFootstep() {
	if len(p.walk) == 0 {
		return
	}
	sound := p.walk[rand.IntN(len(p.walk))] // #nosec G404 -- footstep variety, not security sensitive
	sound.SetVolume(0.5)
	if err := sound.Rewind(); err != nil {
		log.Printf("footstep rewind failed: %v", err)
		return
	}
	sound.Play()
}

func (p *Player) SetPosition(x, y int) {
	p.x = x
	p.y = y
}
