import gleam/io
import gleam/dynamic

@external(erlang, "Elixir.Jido", "cmd")
fn jido_cmd(agent: dynamic.Dynamic, action: dynamic.Dynamic) -> dynamic.Dynamic

pub fn main() {
  io.println("Hello Jido")
}
