import gleam/string
import gleam/list

pub type ReasoningStep {
  ReasoningStep(id: Int, action: String, status: String)
}

pub type ReasoningContext {
  ReasoningContext(query: String, steps: List(ReasoningStep))
}

pub fn start_ultrathink() -> String {
  let ctx = ReasoningContext(
    query: "Optimize search algorithm",
    steps: [
      ReasoningStep(1, "Analyze current search complexity", "done"),
      ReasoningStep(2, "Identify bottleneck in indexing", "done"),
      ReasoningStep(3, "Propose concurrent tree traversal", "pending"),
    ]
  )
  
  let steps_str = 
    ctx.steps
    |> list.map(fn(s) { "Step " <> string.inspect(s.id) <> " [" <> s.status <> "]: " <> s.action })
    |> string.join("\n")

  "ULTRATHINK initialized for query: '" <> ctx.query <> "'\n" <> steps_str
}
