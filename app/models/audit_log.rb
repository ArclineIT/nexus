class AuditLog < ApplicationRecord
  belongs_to :user, optional: true

  validates :action, presence: true

  scope :recent, -> { order(created_at: :desc) }

  # Records a security-relevant event. Accepts an optional ActionDispatch
  # request so callers don't have to pluck IP/user-agent themselves.
  def self.record(action:, user: nil, resource: nil, request: nil, metadata: nil)
    create!(
      action: action,
      user: user,
      resource: resource,
      ip_address: request&.remote_ip,
      user_agent: request&.user_agent,
      metadata: metadata&.to_json
    )
  end
end
