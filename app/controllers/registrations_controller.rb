class RegistrationsController < ApplicationController
  layout "auth"
  allow_unauthenticated_access only: %i[ new create ]

  def new
    @user = User.new
  end

  def create
    @user = User.new(registration_params)

    if @user.save
      @user.roles << Role.find_by(name: "member") if Role.exists?(name: "member")
      start_new_session_for @user
      AuditLog.record(action: "auth.signup", user: @user, request: request)
      redirect_to dashboard_path, notice: "Welcome to Nexus."
    else
      render :new, status: :unprocessable_entity
    end
  end

  private

  def registration_params
    params.require(:user).permit(:display_name, :email_address, :password, :password_confirmation)
  end
end
