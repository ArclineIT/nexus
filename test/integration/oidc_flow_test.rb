require "test_helper"

# Exercises the complete Authorization Code flow: browser login, consent
# (auto-approved for first-party apps), code exchange, ID Token verification
# and the UserInfo endpoint.
class OidcFlowTest < ActionDispatch::IntegrationTest
  setup do
    @app = connected_apps(:portal)
    @user = users(:admin)
    sign_in_as @user
  end

  test "issues a verifiable ID token and userinfo for a signed-in user" do
    code = request_authorization_code

    tokens = exchange_code_for_tokens(code)
    assert tokens["access_token"].present?
    assert tokens["id_token"].present?
    assert tokens["refresh_token"].present?

    claims = decode_id_token(tokens["id_token"])
    assert_equal "http://www.example.com", claims["iss"]
    assert_equal @app.uid, claims["aud"]
    assert_equal @user.id.to_s, claims["sub"]
    assert_equal @user.email_address, claims["email"]
    assert_includes claims["roles"], "admin"

    info = fetch_userinfo(tokens["access_token"])
    assert_equal @user.email_address, info["email"]
    assert_equal @user.display_name, info["name"]
    assert_includes info["roles"], "admin"
  end

  private

  def request_authorization_code
    get "/oauth/authorize", params: {
      client_id: @app.uid,
      redirect_uri: @app.redirect_uri,
      response_type: "code",
      scope: "openid profile email roles",
      state: "state-123"
    }

    assert_response :redirect
    location = URI.parse(response.location)
    assert_equal "portal.arcline.it", location.host

    query = Rack::Utils.parse_query(location.query)
    assert_equal "state-123", query["state"]
    assert query["code"].present?
    query["code"]
  end

  def exchange_code_for_tokens(code)
    post "/oauth/token", params: {
      grant_type: "authorization_code",
      code: code,
      redirect_uri: @app.redirect_uri,
      client_id: @app.uid,
      client_secret: @app.secret
    }

    assert_response :success
    JSON.parse(response.body)
  end

  def decode_id_token(id_token)
    get "/oauth/discovery/keys"
    jwk = JSON.parse(response.body)["keys"].first
    public_key = JWT::JWK.import(jwk).public_key
    claims, = JWT.decode(id_token, public_key, true, algorithm: "RS256")
    claims
  end

  def fetch_userinfo(access_token)
    get "/oauth/userinfo", headers: { "Authorization" => "Bearer #{access_token}" }
    assert_response :success
    JSON.parse(response.body)
  end
end
