ENV["RAILS_ENV"] ||= "test"
require_relative "../config/environment"
require "rails/test_help"

module ActiveSupport
  class TestCase
    # Run tests in parallel with specified workers
    parallelize(workers: :number_of_processors)

    # Setup all fixtures in test/fixtures/*.yml for all tests in alphabetical order.
    fixtures :all

    # ConnectedApp reuses Doorkeeper's table, so point the fixture at the model.
    set_fixture_class connected_apps: ConnectedApp

    # Add more helper methods to be used by all tests here...
    def sign_in_as(user, password: "password123")
      post session_url, params: { email_address: user.email_address, password: password }
    end
  end
end
