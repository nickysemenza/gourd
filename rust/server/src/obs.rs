pub fn initialize_tracing(_s: &str) {
    init_tracing_opentelemetry::tracing_subscriber_ext::init_subscribers().expect("failed to init")
}
