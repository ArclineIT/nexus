module Admin
  class UsersController < BaseController
    def index
      @users = User.includes(:roles).ordered
      @roles = Role.ordered
    end

    def update
      @user = User.find(params[:id])
      @user.role_ids = Array(params[:role_ids]).reject(&:blank?)
      AuditLog.record(action: "admin.user.roles_updated", user: Current.user,
                      resource: @user.email_address, request: request)
      redirect_to admin_users_path, notice: "#{@user.display_name} updated."
    end
  end
end
