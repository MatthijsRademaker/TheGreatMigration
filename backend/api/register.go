package api

import (
	"github.com/danielgtaylor/huma/v2"
)

// RegisterAll registers every API endpoint with the given Huma API instance.
//
// store may be nil when only the OpenAPI specification is needed (no database
// connection required). All registration functions only use store inside handler
// closures, which are never invoked during spec generation.
func RegisterAll(api huma.API, store Store) {
	registerHello(api)
	registerDashboardPeopleAvailability(api, store)
	registerPlanningWindow(api, store)
	registerTasksBacklog(api, store)
	registerTasksEndpoints(api, store)
	registerDailySchedule(api, store)
	registerScheduleCardEndpoints(api, store)
	registerPeopleEndpoints(api, store)
	registerRoomsAreas(api, store)
}
