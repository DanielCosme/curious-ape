package entity_test

import (
	"testing"

	"github.com/danielcosme/curious-ape/internal/core/entity"
	"gotest.tools/v3/assert"
)

func TestCalculateHabitStatus(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		logs     []*entity.HabitLog
		expected entity.HabitStatus
	}{
		"Non-automated unsuccesful overrides automated succesful": {
			expected: entity.HabitStatusNotDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusNotDone,
			logs: []*entity.HabitLog{
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
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
				{
					Success:     true,
					IsAutomated: true,
				},
			},
		},
		"One unsuccessful automated": {
			expected: entity.HabitStatusNotDone,
			logs: []*entity.HabitLog{
				{
					Success:     false,
					IsAutomated: true,
				},
			},
		},
		"One successful nonautomated": {
			expected: entity.HabitStatusDone,
			logs: []*entity.HabitLog{
				{
					Success:     true,
					IsAutomated: false,
				},
			},
		},
		"One unsuccessful nonautomated": {
			expected: entity.HabitStatusNotDone,
			logs: []*entity.HabitLog{
				{
					Success:     false,
					IsAutomated: false,
				},
			},
		},
		"No info with empty logs": {
			expected: entity.HabitStatusNoInfo,
			logs:     []*entity.HabitLog{},
		},
		"No info when logs are nil": {
			expected: entity.HabitStatusNoInfo,
			logs:     nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			status := entity.CalculateHabitStatus(tc.logs)
			assert.Equal(t, tc.expected, status)
		})
	}
}
