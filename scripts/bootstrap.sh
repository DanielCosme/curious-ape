#!/bin/sh

cargo install cargo-watch
cargo install cargo-tarpaulin
rustup component add clippy
rustup component add rustfmt
cargo install cargo-audit

cargo install --version="~0.7" sqlx-cli --no-default-features \
    --features rustls,postgres
