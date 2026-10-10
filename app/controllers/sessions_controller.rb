class SessionsController < ApplicationController
  layout "auth"
  allow_unauthenticated_access only: %i[ new create ]
  rate_limit to: 10, within: 3.minutes, only: :create,
             with: -> { redirect_to new_session_path, alert: "Try again later." }

  def new
  end

  def create
    user = User.authenticate_by(params.permit(:email_address, :password))

    if user&.active?
      start_new_session_for user
      user.update_column(:last_login_at, Time.current)
      AuditLog.record(action: "auth.login", user: user, request: request)
      redirect_to after_authentication_url, notice: "Signed in successfully."
    else
      AuditLog.record(action: "auth.login_failed", resource: params[:email_address], request: request)
      redirect_to new_session_path, alert: "Invalid email or password."
    end
  end

  def destroy
    AuditLog.record(action: "auth.logout", user: Current.user, request: request) if Current.user
    terminate_session
    redirect_to new_session_path, notice: "Signed out.", status: :see_other
  end
end
