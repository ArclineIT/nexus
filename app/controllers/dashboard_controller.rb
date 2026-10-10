class DashboardController < ApplicationController
  def show
    @apps = ConnectedApp.active.ordered
    @recent_events = Current.user.audit_logs.recent.limit(5)
  end
end
