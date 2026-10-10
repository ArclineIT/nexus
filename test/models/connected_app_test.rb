require "test_helper"

class ConnectedAppTest < ActiveSupport::TestCase
  test "rejects an invalid slug" do
    app = ConnectedApp.new(name: "Test", slug: "Bad Slug",
                           redirect_uri: "https://x.arcline.it/cb", scopes: "openid")
    assert_not app.valid?
    assert app.errors[:slug].any?
  end

  test "active scope excludes inactive apps" do
    assert_includes ConnectedApp.active, connected_apps(:portal)
    assert_not_includes ConnectedApp.active, connected_apps(:inactive_app)
  end

  test "builds the discovery URL from the configured host" do
    assert_equal "https://nexus.arcline.it/.well-known/openid-configuration",
                 ConnectedApp.discovery_url
  end
end
