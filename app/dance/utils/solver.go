package utils

import (
	"github.com/wieku/danser-go/app/beatmap/objects"
	"sort"
)

func Solve2B(queue []objects.IHitObject) []objects.IHitObject {
	// Resolving 2B conflicts
	for i := 0; i < len(queue); i++ {
		if s, ok := queue[i].(*objects.Slider); ok {
			found := false

			// We need to loop backwards to look for overlapping spinners (p) that are separated by circles:
			// --ppppppppppppppppp------
			// ----------c--c-----------
			// ---------------ssssssss--
			// Looking just by i-1 (like i+1 in forward detection) wouldn't detect that scenario because objects
			// are not sorted by end times
			for j := i - 1; j >= 0; j-- {
				if o := queue[i-1]; o.GetEndTime() >= s.GetStartTime() {
					queue = PreprocessQueue(i, queue, true)
					found = true
					break
				}
			}

			// If no conflict was detected in the past then look one object ahead, no looping is needed in this scenario
			if !found && i+1 < len(queue) {
				if o := queue[i+1]; o.GetStartTime() <= s.GetEndTime() {
					queue = PreprocessQueue(i, queue, true)
				}
			}
		}
	}

	// Second 2B pass for spinners
	for i := 0; i < len(queue); i++ {
		if s, ok := queue[i].(*objects.Spinner); ok {
			var subSpinners []objects.IHitObject

			startTime := s.GetStartTime()

			for j := i + 1; j < len(queue); j++ {
				o := queue[j]

				if o.GetStartTime() >= s.GetEndTime() {
					break
				}

				if endTime := o.GetStartTime() - 30; endTime > startTime {
					subSpinners = append(subSpinners, objects.NewDummySpinner(startTime, endTime))
				}

				startTime = o.GetEndTime() + 30
			}

			if subSpinners != nil && len(subSpinners) > 0 {
				if s.GetEndTime() > startTime {
					subSpinners = append(subSpinners, objects.NewDummySpinner(startTime, s.GetEndTime()))
				}

				queue = append(queue[:i], append(subSpinners, queue[i+1:]...)...)
				sort.SliceStable(queue, func(i, j int) bool { return queue[i].GetStartTime() < queue[j].GetStartTime() })
			}
		}
	}

	return queue
}
