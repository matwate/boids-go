package main

import (
	"log"
	"math/rand"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type BoidParams struct {
	SeparationStr float32
	CohesionStr   float32
	AlignmentStr  float32
	CloseRadius   float32
	Radius        float32
	MaxVel        float32
}

func RandomParams() BoidParams {
	return BoidParams{
		SeparationStr: rand.Float32() * 0.1,
		CohesionStr:   rand.Float32() * 0.1,
		AlignmentStr:  rand.Float32() * 0.1,
		CloseRadius:   rand.Float32()*30 + 10,
		Radius:        rand.Float32()*50 + 30,
		MaxVel:        rand.Float32()*4 + 1,
	}
}

func Fitness(params BoidParams) float32 {
	// Create a set of boids with these parameters
	boids := make([]*Boid, 50)
	for i := range boids {
		boids[i] = NewBoid(rl.NewVector2(rand.Float32()*float32(rl.GetScreenWidth()), rand.Float32()*float32(rl.GetScreenHeight())))
		boids[i].SeparationStr = params.SeparationStr
		boids[i].CohesionStr = params.CohesionStr
		boids[i].AlignmentStr = params.AlignmentStr
		boids[i].CloseRadius = params.CloseRadius
		boids[i].Radius = params.Radius
		boids[i].MaxVel = params.MaxVel
	}

	// Run simulation for a number of steps
	for step := 0; step < 1000; step++ {
		for _, b := range boids {
			b.UpdateForces(boids)
			b.UpdateMove()
		}
	}

	// Calculate average speed
	totalSpeed := float32(0)
	for _, b := range boids {
		totalSpeed += rl.Vector2Length(b.Vel)
	}
	return totalSpeed / float32(len(boids))
}

func Evolve(generations int, populationSize int) BoidParams {
	population := make([]BoidParams, populationSize)
	for i := range population {
		population[i] = RandomParams()
	}

	for gen := 0; gen < generations; gen++ {
		// Evaluate fitness
		fitnesses := make([]float32, populationSize)
		for i, params := range population {
			fitnesses[i] = Fitness(params)
		}

		// Sort population by fitness
		sort.Slice(population, func(i, j int) bool {
			return fitnesses[i] > fitnesses[j]
		})

		// Keep top half, replace bottom half with children
		for i := populationSize / 2; i < populationSize; i++ {
			parent1 := population[rand.Intn(populationSize/2)]
			parent2 := population[rand.Intn(populationSize/2)]
			child := BoidParams{
				SeparationStr: (parent1.SeparationStr + parent2.SeparationStr) / 2,
				CohesionStr:   (parent1.CohesionStr + parent2.CohesionStr) / 2,
				AlignmentStr:  (parent1.AlignmentStr + parent2.AlignmentStr) / 2,
				CloseRadius:   (parent1.CloseRadius + parent2.CloseRadius) / 2,
				Radius:        (parent1.Radius + parent2.Radius) / 2,
				MaxVel:        (parent1.MaxVel + parent2.MaxVel) / 2,
			}
			// Mutate
			if rand.Float32() < 0.1 {
				child.SeparationStr *= rand.Float32()*0.5 + 0.75
			}
			// ... (similar mutations for other parameters)
			population[i] = child
			// Print
			log.Printf("Generation %d, child %d: %v", gen, i, child)
		}

		// Print best fitness
		log.Printf("Best fitness in generation %d: %f", gen, fitnesses[0])
	}

	// Return best params
	return population[0]
}

func main_() {
	bestParams := Evolve(50, 100)
	// Use bestParams to set up your boids
	// ...
	log.Println(bestParams)
}
