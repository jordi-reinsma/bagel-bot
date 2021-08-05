use std::{borrow::Borrow, collections::HashMap};

use reqwest::Url;
use serde_json::Value;

use crate::error::Result;
use crate::method::Method;

pub struct SlackClient<'a> {
    api_key: &'a str,
    http_client: reqwest::Client,
}

impl<'a> SlackClient<'a> {
    pub fn from_key(api_key: &'a str) -> Self {
        Self {
            api_key,
            http_client: reqwest::Client::new(),
        }
    }

    // todo: error treatment
    pub async fn send<P, K, V>(&self, method: Method, parameters: P) -> Result<serde_json::Value>
    where
        P: IntoIterator + Send,
        K: AsRef<str>,
        V: AsRef<str>,
        P::Item: Borrow<(K, V)>,
    {
        let mut url: Url = method.into();

        // Adds a sequence of name/value pairs in `application/x-www-form-urlencoded` syntax
        // to the URL
        url.query_pairs_mut().extend_pairs(parameters);

        let response = self
            .http_client
            .post(url)
            .bearer_auth(self.api_key)
            .send()
            .await?
            .text()
            .await?;

        Ok(serde_json::from_str(&response)?)
    }
}
