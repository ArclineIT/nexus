module Admin
  class ConnectedAppsController < BaseController
    before_action :set_app, only: %i[ show edit update regenerate_secret ]

    def index
      @apps = ConnectedApp.ordered
    end

    def show
    end

    def edit
    end

    def update
      if @app.update(app_params)
        redirect_to admin_connected_app_path(@app), notice: "#{@app.name} updated."
      else
        render :edit, status: :unprocessable_entity
      end
    end

    def regenerate_secret
      @app.update!(secret: SecureRandom.hex(32))
      AuditLog.record(action: "admin.app.secret_rotated", user: Current.user,
                      resource: @app.slug, request: request)
      redirect_to admin_connected_app_path(@app),
                  notice: "New client secret generated. Copy it now — it will not be shown again."
    end

    private

    def set_app
      @app = ConnectedApp.find_by!(slug: params[:id])
    end

    def app_params
      params.require(:connected_app).permit(:name, :description, :homepage_url, :category, :active, :redirect_uri, :scopes)
    end
  end
end
