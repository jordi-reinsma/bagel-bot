mod client;
mod method;

use std::env;

use client::SlackClient;

const CHANNEL_NAMES: &str = include_str!("../channel-names");

fn main() {
    let api_key = env::var("API_KEY").expect("API_KEY not found");
    let client = SlackClient::from_key(&api_key);
}
