package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MichaelGenchev/smart-home-system/pkg/proto"
	"github.com/MichaelGenchev/smart-home-system/internal/gateway/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"

)

type UserHandler struct {
	client proto.UserServiceClient
}

func NewUserHandler(conn *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		client: proto.NewUserServiceClient(conn),
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users")

	switch{
	case path == "" || path == "/":
		switch r.Method {
		case http.MethodPost:
			h.CreateUser(w,r)
		default:
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	case strings.HasPrefix(path, "/"):
		id := strings.TrimPrefix(path, "/")
		switch r.Method {
		case http.MethodGet:
			h.GetUser(w,r,id)
		case http.MethodPut:
			h.UpdateUser(w,r,id)
		case http.MethodDelete:
			h.DeleteUser(w,r,id)
		default:
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	default:
		utils.RespondWithError(w, http.StatusNotFound, "Not found")
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	resp, err := h.client.CreateUser(r.Context(), &req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	resp, err := h.client.GetUser(r.Context(), &proto.GetUserRequest{Id: id})
	if err != nil {
		if st, ok := status.FromError(err); ok  && st.Code() == codes.NotFound {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, resp)
}	

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	var req proto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	req.Id = id

	user, err := h.client.UpdateUser(r.Context(), &req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code(){
			case codes.NotFound:
				utils.RespondWithError(w, http.StatusNotFound, "User not found")
			case codes.InvalidArgument:
				utils.RespondWithError(w, http.StatusBadRequest, st.Message())
			default:
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
			}
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError,"Failed to update user")
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	_, err := h.client.DeleteUser(r.Context(), &proto.DeleteUserRequest{Id:id})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message":"User deleted successfully"})
}

