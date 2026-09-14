// Package api wires HTTP routes to IAM service operations.
package api

import (
	"net/http"
	"time"

	"github.com/jasonmiller-cc/parameters-core/pkg/auth"
	apierrors "github.com/jasonmiller-cc/parameters-core/pkg/errors"
	"github.com/jasonmiller-cc/parameters-core/pkg/response"
	"github.com/jasonmiller-cc/parameters-iam/internal/service"
)

// Handler dispatches HTTP requests to the IAMService.
type Handler struct {
	svc *service.IAMService
	jwt *auth.JWTConfig
}

// New creates a Handler backed by svc and using jwt for token operations.
func New(svc *service.IAMService, jwt *auth.JWTConfig) *Handler {
	return &Handler{svc: svc, jwt: jwt}
}

// Register attaches all /api/v1 routes to mux.
// Uses Go 1.22 method+pattern syntax so method dispatch is handled by the mux.
func (h *Handler) Register(mux *http.ServeMux) {
	// Users
	mux.HandleFunc("POST /api/v1/users", h.provisionUser)
	mux.HandleFunc("GET /api/v1/users", h.listUsers)
	mux.HandleFunc("GET /api/v1/users/{uid}", h.getUser)
	mux.HandleFunc("DELETE /api/v1/users/{uid}", h.deprovisionUser)
	mux.HandleFunc("POST /api/v1/users/{uid}/roles", h.assignRoles)
	mux.HandleFunc("DELETE /api/v1/users/{uid}/roles/{role}", h.removeRole)
	// Roles
	mux.HandleFunc("GET /api/v1/roles", h.listRoles)
	mux.HandleFunc("POST /api/v1/roles", h.createRole)
	mux.HandleFunc("DELETE /api/v1/roles/{name}", h.deleteRole)
	mux.HandleFunc("GET /api/v1/roles/{name}/permissions", h.getRolePermissions)
	mux.HandleFunc("PUT /api/v1/roles/{name}/permissions", h.setRolePermissions)
	// Auth
	mux.HandleFunc("POST /api/v1/auth/token", h.issueToken)
	mux.HandleFunc("POST /api/v1/auth/refresh", h.refreshToken)
	mux.HandleFunc("POST /api/v1/auth/revoke", h.revokeToken)
}

// --- Users ---

func (h *Handler) provisionUser(w http.ResponseWriter, r *http.Request) {
	var req service.ProvisionRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Err(w, err)
		return
	}
	user, err := h.svc.ProvisionUser(r.Context(), req)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.Created(w, user)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListUsers(r.Context())
	if err != nil {
		response.Err(w, err)
		return
	}
	response.List(w, users, len(users), len(users), 0)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	user, err := h.svc.GetUser(r.Context(), uid)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.OK(w, user)
}

func (h *Handler) deprovisionUser(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	if err := h.svc.DeprovisionUser(r.Context(), uid); err != nil {
		response.Err(w, err)
		return
	}
	response.NoContent(w)
}

func (h *Handler) assignRoles(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	var body struct {
		Roles []string `json:"roles"`
	}
	if err := response.DecodeJSON(r, &body); err != nil {
		response.Err(w, err)
		return
	}
	user, err := h.svc.AssignRoles(r.Context(), uid, body.Roles)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.OK(w, user)
}

func (h *Handler) removeRole(w http.ResponseWriter, r *http.Request) {
	uid := r.PathValue("uid")
	role := r.PathValue("role")
	if err := h.svc.RemoveRole(r.Context(), uid, role); err != nil {
		response.Err(w, err)
		return
	}
	response.NoContent(w)
}

// --- Roles ---

func (h *Handler) listRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.ListRoles(r.Context())
	if err != nil {
		response.Err(w, err)
		return
	}
	response.List(w, roles, len(roles), len(roles), 0)
}

func (h *Handler) createRole(w http.ResponseWriter, r *http.Request) {
	var role service.Role
	if err := response.DecodeJSON(r, &role); err != nil {
		response.Err(w, err)
		return
	}
	created, err := h.svc.CreateRole(r.Context(), role)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.Created(w, created)
}

func (h *Handler) deleteRole(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := h.svc.DeleteRole(r.Context(), name); err != nil {
		response.Err(w, err)
		return
	}
	response.NoContent(w)
}

func (h *Handler) getRolePermissions(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	perms, err := h.svc.GetRolePermissions(r.Context(), name)
	if err != nil {
		response.Err(w, err)
		return
	}
	response.OK(w, perms)
}

func (h *Handler) setRolePermissions(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var perms []service.Permission
	if err := response.DecodeJSON(r, &perms); err != nil {
		response.Err(w, err)
		return
	}
	if err := h.svc.SetRolePermissions(r.Context(), name, perms); err != nil {
		response.Err(w, err)
		return
	}
	response.NoContent(w)
}

// --- Auth ---

// issueToken authenticates username/password via Kerberos (stub) and returns
// a signed JWT access token and refresh token using parameters-core's JWTConfig.
func (h *Handler) issueToken(w http.ResponseWriter, r *http.Request) {
	var req service.TokenRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Err(w, err)
		return
	}
	user, err := h.svc.ValidateCredentials(r.Context(), req.Username, req.Password)
	if err != nil {
		response.Err(w, err)
		return
	}

	accessClaims := auth.Claims{
		Subject: user.UID,
		Email:   user.Email,
		Roles:   user.Roles,
		Service: "parameters-iam",
	}
	accessToken, err := h.jwt.Sign(accessClaims)
	if err != nil {
		response.Err(w, apierrors.Internal(err))
		return
	}

	// Refresh token carries only identity, no role claims, with the same TTL.
	// A production implementation would use a longer TTL and a separate secret.
	refreshClaims := auth.Claims{
		Subject: user.UID,
		Service: "parameters-iam",
	}
	refreshToken, err := h.jwt.Sign(refreshClaims)
	if err != nil {
		response.Err(w, apierrors.Internal(err))
		return
	}

	ttl := h.jwt.TTL
	if ttl == 0 {
		ttl = 24 * time.Hour
	}
	response.OK(w, service.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    ttl.String(),
	})
}

// refreshToken validates a refresh JWT and issues a new access token.
func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := response.DecodeJSON(r, &body); err != nil {
		response.Err(w, err)
		return
	}
	if body.RefreshToken == "" {
		response.Err(w, apierrors.BadRequest("refresh_token is required"))
		return
	}
	claims, err := h.jwt.Verify(body.RefreshToken)
	if err != nil {
		response.Err(w, err)
		return
	}
	newClaims := auth.Claims{
		Subject: claims.Subject,
		Email:   claims.Email,
		Roles:   claims.Roles,
		Service: "parameters-iam",
	}
	accessToken, err := h.jwt.Sign(newClaims)
	if err != nil {
		response.Err(w, apierrors.Internal(err))
		return
	}
	ttl := h.jwt.TTL
	if ttl == 0 {
		ttl = 24 * time.Hour
	}
	response.OK(w, service.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   ttl.String(),
	})
}

// revokeToken records a token as revoked.
// A full implementation persists the JTI to a blocklist checked by JWTMiddleware.
func (h *Handler) revokeToken(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	if err := response.DecodeJSON(r, &body); err != nil {
		response.Err(w, err)
		return
	}
	// TODO: add token JTI to a persistent blocklist.
	response.NoContent(w)
}
