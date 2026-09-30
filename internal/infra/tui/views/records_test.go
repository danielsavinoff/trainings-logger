package views

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/testutils"
)

func sampleModel() RecordsModel {
	date := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	records := []exerciseRecord{{name: "Squat", bodyWeightShare: 0.75}, {name: "Bench press"}}
	for i := 0; i < 24; i++ {
		records[0].sessions = append(records[0].sessions, workoutSession{
			date: date.AddDate(0, 0, -i),
			sets: []workoutSet{{reps: 8, weight: 80, rest: 150 * time.Second}},
		})
	}

	return newRecordsModel(records)
}

func press(m RecordsModel, keys ...string) RecordsModel {
	for _, key := range keys {
		updated, _ := m.Update(testutils.Key(key))
		m = updated.(RecordsModel)
	}

	return m
}

// These snapshots were captured before separating views and components.
// Include the native cursor as well as the rendered screen in the comparison.
func TestRecordsRendering(t *testing.T) {
	for _, test := range []struct {
		name          string
		width, height int
		keys          []string
	}{
		{"records", 100, 24, nil},
		{"wide", 161, 40, nil},
		{"compact", 40, 12, nil},
		{"exercise_dialog", 100, 24, []string{"e"}},
		{"insert", 100, 24, []string{"e", "i"}},
		{"compact_dialog", 40, 12, []string{"e", "j", "i"}},
		{"set_dialog", 100, 24, []string{"g", "s", "e"}},
		{"pagination", 100, 24, []string{"g", "s", "l"}},
		{"unicode_input", 100, 24, []string{"n", "i", "深蹲", "é"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			updated, _ := sampleModel().Update(tea.WindowSizeMsg{Width: test.width, Height: test.height})
			m := press(updated.(RecordsModel), test.keys...)
			v := m.View()
			actual, err := json.MarshalIndent(struct {
				Content   string
				Cursor    *tea.Cursor
				AltScreen bool
				MouseMode tea.MouseMode
			}{v.Content, v.Cursor, v.AltScreen, v.MouseMode}, "", "  ")
			if err != nil {
				t.Fatal(err)
			}

			path := filepath.Join("testdata", test.name+".golden.json")
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err := os.MkdirAll("testdata", 0755); err != nil {
					t.Fatal(err)
				}

				if err := os.WriteFile(path, actual, 0644); err != nil {
					t.Fatal(err)
				}
			}

			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			if string(actual) != string(want) {
				t.Fatalf("rendering differs from %s", path)
			}
		})
	}
}
