@0x9c0139b8202143ba;

using Go = import "/go.capnp";

$Go.package("main");
$Go.import("github.com/iwasaki-kenta/capnproto-benchmark");

struct Message {
    bytes @0: Data;
}

interface Benchmark {
    send @0 (req: Message) -> (res: Message);
}