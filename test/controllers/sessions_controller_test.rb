require "test_helper"

class SessionsControllerTest < ActionDispatch::IntegrationTest
  test "signs in with valid credentials" do
    assert_difference -> { AuditLog.where(action: "auth.login").count }, 1 do
      post session_url, params: { email_address: users(:admin).email_address, password: "password123" }
    end

    assert_redirected_to root_url
    assert users(:admin).reload.last_login_at.present?
  end

  test "rejects invalid credentials" do
    post session_url, params: { email_address: users(:admin).email_address, password: "wrong-password" }

    assert_redirected_to new_session_url
    assert_equal "Invalid email or password.", flash[:alert]
  end

  test "signs out" do
    sign_in_as users(:admin)
    delete logout_url
    assert_redirected_to new_session_url
  end
end
