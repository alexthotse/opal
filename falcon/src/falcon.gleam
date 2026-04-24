import gleam/erlang/process
import mist
import ports/connect_server
import wisp

pub fn main() {
  wisp.configure_logger()

  let assert Ok(_) =
    wisp.mist_handler(connect_server.handle_request, "secret")
    |> mist.new
    |> mist.port(8080)
    |> mist.start_http

  process.sleep_forever()
}
