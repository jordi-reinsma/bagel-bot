mod client;
mod method;
mod partition;

use std::env;
use std::collections::HashMap;

use client::SlackClient;

const CHANNEL_NAMES: &str = include_str!("../channel-names");

fn main() {
    let api_key = env::var("API_KEY").expect("API_KEY not found");
    let client = SlackClient::from_key(&api_key);

    let mut parameters = HashMap::new();
    parameters.insert(
        "channel",
        "C025U17MQ0J",
    );
    parameters.insert(
        "api-key",
        api_key,
    );

    let result = client.send(method::Method::ListMembersOfChannel, parameters);
    dbg!(result); // doest not work because result is a Future
}
