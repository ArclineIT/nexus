class CreateAuditLogs < ActiveRecord::Migration[8.1]
  def change
    create_table :audit_logs do |t|
      t.references :user, null: true, foreign_key: true
      t.string :action, null: false
      t.string :resource
      t.string :ip_address
      t.string :user_agent
      t.text :metadata
      t.datetime :created_at, null: false
    end

    add_index :audit_logs, :action
    add_index :audit_logs, :created_at
  end
end
