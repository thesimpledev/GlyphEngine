package player

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/exp/rand"
)

// //go:embed assets/audio/footstep_carpet_000.wav
// var step1 []byte

// //go:embed assets/audio/footstep_carpet_001.wav
// var step2 []byte

// //go:embed assets/audio/footstep_carpet_002.wav
// var step3 []byte

const (
	MOVE_COOLDOWN = 15
)

type Level interface {
	IsWalkable(x, y int) bool
	UpdateBoard(x, y, dx, dy int)
}

type Camera interface {
	UpdateCamera(x, y int)
}

type Player struct {
	x, y             int
	walk             []*audio.Player
	movementCooldown int
	Level            Level
	Camera           Camera
}

type PlayerMove struct {
	x int
	y int
}

func New() *Player {
	return &Player{}
}

func (p *Player) Update() {
	if p.movementCooldown > 0 {
		p.movementCooldown--
		return
	}

	currentMove := PlayerMove{0, 0}
	hasMoveInput := false
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		currentMove = PlayerMove{0, -1}
		hasMoveInput = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		currentMove = PlayerMove{-1, 0}
		hasMoveInput = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		currentMove = PlayerMove{0, 1}
		hasMoveInput = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		currentMove = PlayerMove{1, 0}
		hasMoveInput = true
	}

	if hasMoveInput {
		p.movementCooldown = MOVE_COOLDOWN
		newPosition := PlayerMove{p.x + currentMove.x, p.y + currentMove.y}
		p.move(newPosition)
	}

}

func (p *Player) move(newPosition PlayerMove) {
	if p.Level.IsWalkable(newPosition.x, newPosition.y) {
		p.Level.UpdateBoard(p.x, p.y, newPosition.x, newPosition.y)
		p.x = newPosition.x
		p.y = newPosition.y
		// p.PlayFootstep()
		p.Level.UpdateCamera(p.x, p.y)
	}
}

func (p *Player) PlayFootstep() {
	index := rand.Intn(len(p.walk))
	sound := p.walk[index]
	sound.SetVolume(0.5)
	sound.Rewind()
	sound.Play()
}

func (p *Player) SetPosition(x, y int) {
	p.x = x
	p.y = y
}

func loadAudio() {
	// d1, err := wav.Decode(ac, bytes.NewReader(step1))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// p1, err := ac.NewPlayer(d1)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// d2, err := wav.Decode(ac, bytes.NewReader(step2))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// p2, err := ac.NewPlayer(d2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// d3, err := wav.Decode(ac, bytes.NewReader(step3))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// p3, err := ac.NewPlayer(d3)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// walk: []*audio.Player{p1, p2, p3},
}
