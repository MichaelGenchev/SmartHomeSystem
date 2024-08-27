package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MichaelGenchev/smart-home-system/pkg/proto"
	"github.com/MichaelGenchev/smart-home-system/internal/gateway/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type DeviceHandler struct {
	client proto.DeviceServiceClient 
}


func NewDeviceHandler(conn *grpc.ClientConn) *DeviceHandler {
	return &DeviceHandler{
		client: proto.NewDeviceServiceClient(conn),
	}
}


func (h *DeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/devices")

	switch{
	case path == "" || path == "/":
		switch r.Method {
		case http.MethodGet:
			h.GetDevices(w,r)
		case http.MethodPost:
			h.CreateDevice(w,r)
		default:
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	case strings.HasPrefix(path, "/"):
		id := strings.TrimPrefix(path, "/")
		switch r.Method {
		case http.MethodGet:
			h.GetDevice(w,r,id)
		case http.MethodPut:
			h.UpdateDeviceState(w,r,id)
		default: 
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	default:
		utils.RespondWithError(w, http.StatusNotFound, "Not found")
	}
}


func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	var req proto.ListDevicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	resp, err := h.client.ListDevices(r.Context(), &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get devices")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}


func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	device, err := h.client.CreateDevice(r.Context(), &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create device")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, device)
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request, id string) {
	device, err := h.client.GetDevice(r.Context(), &proto.GetDeviceRequest{Id: id})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get device")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, device)
}

func (h *DeviceHandler) UpdateDeviceState(w http.ResponseWriter, r *http.Request, id string) {
	var req proto.UpdateDeviceStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	req.Id = id

	device, err := h.client.UpdateDeviceState(r.Context(), &req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				utils.RespondWithError(w, http.StatusNotFound, "Device not found")
			case codes.InvalidArgument:
				utils.RespondWithError(w, http.StatusBadRequest, st.Message())
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update device")
			}
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update device")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, device)
}