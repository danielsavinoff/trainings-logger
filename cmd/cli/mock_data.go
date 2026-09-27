package main

import "time"

type workoutSet struct {
	reps     int
	weightKG float64
	rest     time.Duration
}

type workoutSession struct {
	date time.Time
	sets []workoutSet
}

type mockExercise struct {
	name         string
	bodyWeightKG float64
	sessions     []workoutSession
}

func mockExercises() []mockExercise {
	// Recent working sets for an example lifter weighing 80 kg.
	fixtures := []struct {
		name   string
		weight float64
		reps   int
		rest   time.Duration
		day    int
	}{
		{"back squat", 90, 5, 3 * time.Minute, 23},
		{"bench press", 65, 8, 2*time.Minute + 30*time.Second, 22},
		{"deadlift", 110, 5, 3 * time.Minute, 21},
		{"overhead press", 35, 8, 2 * time.Minute, 22},
		{"barbell row", 55, 10, 90 * time.Second, 21},
		{"romanian deadlift", 75, 8, 2 * time.Minute, 23},
		{"lat pulldown", 50, 10, 90 * time.Second, 21},
		{"leg press", 140, 12, 2 * time.Minute, 23},
	}
	exercises := make([]mockExercise, 0, len(fixtures))
	for _, fixture := range fixtures {
		exercise := mockExercise{name: fixture.name, bodyWeightKG: 80}
		for week := 0; week < 3; week++ {
			weight := fixture.weight - float64(week)*2.5
			exercise.sessions = append(exercise.sessions, workoutSession{
				date: time.Date(2026, time.September, fixture.day-week*7, 0, 0, 0, 0, time.UTC),
				sets: []workoutSet{
					{fixture.reps, weight, fixture.rest},
					{fixture.reps, weight, fixture.rest},
					{fixture.reps - 1, weight, fixture.rest + 30*time.Second},
				},
			})
		}
		exercises = append(exercises, exercise)
	}
	return exercises
}
