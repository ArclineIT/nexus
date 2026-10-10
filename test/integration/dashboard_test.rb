require "test_helper"

class DashboardTest < ActionDispatch::IntegrationTest
  test "redirects unauthenticated visitors to login" do
    get dashboard_url
    assert_redirected_to new_session_url
  end

  test "lists active connected apps for signed-in users" do
    sign_in_as users(:admin)

    get dashboard_url
    assert_response :success
    assert_match "Customer Portal", response.body
    assert_no_match "Website", response.body
  end

  test "non-admins cannot reach the admin area" do
    sign_in_as users(:member)

    get admin_root_url
    assert_redirected_to root_url
  end

  test "admins can manage connected apps" do
    sign_in_as users(:admin)

    get admin_root_url
    assert_response :success
    assert_match "Customer Portal", response.body
  end

  test "admin can render the application detail, edit and users pages" do
    sign_in_as users(:admin)
    app = connected_apps(:portal)

    get admin_connected_app_url(app)
    assert_response :success
    assert_match app.uid, response.body

    get edit_admin_connected_app_url(app)
    assert_response :success

    get admin_users_url
    assert_response :success
    assert_match users(:admin).email_address, response.body
  end
end
