require "test_helper"

class OidcTest < ActionDispatch::IntegrationTest
  test "discovery document advertises issuer, endpoints and scopes" do
    get "/.well-known/openid-configuration"
    assert_response :success

    doc = JSON.parse(response.body)
    assert_equal "http://www.example.com", doc["issuer"]
    assert_equal "http://www.example.com/oauth/authorize", doc["authorization_endpoint"]
    assert_equal "http://www.example.com/oauth/token", doc["token_endpoint"]
    assert_equal "http://www.example.com/oauth/userinfo", doc["userinfo_endpoint"]
    assert_includes doc["scopes_supported"], "roles"
    assert_includes doc["response_types_supported"], "code"
  end

  test "JWKS exposes the RSA signing key" do
    get "/oauth/discovery/keys"
    assert_response :success

    keys = JSON.parse(response.body)["keys"]
    assert_equal 1, keys.length
    assert_equal "RSA", keys.first["kty"]
  end

  test "authorization endpoint sends anonymous users to login" do
    app = connected_apps(:portal)

    get "/oauth/authorize", params: {
      client_id: app.uid,
      redirect_uri: app.redirect_uri,
      response_type: "code",
      scope: "openid profile email"
    }

    assert_response :redirect
    assert_match %r{/login}, response.location
  end
end
