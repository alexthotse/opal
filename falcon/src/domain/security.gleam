import gleam/string
import gleam/list

pub type SecurityRule {
  SecurityRule(pattern: String, action: String)
}

pub fn bash_classifier() -> String {
  let rules = [
    SecurityRule("rm -rf /", "BLOCK"),
    SecurityRule("ls -la", "ALLOW")
  ]
  let check_command = "ls -la"
  let rule_str = rules
    |> list.map(fn(r) { "Rule '" <> r.pattern <> "' -> " <> r.action })
    |> string.join("\n")
  
  "Bash Classifier Safe check for '" <> check_command <> "':\n" <> rule_str <> "\nResult: SAFE"
}
