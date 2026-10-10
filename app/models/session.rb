class Session < ApplicationRecord
  # Browser login sessions expire after this duration (matches the cookie lifetime).
  TTL = 30.days

  belongs_to :user

  scope :active,  -> { where(created_at: TTL.ago..) }
  scope :expired, -> { where(created_at: ...TTL.ago) }

  def expired?
    created_at < TTL.ago
  end
end
