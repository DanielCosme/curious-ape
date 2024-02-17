use tracing::{subscriber::set_global_default, Subscriber};
use tracing_bunyan_formatter::{BunyanFormattingLayer, JsonStorageLayer};
use tracing_log::LogTracer;
use tracing_subscriber::fmt::MakeWriter;
use tracing_subscriber::{layer::SubscriberExt, EnvFilter, Registry};

pub fn get_tracing_subscriber<Writer>(
    name: String,
    env_filter: String,
    writer: Writer,
) -> impl Subscriber + Send + Sync
where
    Writer: for<'a> MakeWriter<'a> + Send + Sync + 'static,
{
    let env_filter =
        EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new(env_filter));
    let fomatting_layer = BunyanFormattingLayer::new(name, writer);
    Registry::default()
        .with(env_filter)
        .with(JsonStorageLayer)
        .with(fomatting_layer)
}

pub fn init_subscriber(subscriber: impl Subscriber + Send + Sync) {
    // Redirect all log's events to the tracing subscriber.
    LogTracer::init().expect("Failed to set logger");
    set_global_default(subscriber).expect("Failed to set tracing subscriber");
}
