package service

import (
	"github.com/ts-dmitry/cronpad/backend/repository"
	"github.com/ts-dmitry/cronpad/backend/utils"
	"sort"
	"time"
)

func AddEventProperly(event repository.Event, existedEvents []repository.Event, uuidProvider utils.UuidProvider) []repository.Event {
	if len(existedEvents) == 0 {
		return []repository.Event{event}
	}

	eventBlocks := createEventBlocks(existedEvents, event)

	result := make([]repository.Event, 0)

	currentEvent := event

	for i := range eventBlocks {
		eventBlock := eventBlocks[i]

		if beforeOrEquals(eventBlock.end, currentEvent.Start) {
			// no conflicts: [--] (--)
			result = append(result, eventBlock.events...)

			continue
		}

		if beforeOrEquals(eventBlock.start, currentEvent.Start) && afterOrEquals(eventBlock.end, currentEvent.End) {
			// new event is inside an existed event block: [--(--)--]
			return existedEvents
		}

		if eventBlock.end.Before(currentEvent.End) && eventBlock.start.After(currentEvent.Start) {
			// event block is inside new event: (--[--]--)
			currentEvent.End = eventBlock.start
			result = append(result, currentEvent)

			currentEvent = event.Copy()
			currentEvent.ID = uuidProvider.New()
			currentEvent.Start = eventBlock.end

			result = append(result, eventBlock.events...)

			continue
		}

		if beforeOrEquals(eventBlock.start, currentEvent.Start) && eventBlock.end.After(currentEvent.Start) {
			// only start of event is inside an event block: [--(--]--)
			result = append(result, eventBlock.events...)
			currentEvent.Start = eventBlock.end
			continue
		}

		if currentEvent.Start.Before(eventBlock.start) && afterOrEquals(currentEvent.End, eventBlock.start) {
			// only end of event is inside event block: (--[--)--]
			currentEvent.End = eventBlock.start
			result = append(result, currentEvent)

			for j := i; j < len(eventBlocks); j++ {
				result = append(result, eventBlocks[j].events...)
			}

			return result
		}

		if afterOrEquals(eventBlock.start, currentEvent.End) {
			// no conflicts: (--) [--]
			result = append(result, currentEvent)

			for j := i; j < len(eventBlocks); j++ {
				result = append(result, eventBlocks[j].events...)
			}

			return result
		}
	}

	result = append(result, currentEvent)
	return result
}

type eventBlock struct {
	start  time.Time
	end    time.Time
	events []repository.Event
}

func createEventBlocks(events []repository.Event, targetEvent repository.Event) []eventBlock {
	sortedEvents := SortEventsByStartDate(events)

	currentEventBlock := createEventBlockFromEvent(events[0])
	result := make([]eventBlock, 0)

	for i := 1; i < len(sortedEvents); i++ {
		event := events[i]

		if event.ID == targetEvent.ID {
			continue
		}

		if currentEventBlock.end.Equal(event.Start) {
			currentEventBlock.end = event.End
			currentEventBlock.events = append(currentEventBlock.events, event)
		} else {
			result = append(result, currentEventBlock)
			currentEventBlock = createEventBlockFromEvent(event)
		}
	}
	result = append(result, currentEventBlock)

	return result
}

func createEventBlockFromEvent(event repository.Event) eventBlock {
	return eventBlock{
		start:  event.Start,
		end:    event.End,
		events: []repository.Event{event},
	}
}

func SortEventsByStartDate(events []repository.Event) []repository.Event {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})
	return events
}

func beforeOrEquals(comparable time.Time, reference time.Time) bool {
	return !comparable.After(reference)
}

func afterOrEquals(comparable time.Time, reference time.Time) bool {
	return !comparable.Before(reference)
}
