use reqwest::Client;

pub struct FitbitClient {
    http_client: Client,
    token: String,
}

impl FitbitClient {
    pub fn new() -> Self {
        Self {
            http_client: Client::new(),
            token: String::new(),
        }
    }
}

// Fitibt Client
//    - Make requests with custom http cllient.
//    - Receives Token.
//    - Should be able to exchange the token if it expired?
//
// Auth for all providers.
//    - Read environment variables.
//    - Authenticate & Authorize.
//    - Reads token and refreshes it?
//
// Could I build stand alone clients?
