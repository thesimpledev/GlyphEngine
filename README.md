# GlyphEngine

This is a small, badly written, text based 2D engine written in Go using the EbitEngine Library.

I wrote this in a few days during the Pirate Software Game Jam. It started out beautifully structured using all of the best practices.

You can find the game I created with it here [The Social Contract](https://apocalypsetheory.itch.io/the-social-contract)

As life happened and time grew short it became a tangled ball of mud.

still I had fun writing it, I may even come back and clean it up more in the future.

## How to Use

I will fill this in later, for now just read over the code. I removed most of the assets since they

did not belong to me. Footstep sounds are optional: embed your own WAV files and pass them to
`Player.LoadFootsteps` with the game's audio context, for example in `main.go`:

```go
//go:embed assets/audio/step1.wav
var step1 []byte

if err := p.LoadFootsteps(g.AudioContext, step1); err != nil {
	log.Fatal(err)
}
```

You will also need to update the image references that I did not include for the intro scene in `internal/intro/intro.go`

## Why share

I had fun with it, and I am hoping someone else will find this interesting or useful and do something cool with it despite how badly written it is.

## ToDo List

- [x] Fix Restart Not resetting Camera
- [ ] 100% test coverage `player.go`
- [ ] Add New Audio for `player.go`
- [ ] 100% test coverage `game.go`
- [ ] 100% test coverage `intro.go`
- [ ] refactor `camera.go` into its own package
- [ ] refactor `dialog.go` into its own package
- [ ] 100% test coverage `level.go`
- [ ] 100% test coverage `camera.go`
- [ ] 100% test coverage `dialog.go`
- [ ] Patrolling NPCs
- [ ] Trust Meter
- [ ] Timed Events
- [ ] Input
