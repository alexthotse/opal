import gleam/dynamic
import wisp
import gleam/dict
import gleam/bytes_tree

import domain/reasoning
import domain/planning
import adapters/jido_agent
import domain/search

@external(erlang, "gleam_stdlib", "identity")
fn to_dynamic(a: a) -> dynamic.Dynamic

@external(erlang, "msgpack", "pack")
fn msgpack_pack(term: dynamic.Dynamic) -> Result(BitArray, dynamic.Dynamic)

pub fn handle_request(req: wisp.Request) -> wisp.Response {
  // Ensure application/msgpack
  use <- wisp.require_content_type(req, "application/msgpack")
  use _body <- wisp.require_bit_array_body(req)

  let path = req.path

  // A naive routing logic for Connect RPC:
  let result = case path {
    "/falcon.v1.FalconService/Ping" -> "pong"
    "/falcon.v1.FalconService/StartUltrathink" -> reasoning.start_ultrathink()
    "/falcon.v1.FalconService/StartUltraplan" -> planning.start_ultraplan()
    "/falcon.v1.FalconService/DispatchAction" -> jido_agent.dispatch_action("test_action")
    "/falcon.v1.FalconService/QuickSearch" -> search.quick_search()
    _ -> "not_found"
  }

  // Pack a simple map
  let response_map = dict.from_list([
    #("Result", result),
    #("Message", result)
  ])

  let response_payload = case msgpack_pack(to_dynamic(response_map)) {
    Ok(bin) -> bin
    Error(_) -> <<>>
  }

  wisp.response(200)
  |> wisp.set_header("content-type", "application/msgpack")
  |> wisp.set_body(wisp.Bytes(bytes_tree.from_bit_array(response_payload)))
}
