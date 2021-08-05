mod client;
mod error;
mod method;

use std::{collections::HashMap, env};

use client::SlackClient;
use method::Method;

const CHANNEL_NAMES: &str = include_str!("../channel-names");

#[tokio::main]
async fn main() {
    // Get it at `https://api.slack.com/apps/<your-bot-id>/oauth?`
    let oauth_token = env::var("SLACK_OAUTH_TOKEN").expect("SLACK_OAUTH_TOKEN not found");
    let client = SlackClient::from_key(&oauth_token);

    // let mut parameters = HashMap::new();
    // parameters.insert("channel", "");

    let response = client.members_of_channel("C02A4KVB43F").await.unwrap();

    dbg!(response);
}
