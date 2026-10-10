# A ConnectedApp is an Arcline application registered as an OAuth2 / OpenID
# Connect client with Nexus. It is the "spoke" of the Nexus hub. The model
# reuses Doorkeeper's `oauth_applications` table and mixin so the token and
# grant lifecycle is handled by Doorkeeper.
class ConnectedApp < ApplicationRecord
  include ::Doorkeeper::Orm::ActiveRecord::Mixins::Application

  self.table_name = "oauth_applications"

  CATEGORIES = %w[core infrastructure operations customer internal].freeze

  validates :slug, presence: true, uniqueness: true,
                   format: { with: /\A[a-z0-9]+(?:-[a-z0-9]+)*\z/ }
  validates :homepage_url, format: { with: %r{\Ahttps?://\S+\z}, allow_blank: true }
  validates :category, inclusion: { in: CATEGORIES }, allow_blank: true

  scope :active, -> { where(active: true) }
  scope :ordered, -> { order(:name) }

  def to_param
    slug
  end

  # Absolute URL of the OpenID Connect discovery document for this app, handy
  # when wiring a new client.
  def self.discovery_url
    "https://#{ENV.fetch('NEXUS_HOST', 'nexus.arcline.it')}/.well-known/openid-configuration"
  end
end
