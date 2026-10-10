class User < ApplicationRecord
  MIN_PASSWORD_LENGTH = 8

  has_secure_password

  has_many :sessions, dependent: :destroy
  has_many :user_roles, dependent: :destroy
  has_many :roles, through: :user_roles
  has_many :audit_logs, dependent: :nullify

  normalizes :email_address, with: ->(e) { e.strip.downcase }

  validates :email_address, presence: true, uniqueness: true
  validates :display_name, presence: true
  validates :password, length: { minimum: MIN_PASSWORD_LENGTH }, allow_nil: true

  scope :active, -> { where(active: true) }
  scope :ordered, -> { order(:display_name) }

  # Self-expiring token used by the forgot-password flow. Invalidated by a
  # password change because the salt changes.
  generates_token_for :password_reset, expires_in: 15.minutes do
    password_salt&.last(10)
  end

  # Roles and permissions, flattened for the ID token and admin views.

  def admin?
    roles.exists?(name: "admin")
  end

  def has_role?(*names)
    roles.where(name: names).exists?
  end

  def role_names
    roles.pluck(:name)
  end

  def permission_names
    Permission.joins(:role_permissions)
      .where(role_permissions: { role_id: roles.select(:id) })
      .distinct
      .pluck(:name)
  end
end
