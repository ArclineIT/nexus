class Permission < ApplicationRecord
  has_many :role_permissions, dependent: :destroy
  has_many :roles, through: :role_permissions

  validates :name, presence: true, uniqueness: true
  validates :resource, :action, presence: true
  validates :action, uniqueness: { scope: :resource }

  scope :ordered, -> { order(:resource, :action) }
end
