package core_test

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCalculateHabitStatus(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		logs     []*core.HabitLog
		expected core.HabitState
	}{
		"Non-automated unsuccessful overrides automated successful": {
			expected: core.HabitStateNotDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateNotDone,
			logs: []*core.HabitLog{
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
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"One unsuccessful automated": {
			expected: core.HabitStateNotDone,
			logs: []*core.HabitLog{
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"One successful non-automated": {
			expected: core.HabitStateDone,
			logs: []*core.HabitLog{
				{
					Success:     true,
					IsAutomated: false,
				},
			},
		},
		"One unsuccessful non-automated": {
			expected: core.HabitStateNotDone,
			logs: []*core.HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
			},
		},
		"No info with empty logs": {
			expected: core.HabitStateNoInfo,
			logs:     []*core.HabitLog{},
		},
		"No info when logs are nil": {
			expected: core.HabitStateNoInfo,
			logs:     nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			status := core.CalculateHabitStatus(tc.logs)
			assert.Equal(t, tc.expected, status)
		})
	}
}
