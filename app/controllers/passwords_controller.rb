class PasswordsController < ApplicationController
  layout "auth"
  allow_unauthenticated_access
  before_action :set_user_by_token, only: %i[ edit update ]
  rate_limit to: 10, within: 3.minutes, only: :create,
             with: -> { redirect_to new_password_path, alert: "Try again later." }

  def new
  end

  def create
    if user = User.find_by(email_address: params[:email_address].to_s.strip.downcase)
      PasswordsMailer.reset(user).deliver_later
    end

    redirect_to new_session_path,
                notice: "Password reset instructions sent (if that account exists)."
  end

  def edit
  end

  def update
    if @user.update(params.permit(:password, :password_confirmation))
      @user.sessions.destroy_all
      AuditLog.record(action: "auth.password_reset", user: @user, request: request)
      redirect_to new_session_path, notice: "Your password has been reset."
    else
      redirect_to edit_password_path(params[:token]),
                  alert: @user.errors.full_messages.to_sentence
    end
  end

  private

  def set_user_by_token
    @user = User.find_by_password_reset_token!(params[:token])
  rescue ActiveSupport::MessageVerifier::InvalidSignature
    redirect_to new_password_path, alert: "That reset link is invalid or has expired."
  end
end
