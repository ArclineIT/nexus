Rails.application.routes.draw do
  # OpenID Connect discovery / userinfo / JWKS endpoints.
  use_doorkeeper_openid_connect
  # OAuth2 authorization, token, introspection and revocation endpoints.
  # Applications are managed through the Nexus admin UI, not Doorkeeper's.
  use_doorkeeper do
    skip_controllers :applications, :authorized_applications
  end

  root "dashboard#show"

  # Authentication
  get    "login",  to: "sessions#new",     as: :new_session
  post   "login",  to: "sessions#create",  as: :session
  delete "logout", to: "sessions#destroy", as: :logout

  # Registration
  get  "signup", to: "registrations#new",    as: :new_registration
  post "signup", to: "registrations#create", as: :registration

  # Password reset
  get   "forgot-password",        to: "passwords#new",    as: :new_password
  post  "forgot-password",        to: "passwords#create", as: :passwords
  get   "reset-password/:token",  to: "passwords#edit",   as: :edit_password
  patch "reset-password/:token",  to: "passwords#update", as: :password
  put   "reset-password/:token",  to: "passwords#update"

  get "dashboard", to: "dashboard#show", as: :dashboard

  namespace :admin do
    root "connected_apps#index"
    resources :connected_apps, only: %i[ index show edit update ] do
      post :regenerate_secret, on: :member
    end
    resources :users, only: %i[ index update ]
  end

  # Health check for load balancers and uptime monitors.
  get "up" => "rails/health#show", as: :rails_health_check
end
