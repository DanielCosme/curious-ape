package core

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCalculateHabitStatus(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		logs     []HabitLog
		expected HabitState
	}{
		"Non-automated unsuccessful overrides automated successful": {
			expected: HabitStateNotDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
				{
					Success:     true,
					IsAutomated: true,
				},
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"Non-automated successful overrides automated unsuccessful": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     true,
					IsAutomated: false,
				},
				{
					Success:     false,
					IsAutomated: true,
				},
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"Automated successful": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: true,
				},
				{
					Success:     true,
					IsAutomated: true,
				},
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"Non-automated successful": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
				{
					Success:     false,
					IsAutomated: false,
				},
				{
					Success:     true,
					IsAutomated: false,
				},
			},
		},
		"Automated unsuccessful": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: true,
				},
				{
					Success:     false,
					IsAutomated: true,
				},
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"All successful": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     true,
					IsAutomated: false,
				},
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"All unsuccessful": {
			expected: HabitStateNotDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"One successful automated": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"One unsuccessful automated": {
			expected: HabitStateNotDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"One successful non-automated": {
			expected: HabitStateDone,
			logs: []HabitLog{
				{
					Success:     true,
					IsAutomated: false,
				},
			},
		},
		"One unsuccessful non-automated": {
			expected: HabitStateNotDone,
			logs: []HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
			},
		},
		"No info with empty logs": {
			expected: HabitStateNoInfo,
			logs:     []HabitLog{},
		},
		"No info when logs are nil": {
			expected: HabitStateNoInfo,
			logs:     nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			status, _ := calculateHabitState(tc.logs)
			assert.Equal(t, tc.expected, status)
		})
	}
}
