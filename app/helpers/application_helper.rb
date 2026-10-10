module ApplicationHelper
  # Two-letter mark for an app card, e.g. "Customer Portal" => "CP".
  def app_initials(name)
    name.to_s.split(/[\s-]+/).first(2).map { |word| word[0] }.join.upcase
  end

  def category_label(category)
    category.to_s.tr("-", " ").split.map(&:capitalize).join(" ")
  end
end
