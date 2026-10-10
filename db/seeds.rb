# Seeds are idempotent: safe to run in any environment at any time.
#   bin/rails db:seed

# ---------------------------------------------------------------------------
# Roles & permissions
# ---------------------------------------------------------------------------

permissions = [
  { name: "apps.read",    resource: "apps",  action: "read",    description: "View connected applications" },
  { name: "apps.manage",  resource: "apps",  action: "manage",  description: "Register and edit connected applications" },
  { name: "users.manage", resource: "users", action: "manage",  description: "Manage users and role assignments" },
  { name: "audit.read",   resource: "audit", action: "read",    description: "Read the audit log" }
].map do |attrs|
  Permission.find_or_create_by!(resource: attrs[:resource], action: attrs[:action]) do |p|
    p.name = attrs[:name]
    p.description = attrs[:description]
  end
end
permission_by_name = permissions.index_by(&:name)

roles = {
  "admin"    => %w[apps.read apps.manage users.manage audit.read],
  "operator" => %w[apps.read audit.read],
  "member"   => %w[apps.read],
  "viewer"   => %w[apps.read]
}.map do |name, permission_names|
  role = Role.find_or_create_by!(name: name)
  role.permissions = permission_names.map { |n| permission_by_name.fetch(n) }
  role
end
role_by_name = roles.index_by(&:name)

# ---------------------------------------------------------------------------
# Bootstrap admin user
# ---------------------------------------------------------------------------

admin_email = ENV.fetch("NEXUS_ADMIN_EMAIL", "admin@arcline.it")
admin_password = ENV["NEXUS_ADMIN_PASSWORD"]
admin_password ||= "password123" if Rails.env.local?

admin = User.find_by(email_address: admin_email)

if admin.nil? && admin_password.present?
  admin = User.create!(
    email_address: admin_email,
    display_name: ENV.fetch("NEXUS_ADMIN_NAME", "Nexus Admin"),
    password: admin_password,
    password_confirmation: admin_password,
    email_verified: true,
    active: true
  )
  puts "Created admin user #{admin_email}"
elsif admin.nil?
  warn "NEXUS_ADMIN_PASSWORD not set; skipping admin user creation."
end

if admin && !admin.roles.exists?(name: "admin")
  admin.roles << role_by_name.fetch("admin")
  puts "Granted admin role to #{admin_email}"
end

# ---------------------------------------------------------------------------
# Connected applications (the "spokes" of the Nexus hub)
# ---------------------------------------------------------------------------

apps = [
  { slug: "portal",       name: "Customer Portal",  category: "customer",
    description: "SSL expiry monitoring and support tickets for Arcline customers.",
    homepage_url: "https://portal.arcline.it",       redirect_uri: "https://portal.arcline.it/auth/nexus/callback" },
  { slug: "billing",      name: "Billing",          category: "customer",
    description: "Subscriptions, invoices and payments.",
    homepage_url: "https://billing.arcline.it",      redirect_uri: "https://billing.arcline.it/auth/nexus/callback" },
  { slug: "git",          name: "Git",              category: "core",
    description: "Arcline source control (Forgejo).",
    homepage_url: "https://git.arcline.it",          redirect_uri: "https://git.arcline.it/user/oauth2/nexus/callback" },
  { slug: "uptime",       name: "Monitoring",       category: "operations",
    description: "Uptime and service health monitoring.",
    homepage_url: "https://uptime.arcline.it",       redirect_uri: "https://uptime.arcline.it/auth/nexus/callback" },
  { slug: "status",       name: "Status Page",      category: "operations",
    description: "Public and internal service status.",
    homepage_url: "https://status.arcline.it",       redirect_uri: "https://status.arcline.it/auth/nexus/callback" },
  { slug: "email",        name: "Email",            category: "infrastructure",
    description: "Arcline transactional email service.",
    homepage_url: "https://email.arcline.it",        redirect_uri: "https://email.arcline.it/auth/nexus/callback" },
  { slug: "docs",         name: "Docs",             category: "internal",
    description: "Internal documentation and runbooks.",
    homepage_url: "https://docs.arcline.it",         redirect_uri: "https://docs.arcline.it/auth/nexus/callback" },
  { slug: "migrate",      name: "Migrate",          category: "operations",
    description: "Server migration tooling.",
    homepage_url: "https://migrate.arcline.it",      redirect_uri: "https://migrate.arcline.it/auth/nexus/callback" },
  { slug: "provisioner",  name: "Provisioner",      category: "infrastructure",
    description: "Automated infrastructure provisioning.",
    homepage_url: "https://provisioner.arcline.it",  redirect_uri: "https://provisioner.arcline.it/auth/nexus/callback" },
  { slug: "vault",        name: "Vault",            category: "infrastructure",
    description: "Secrets management.",
    homepage_url: "https://vault.arcline.it",        redirect_uri: "https://vault.arcline.it/auth/nexus/callback" },
  { slug: "audit",        name: "Audit",            category: "internal",
    description: "Security and compliance audit trail.",
    homepage_url: "https://audit.arcline.it",        redirect_uri: "https://audit.arcline.it/auth/nexus/callback" },
  { slug: "check",        name: "Check",            category: "operations",
    description: "Configuration and health checks.",
    homepage_url: "https://check.arcline.it",        redirect_uri: "https://check.arcline.it/auth/nexus/callback" },
  { slug: "dns",          name: "DNS",              category: "infrastructure",
    description: "DNS zone management.",
    homepage_url: "https://dns.arcline.it",          redirect_uri: "https://dns.arcline.it/auth/nexus/callback" },
  { slug: "network",      name: "Network Planning", category: "infrastructure",
    description: "Network topology, configurations and planning.",
    homepage_url: "https://network.arcline.it",     redirect_uri: "https://network.arcline.it/auth/nexus/callback" },
  { slug: "website",      name: "Website",          category: "customer", active: false,
    description: "Public arcline.it marketing site.",
    homepage_url: "https://arcline.it",              redirect_uri: "https://arcline.it/auth/nexus/callback" }
]

apps.each do |attrs|
  app = ConnectedApp.find_or_initialize_by(slug: attrs[:slug])
  app.assign_attributes(
    name: attrs[:name],
    description: attrs[:description],
    homepage_url: attrs[:homepage_url],
    category: attrs[:category],
    active: attrs.fetch(:active, true),
    redirect_uri: attrs[:redirect_uri],
    scopes: "openid profile email roles offline_access"
  )
  app.save!
end

puts "Seeded #{ConnectedApp.count} connected apps, #{Role.count} roles, #{Permission.count} permissions."
