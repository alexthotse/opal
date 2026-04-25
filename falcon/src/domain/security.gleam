import gleam/string

pub type SecurityRule {
  SecurityRule(pattern: String, action: String)
}

pub fn bash_classifier(command: String) -> String {
  let is_unsafe = string.contains(command, "rm -rf") || string.contains(command, "sudo")
  
  case is_unsafe {
    True -> "Bash Classifier Safe check for '" <> command <> "':\nResult: BLOCK"
    False -> "Bash Classifier Safe check for '" <> command <> "':\nResult: ALLOW"
  }
}
