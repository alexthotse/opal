import gleam/string
import gleam/list

pub type Node {
  Node(id: String, role: String)
}

pub type PlanGraph {
  PlanGraph(nodes: List(Node), dependencies: List(#(String, String)))
}

pub fn start_ultraplan() -> String {
  let plan = PlanGraph(
    nodes: [
      Node("db_agent", "Database Migrations"),
      Node("api_agent", "API Endpoints"),
      Node("ui_agent", "Frontend Updates")
    ],
    dependencies: [
      #("db_agent", "api_agent"),
      #("api_agent", "ui_agent")
    ]
  )

  let nodes_str = plan.nodes |> list.map(fn(n) { n.id <> " (" <> n.role <> ")" }) |> string.join(", ")
  let deps_str = plan.dependencies |> list.map(fn(d) { d.0 <> " -> " <> d.1 }) |> string.join(", ")

  "ULTRAPLAN initialized: Building architectural graph...\nNodes: " <> nodes_str <> "\nDependencies: " <> deps_str
}
