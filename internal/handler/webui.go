package handler

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"git.arcline.it/ArclineIT/nexus/internal/auth"
	"git.arcline.it/ArclineIT/nexus/internal/config"
	"github.com/google/uuid"
)

//go:embed templates
var templateFS embed.FS

// UITemplateData holds data passed to UI templates.
type UITemplateData struct {
	Error       string
	Success     string
	Email       string
	DisplayName string
	Token       string
}

// UIHandler serves the web UI pages and handles HTMX form submissions.
type UIHandler struct {
	cfg       *config.Config
	templates *template.Template
}

// NewUIHandler creates a new UIHandler.
func NewUIHandler(cfg *config.Config) (*UIHandler, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &UIHandler{cfg: cfg, templates: tmpl}, nil
}

// ServePage renders a full page (base + named template).
func (h *UIHandler) ServePage(w http.ResponseWriter, data any, templateName string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
	}
}

// ServeFragment renders only the content fragment (for HTMX swaps).
func (h *UIHandler) ServeFragment(w http.ResponseWriter, data any, templateName string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "failed to render fragment", http.StatusInternalServerError)
	}
}

// Root redirects / to /login.
func (h *UIHandler) Root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

// LoginPage serves GET /login.
func (h *UIHandler) LoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &UITemplateData{}
		// Check for success message from signup
		if msg := r.URL.Query().Get("registered"); msg == "1" {
			data.Success = "Account created successfully. Please sign in."
		}
		if msg := r.URL.Query().Get("reset"); msg == "1" {
			data.Success = "Password reset successfully. Please sign in."
		}
		h.ServePage(w, data, "login")
	}
}

// LoginSubmit handles POST /login (HTMX form submission).
func (h *UIHandler) LoginSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		data := &UITemplateData{Email: email}

		if email == "" || password == "" {
			data.Error = "Email and password are required."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "login")
			return
		}

		// TODO: validate credentials against database
		// For now, accept any credentials and generate tokens
		userID := uuid.New()

		tokens, err := auth.GenerateTokenPair(h.cfg, userID, email)
		if err != nil {
			data.Error = "Something went wrong. Please try again."
			w.WriteHeader(http.StatusInternalServerError)
			h.ServeFragment(w, data, "login")
			return
		}

		// Set access token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "nexus_access_token",
			Value:    tokens.AccessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(h.cfg.Auth.AccessTokenDuration.Seconds()),
		})

		// Set refresh token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "nexus_refresh_token",
			Value:    tokens.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(h.cfg.Auth.RefreshTokenDuration.Seconds()),
		})

		// Tell HTMX to redirect to the dashboard
		w.Header().Set("HX-Redirect", "/dashboard")
		w.WriteHeader(http.StatusOK)
	}
}

// ---------------------------------------------------------------------------
// Signup
// ---------------------------------------------------------------------------

// SignupPage serves GET /signup.
func (h *UIHandler) SignupPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServePage(w, &UITemplateData{}, "signup")
	}
}

// SignupSubmit handles POST /signup (HTMX form submission).
func (h *UIHandler) SignupSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		displayName := strings.TrimSpace(r.FormValue("display_name"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirm")

		data := &UITemplateData{
			Email:       email,
			DisplayName: displayName,
		}

		// Validate
		if displayName == "" {
			data.Error = "Full name is required."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "signup")
			return
		}
		if email == "" {
			data.Error = "Email address is required."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "signup")
			return
		}
		if password == "" {
			data.Error = "Password is required."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "signup")
			return
		}
		if len(password) < 8 {
			data.Error = "Password must be at least 8 characters."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "signup")
			return
		}
		if password != passwordConfirm {
			data.Error = "Passwords do not match."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "signup")
			return
		}

		// TODO: check if email already exists in database
		// TODO: hash password with bcrypt and store user

		// Redirect to login with success message
		w.Header().Set("HX-Redirect", "/login?registered=1")
		w.WriteHeader(http.StatusOK)
	}
}

// ---------------------------------------------------------------------------
// Forgot Password
// ---------------------------------------------------------------------------

// ForgotPasswordPage serves GET /forgot-password.
func (h *UIHandler) ForgotPasswordPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServePage(w, &UITemplateData{}, "forgot-password")
	}
}

// ForgotPasswordSubmit handles POST /forgot-password (HTMX form submission).
func (h *UIHandler) ForgotPasswordSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := strings.TrimSpace(r.FormValue("email"))

		data := &UITemplateData{Email: email}

		if email == "" {
			data.Error = "Email address is required."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "forgot-password")
			return
		}

		// TODO: look up user in database, generate reset token, send email
		// For now, always show success to prevent email enumeration
		data.Success = "If an account exists for " + email + ", you will receive a password reset link shortly."

		h.ServeFragment(w, data, "forgot-password")
	}
}

// ---------------------------------------------------------------------------
// Reset Password
// ---------------------------------------------------------------------------

// ResetPasswordPage serves GET /reset-password.
func (h *UIHandler) ResetPasswordPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		data := &UITemplateData{Token: token}

		if token == "" {
			data.Error = "Invalid or missing reset token."
			h.ServePage(w, data, "reset-password")
			return
		}

		// TODO: validate reset token exists and hasn't expired
		h.ServePage(w, data, "reset-password")
	}
}

// ResetPasswordSubmit handles POST /reset-password (HTMX form submission).
func (h *UIHandler) ResetPasswordSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("password_confirm")

		data := &UITemplateData{Token: token}

		if token == "" {
			data.Error = "Invalid or missing reset token."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "reset-password")
			return
		}

		if password == "" {
			data.Error = "Password is required."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "reset-password")
			return
		}
		if len(password) < 8 {
			data.Error = "Password must be at least 8 characters."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "reset-password")
			return
		}
		if password != passwordConfirm {
			data.Error = "Passwords do not match."
			w.WriteHeader(http.StatusBadRequest)
			h.ServeFragment(w, data, "reset-password")
			return
		}

		// TODO: validate reset token, look up user, hash new password, save

		// Redirect to login with success message
		w.Header().Set("HX-Redirect", "/login?reset=1")
		w.WriteHeader(http.StatusOK)
	}
}

// ---------------------------------------------------------------------------
// Dashboard
// ---------------------------------------------------------------------------

// DashboardPage serves GET /dashboard — simple placeholder for now.
func (h *UIHandler) DashboardPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value("user_id").(string)
		userEmail, _ := r.Context().Value("user_email").(string)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Simple inline dashboard — can be moved to a template later
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Nexus — Dashboard</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <link rel="icon" type="image/svg+xml" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'><text y='28' font-size='28'>⚡</text></svg>">
</head>
<body class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
    <nav class="border-b border-slate-700/50 bg-slate-800/50 backdrop-blur-sm">
        <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
            <div class="flex items-center gap-3">
                <span class="text-2xl">⚡</span>
                <span class="text-white font-semibold text-lg">Nexus</span>
            </div>
            <div class="flex items-center gap-4">
                <span class="text-slate-400 text-sm">` + userEmail + `</span>
                <form hx-post="/logout" hx-target="body" class="inline">
                    <button type="submit" class="text-slate-400 hover:text-white text-sm transition-colors">Sign out</button>
                </form>
            </div>
        </div>
    </nav>
    <main class="max-w-6xl mx-auto px-4 py-12">
        <div class="bg-slate-800/50 backdrop-blur-sm border border-slate-700/50 rounded-xl shadow-2xl p-8">
            <h2 class="text-2xl font-bold text-white mb-2">Welcome back</h2>
            <p class="text-slate-400 mb-6">You are signed in as <span class="text-white font-medium">` + userEmail + `</span>.</p>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="bg-slate-700/50 rounded-lg p-4 border border-slate-600/50">
                    <div class="text-slate-400 text-sm mb-1">User ID</div>
                    <div class="text-white font-mono text-sm break-all">` + userID + `</div>
                </div>
                <div class="bg-slate-700/50 rounded-lg p-4 border border-slate-600/50">
                    <div class="text-slate-400 text-sm mb-1">Connected Apps</div>
                    <div class="text-white text-sm">None yet</div>
                </div>
            </div>
        </div>
    </main>
</body>
</html>`))
	}
}

// ---------------------------------------------------------------------------
// Logout
// ---------------------------------------------------------------------------

// Logout handles POST /logout — clears auth cookies and redirects to login.
func (h *UIHandler) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Clear cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "nexus_access_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "nexus_refresh_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusOK)
	}
}

