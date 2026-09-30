package views

import (
	"time"
)

type workoutSet struct {
	id        string
	createdAt time.Time
	reps      int
	weight    float64
	rest      time.Duration
}

type workoutSession struct {
	date time.Time
	sets []workoutSet
}

type exerciseRecord struct {
	id              string
	name            string
	bodyWeightShare float32
	sessions        []workoutSession
}
