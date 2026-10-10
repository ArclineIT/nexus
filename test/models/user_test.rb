require "test_helper"

class UserTest < ActiveSupport::TestCase
  test "normalizes email address" do
    user = User.new(email_address: "  ADMIN@Arcline.IT ", display_name: "A", password: "password123")
    assert_equal "admin@arcline.it", user.email_address
  end

  test "enforces a minimum password length" do
    user = User.new(email_address: "x@arcline.it", display_name: "X", password: "short")
    assert_not user.valid?
    assert user.errors[:password].any?
  end

  test "admin? reflects the admin role" do
    assert users(:admin).admin?
    assert_not users(:member).admin?
  end

  test "role_names exposes assigned roles" do
    assert_includes users(:admin).role_names, "admin"
  end
end
