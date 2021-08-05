mod client;
mod error;
mod method;
mod partition;

use std::env;

use error::Result;
use client::SlackClient;

const CHANNEL_IDS: &str = include_str!("../channel-ids");

#[tokio::main]
async fn main() {
    let channels: Vec<_> = CHANNEL_IDS.split('\n').collect();

    for channel in channels {
        let _ = dbg!(set_up_meetings(channel).await);
    }
}

async fn set_up_meetings(channel_id: &str) -> Result<()> {
    // Get your token at `https://api.slack.com/apps/<your-bot-id>/oauth?`
    let oauth_token = env::var("SLACK_OAUTH_TOKEN").expect("SLACK_OAUTH_TOKEN not found");
    let client = SlackClient::from_key(&oauth_token);

    let users = client.members_of_channel(channel_id).await?;
    dbg!(&users);

    let channel_id = client.start_direct_message(users).await?;
    dbg!(&channel_id);

    let success = client.post_message(&channel_id, "Olá, abigos!").await?;
    
    Ok(())
}