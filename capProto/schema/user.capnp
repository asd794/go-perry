@0xbf5147d3e5f3e4a1;

using Go = import "/go.capnp";

$Go.package("generated");
$Go.import("go-perry/capProto/generated");

struct User {
    id @0 :UInt64;
    name @1 :Text;
    age @2 :UInt8;
    active @3 :Bool;
}