class AddConnectedAppFieldsToOauthApplications < ActiveRecord::Migration[8.1]
  def change
    add_column :oauth_applications, :slug, :string
    add_column :oauth_applications, :description, :string
    add_column :oauth_applications, :homepage_url, :string
    add_column :oauth_applications, :category, :string
    add_column :oauth_applications, :active, :boolean, null: false, default: true

    add_index :oauth_applications, :slug, unique: true
  end
end
