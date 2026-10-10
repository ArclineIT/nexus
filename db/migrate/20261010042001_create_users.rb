class CreateUsers < ActiveRecord::Migration[8.1]
  def change
    create_table :users do |t|
      t.string  :email_address, null: false
      t.string  :password_digest, null: false
      t.string  :display_name, null: false
      t.boolean :mfa_enabled, null: false, default: false
      t.string  :mfa_secret
      t.boolean :email_verified, null: false, default: false
      t.boolean :active, null: false, default: true
      t.datetime :last_login_at
      t.timestamps
    end

    add_index :users, :email_address, unique: true
  end
end
