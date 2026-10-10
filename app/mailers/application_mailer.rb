class ApplicationMailer < ActionMailer::Base
  default from: ENV.fetch("NEXUS_MAIL_FROM", "Nexus <noreply@arcline.it>")
  layout "mailer"
end
