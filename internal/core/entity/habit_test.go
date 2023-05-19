package entity

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCalculateHabitStatus(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		logs     []*HabitLog
		expected HabitStatus
	}{
		"Non-automated unsuccesful overrides automated succesful": {
			expected: HabitStatusNotDone,
			logs: []*HabitLog{
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
		"Non-automated succesful overrides automated unsuccesful": {
			expected: HabitStatusDone,
			logs: []*HabitLog{
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
		"Automated succesful": {
			expected: HabitStatusDone,
			logs: []*HabitLog{
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
		"Non-automated succesful": {
			expected: HabitStatusDone,
			logs: []*HabitLog{
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
		"Automated unsuccesful": {
			expected: HabitStatusDone,
			logs: []*HabitLog{
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
			expected: HabitStatusDone,
			logs: []*HabitLog{
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
			expected: HabitStatusNotDone,
			logs: []*HabitLog{
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
			expected: HabitStatusDone,
			logs: []*HabitLog{
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"One unsuccessful automated": {
			expected: HabitStatusNotDone,
			logs: []*HabitLog{
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"One successful nonautomated": {
			expected: HabitStatusDone,
			logs: []*HabitLog{
				{
					Success:     true,
					IsAutomated: false,
				},
			},
		},
		"One unsuccessful nonautomated": {
			expected: HabitStatusNotDone,
			logs: []*HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
			},
		},
		"No info with empty logs": {
			expected: HabitStatusNoInfo,
			logs:     []*HabitLog{},
		},
		"No info when logs are nil": {
			expected: HabitStatusNoInfo,
			logs:     nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			status := CalculateHabitStatus(tc.logs)
			assert.Equal(t, tc.expected, status)
		})
	}
}
