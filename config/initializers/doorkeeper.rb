# frozen_string_literal: true

Doorkeeper.configure do
  orm :active_record

  # Use the domain model so OAuth clients are "Connected Apps" in the UI.
  application_class "ConnectedApp"

  # The resource owner is whoever holds a valid Nexus browser session. If
  # there is no session we bounce through the login page and come back.
  resource_owner_authenticator do
    Session.find_by(id: cookies.signed[:session_id])&.user || begin
      session[:return_to_after_authenticating] = request.fullpath
      redirect_to new_session_url
    end
  end

  # Application management is provided by the Nexus admin UI; Doorkeeper's
  # own admin controllers are skipped in config/routes.rb.
  admin_authenticator do
    redirect_to root_url
  end

  realm "Nexus"

  # Authorization Code (browser SSO), Refresh Token and Client Credentials
  # (service-to-service) grants. Implicit and password grants are disabled.
  grant_flows %w[authorization_code refresh_token client_credentials]

  default_scopes :openid
  optional_scopes :profile, :email, :roles, :offline_access

  access_token_expires_in 2.hours
  authorization_code_expires_in 10.minutes
  use_refresh_token

  # First-party Arcline apps are trusted, so skip the per-app consent screen
  # and deliver the "one login, every tool" experience.
  skip_authorization do
    true
  end

  # Require PKCE for public (non-confidential) clients.
  force_pkce

  # Non-native redirect URIs must use HTTPS outside development.
  force_ssl_in_redirect_uri { |uri| Rails.env.production? && uri.host != "localhost" }
end
