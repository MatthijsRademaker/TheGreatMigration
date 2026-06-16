package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ---------- Domain types ----------

// Room is the API-facing room/area record.
type Room struct {
	ID        string `json:"id" doc:"Stable room identifier, prefixed 'room-'"`
	Name      string `json:"name" doc:"Human-readable room or area name"`
	Type      string `json:"type" enum:"room,area" doc:"One of: room, area"`
	CreatedAt string `json:"createdAt" doc:"ISO 8601 timestamp of creation"`
	UpdatedAt string `json:"updatedAt" doc:"ISO 8601 timestamp of last update"`
}

// CreateRoomInput is the domain-layer input for creating a room.
type CreateRoomInput struct {
	Name string
	Type string
}

// UpdateRoomInput is the domain-layer input for updating a room.
type UpdateRoomInput struct {
	Name string
	Type string
}

// ErrRoomNotFound is returned when a room ID is not found.
var ErrRoomNotFound = errors.New("room not found")

// ---------- Request / Response types ----------

// ListRoomsInput is empty: the endpoint accepts no query parameters.
type ListRoomsInput struct{}

// ListRoomsOutput is the response wrapper for GET /api/rooms.
type ListRoomsOutput struct {
	Body struct {
		Rooms []Room `json:"rooms" doc:"All room/area records"`
	}
}

// CreateRoomRequestBody holds the fields for creating a room.
type CreateRoomRequestBody struct {
	Name string `json:"name" required:"true" doc:"Human-readable room or area name"`
	Type string `json:"type" required:"true" enum:"room,area" doc:"One of: room, area"`
}

// CreateRoomInputHuma is the input for POST /api/rooms.
type CreateRoomInputHuma struct {
	Body CreateRoomRequestBody
}

// CreateRoomOutput is the response wrapper for POST /api/rooms.
type CreateRoomOutput struct {
	Body Room
}

// UpdateRoomRequestBody holds the fields for updating a room.
type UpdateRoomRequestBody struct {
	Name string `json:"name" required:"true" doc:"Updated human-readable room or area name"`
	Type string `json:"type" required:"true" enum:"room,area" doc:"One of: room, area"`
}

// UpdateRoomInputHuma is the input for PUT /api/rooms/{id}.
type UpdateRoomInputHuma struct {
	ID   string `path:"id" doc:"Room identifier"`
	Body UpdateRoomRequestBody
}

// UpdateRoomOutput is the response wrapper for PUT /api/rooms/{id}.
type UpdateRoomOutput struct {
	Body Room
}

// DeleteRoomInput is the input for DELETE /api/rooms/{id}.
type DeleteRoomInput struct {
	ID string `path:"id" doc:"Room identifier"`
}

// DeleteRoomOutput is the empty response for DELETE /api/rooms/{id}.
type DeleteRoomOutput struct{}

// ---------- Handler ----------

func registerRoomsAreas(api huma.API, store Store) {
	// GET /api/rooms — list all rooms.
	huma.Register(api, huma.Operation{
		OperationID: "list-rooms",
		Method:      http.MethodGet,
		Path:        "/api/rooms",
		Summary:     "List rooms and areas",
		Description: "Returns all room and area records ordered by name.",
		Tags:        []string{"Rooms"},
	}, func(ctx context.Context, input *ListRoomsInput) (*ListRoomsOutput, error) {
		rooms, err := store.ListRooms(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list rooms", err)
		}
		resp := &ListRoomsOutput{}
		resp.Body.Rooms = rooms
		return resp, nil
	})

	// POST /api/rooms — create a new room.
	huma.Register(api, huma.Operation{
		OperationID:   "create-room",
		Method:        http.MethodPost,
		Path:          "/api/rooms",
		Summary:       "Create a room or area",
		Description:   "Creates a new room or area record and returns it.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Rooms"},
	}, func(ctx context.Context, input *CreateRoomInputHuma) (*CreateRoomOutput, error) {
		room, err := store.CreateRoom(ctx, CreateRoomInput{
			Name: input.Body.Name,
			Type: input.Body.Type,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create room", err)
		}
		resp := &CreateRoomOutput{}
		resp.Body = *room
		return resp, nil
	})

	// PUT /api/rooms/{id} — update an existing room.
	huma.Register(api, huma.Operation{
		OperationID: "update-room",
		Method:      http.MethodPut,
		Path:        "/api/rooms/{id}",
		Summary:     "Update a room or area",
		Description: "Updates an existing room or area record by ID. Returns 404 if the ID is not found.",
		Tags:        []string{"Rooms"},
	}, func(ctx context.Context, input *UpdateRoomInputHuma) (*UpdateRoomOutput, error) {
		room, err := store.UpdateRoom(ctx, input.ID, UpdateRoomInput{
			Name: input.Body.Name,
			Type: input.Body.Type,
		})
		if err != nil {
			if errors.Is(err, ErrRoomNotFound) {
				return nil, huma.Error404NotFound("room not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to update room", err)
		}
		resp := &UpdateRoomOutput{}
		resp.Body = *room
		return resp, nil
	})

	// DELETE /api/rooms/{id} — delete a room.
	huma.Register(api, huma.Operation{
		OperationID: "delete-room",
		Method:      http.MethodDelete,
		Path:        "/api/rooms/{id}",
		Summary:     "Delete a room or area",
		Description: "Deletes a room or area record by ID. Returns 404 if the ID is not found.",
		Tags:        []string{"Rooms"},
	}, func(ctx context.Context, input *DeleteRoomInput) (*DeleteRoomOutput, error) {
		err := store.DeleteRoom(ctx, input.ID)
		if err != nil {
			if errors.Is(err, ErrRoomNotFound) {
				return nil, huma.Error404NotFound("room not found", err)
			}
			return nil, huma.Error500InternalServerError("failed to delete room", err)
		}
		return &DeleteRoomOutput{}, nil
	})
}
