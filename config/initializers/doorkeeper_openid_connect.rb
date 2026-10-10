# frozen_string_literal: true

# The RSA key used to sign ID Tokens. Prefer NEXUS_OIDC_PRIVATE_KEY (PEM in an
# env var / secret store). Otherwise read config/keys/oidc.pem, generating a
# throwaway key in development and test when none is present.
oidc_signing_key = begin
  env_key = ENV["NEXUS_OIDC_PRIVATE_KEY"].presence
  if env_key
    env_key
  else
    key_path = Rails.root.join("config/keys/oidc.pem")
    if key_path.exist?
      key_path.read
    elsif Rails.env.local?
      require "openssl"
      require "fileutils"
      key = OpenSSL::PKey::RSA.new(2048)
      FileUtils.mkdir_p(key_path.dirname)
      File.write(key_path, key.to_pem)
      Rails.logger.warn("nexus: generated development OIDC signing key at #{key_path}")
      key.to_pem
    else
      raise "NEXUS_OIDC_PRIVATE_KEY is unset and #{key_path} is missing"
    end
  end
end

Doorkeeper::OpenidConnect.configure do
  # The issuer must be stable and identical in the discovery document, the ID
  # Token `iss` claim and the RFC 9207 authorization response. Configure it
  # explicitly per environment rather than deriving it from each request.
  issuer ENV.fetch("NEXUS_ISSUER") {
    Rails.env.local? ? "http://localhost:3000" : "https://nexus.arcline.it"
  }

  signing_key oidc_signing_key

  subject_types_supported [ :public ]

  resource_owner_from_access_token do |access_token|
    User.find_by(id: access_token.resource_owner_id)
  end

  auth_time_from_resource_owner do |resource_owner|
    resource_owner.last_login_at || resource_owner.created_at
  end

  subject do |resource_owner, _application|
    resource_owner.id.to_s
  end

  # Claims are emitted either in the ID Token, the UserInfo response, or both,
  # gated by the scope the client requested.
  claims do
    normal_claim :email, scope: :email, response: %i[ id_token user_info ] do |resource_owner|
      resource_owner.email_address
    end

    normal_claim :email_verified, scope: :email, response: %i[ id_token user_info ] do |resource_owner|
      resource_owner.email_verified
    end

    normal_claim :name, scope: :profile, response: %i[ id_token user_info ] do |resource_owner|
      resource_owner.display_name
    end

    normal_claim :preferred_username, scope: :profile, response: %i[ user_info ] do |resource_owner|
      resource_owner.email_address.split("@").first
    end

    normal_claim :roles, scope: :roles, response: %i[ id_token user_info ] do |resource_owner|
      resource_owner.role_names
    end

    normal_claim :permissions, scope: :roles, response: %i[ user_info ] do |resource_owner|
      resource_owner.permission_names
    end
  end
end
