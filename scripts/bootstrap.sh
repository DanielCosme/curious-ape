#!/bin/sh

rustup toolchain install nightly
rustup component add clippy
rustup component add rustfmt
cargo install cargo-watch
cargo install cargo-tarpaulin
cargo install cargo-audit
cargo install bunyan
cargo install --version="~0.7" sqlx-cli --no-default-features \
    --features rustls,postgres

cargo install cargo-udeps
