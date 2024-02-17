#!/bin/sh
#
set -x
set -eo pipefail

cargo test
cargo tarpaulin --ignore-tests
cargo clippy -- -D warnings
cargo fmt -- --check
cargo audit
# cargo +nightly udeps --all-targets
