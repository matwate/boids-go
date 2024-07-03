package main

import (
	"math"
	"math/rand/v2"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var boids []*Boid

type Boid struct {
	Pos rl.Vector2
	Vel rl.Vector2

	Cohesion   rl.Vector2
	Separation rl.Vector2
	Alignment  rl.Vector2

	CohesionStr   float32
	SeparationStr float32
	AlignmentStr  float32

	MaxVel float32

	CloseRadius float32

	Radius float32

	fov float32

	once sync.Once
}

func NewBoid(pos rl.Vector2) *Boid {
	return &Boid{
		Pos:           pos,
		Vel:           rl.NewVector2((rand.Float32()*2-1)*20, (rand.Float32()*2-1)*20),
		MaxVel:        10,
		CloseRadius:   20,
		Radius:        30,
		fov:           270 * math.Pi / 180,
		SeparationStr: 1.0 / initialSepMod,
		CohesionStr:   1.0 / initialCohMod,
		AlignmentStr:  1.0 / initialAlMod,
		once:          sync.Once{},
	}
}

const (
	initialSepMod = 1.25
	initialCohMod = 6
	initialAlMod  = 0.75
)

var sepMod = initialSepMod
var cohMod = initialCohMod
var alMod = initialAlMod

func (b *Boid) UpdateMove() {

	// Move
	b.Pos = rl.Vector2Add(b.Pos, b.Vel)
	// Add bounds check
	if b.Pos.X > float32(rl.GetScreenWidth()) {
		b.Pos.X = 0
	}

	if b.Pos.X < 0 {
		b.Pos.X = float32(rl.GetScreenWidth())
	}

	if b.Pos.Y > float32(rl.GetScreenHeight()) {
		b.Pos.Y = 0
	}

	if b.Pos.Y < 0 {
		b.Pos.Y = float32(rl.GetScreenHeight())
	}

	// Enforce min and max velocity

}

func (b *Boid) UpdateForces(boids []*Boid) {

	b.once.Do(func() {
		b.Vel = rl.NewVector2(rand.Float32()*2-1, rand.Float32()*2-1)
		b.MaxVel = 5
		b.SeparationStr = float32(1.0 / sepMod)
		b.CohesionStr = float32(1.0 / cohMod)
		b.AlignmentStr = float32(1.0 / alMod)
	})

	var separation rl.Vector2
	var alignment rl.Vector2
	var cohesion rl.Vector2
	var neighborCount int

	for _, other := range boids {
		if other == b {
			continue
		}

		distance := rl.Vector2Distance(b.Pos, other.Pos)

		if distance < b.CloseRadius {
			// Separation
			diff := rl.Vector2Subtract(b.Pos, other.Pos)
			diff = rl.Vector2Scale(diff, 1.0/distance) // Weight by distance
			separation = rl.Vector2Add(separation, diff)
		}

		if distance < b.Radius {
			// Alignment
			alignment = rl.Vector2Add(alignment, other.Vel)

			// Cohesion
			cohesion = rl.Vector2Add(cohesion, other.Pos)

			neighborCount++
		}
	}

	if neighborCount > 0 {
		// Normalize and apply strengths
		separation = rl.Vector2Scale(separation, b.SeparationStr)

		alignment = rl.Vector2Scale(alignment, 1.0/float32(neighborCount))
		alignment = rl.Vector2Subtract(alignment, b.Vel)
		alignment = rl.Vector2Scale(alignment, b.AlignmentStr)

		cohesion = rl.Vector2Scale(cohesion, 1.0/float32(neighborCount))
		cohesion = rl.Vector2Subtract(cohesion, b.Pos)
		cohesion = rl.Vector2Scale(cohesion, b.CohesionStr)

		// Apply forces
		b.Vel = rl.Vector2Add(b.Vel, separation)
		b.Vel = rl.Vector2Add(b.Vel, alignment)
		b.Vel = rl.Vector2Add(b.Vel, cohesion)

		randomAngle := (rand.Float32() - 0.5) * math.Pi * 0.5 * 0.11 // Random angle between -90 and 90 degrees
		b.Vel = rotateVector(b.Vel, randomAngle)
	}

	// Limit velocity

	minVel := float32(1) // Adjust this value as needed
	speed := rl.Vector2Length(b.Vel)
	if speed < minVel {
		b.Vel = rl.Vector2Scale(rl.Vector2Normalize(b.Vel), minVel/speed)
	}

	speed = rl.Vector2Length(b.Vel)
	if speed > b.MaxVel {
		b.Vel = rl.Vector2Scale(rl.Vector2Normalize(b.Vel), b.MaxVel/speed)
	}

}

func main() {
	initRl()

	for !rl.WindowShouldClose() {

		if rl.IsKeyDown(rl.KeyQ) {
			sepMod -= 1
		}

		if rl.IsKeyDown(rl.KeyA) {
			sepMod += 1
		}

		if rl.IsKeyDown(rl.KeyW) {
			cohMod -= 1
		}

		if rl.IsKeyDown(rl.KeyS) {
			cohMod += 1
		}

		if rl.IsKeyDown(rl.KeyE) {
			alMod -= 1
		}

		if rl.IsKeyDown(rl.KeyD) {
			alMod += 1
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		for _, b := range boids {
			b.UpdateForces(boids)
			b.UpdateMove()

			rl.DrawCircleV(b.Pos, 2, lerpColor(rl.White, rl.Black, rl.Vector2Length(b.Vel)/b.MaxVel))
			// rl.DrawLine(int32(b.Pos.X), int32(b.Pos.Y), int32(b.Pos.X+b.Vel.X), int32(b.Pos.Y+b.Vel.Y), lerpColor(rl.Black, rl.White, rl.Vector2Length(b.Vel)/b.MaxVel))
		}

		rl.EndDrawing()

	}
}

func initRl() {
	rl.InitWindow(1600, 900, "Boids")

	rl.SetTargetFPS(60)

	boids = make([]*Boid, 0)

	for i := 0; i < 300; i++ {

		randX := rand.IntN(rl.GetScreenWidth())
		randY := rand.IntN(rl.GetScreenHeight())

		boids = append(boids, NewBoid(rl.NewVector2(float32(randX), float32(randY))))
	}

}
func rotateVector(v rl.Vector2, angle float32) rl.Vector2 {
	cos := float32(math.Cos(float64(angle)))
	sin := float32(math.Sin(float64(angle)))
	return rl.Vector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

func lerpColor(a, b rl.Color, t float32) rl.Color {

	r := uint8(rl.Lerp(float32(a.R), float32(b.R), t))
	g := uint8(rl.Lerp(float32(a.G), float32(b.G), t))
	bl := uint8(rl.Lerp(float32(a.B), float32(b.B), t))
	al := uint8(rl.Lerp(float32(a.A), float32(b.A), t))
	return rl.Color{R: r, G: g, B: bl, A: al}

}
