pub struct EmailClient {
    sender: String
}

impl EmailClient {
    pub async fn send_email(
        &self,
        recipient: String,
        subject: &str,
        html_content: &str,
        text_content: &str
    ) -> Result<(), String> {
        todo!()
    }
}
