FROM rust:1.86 as builder

WORKDIR /usr/src/search
COPY Cargo.* .
RUN mkdir -p src
RUN echo 'fn main() {println!("for cache");}' > src/main.rs
RUN cargo build --release
COPY ./src ./src
# The last modified attribute of main.rs needs to be updated manually,
# otherwise cargo won't rebuild it.
RUN touch -a -m ./src/main.rs
RUN cargo build --release

FROM debian:12-slim

RUN apt-get update && apt-get install -y libssl-dev ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/search/target/release/search /usr/local/bin/search
COPY --from=builder /usr/src/search/.env /.env

EXPOSE 3001

CMD ["search"] 
