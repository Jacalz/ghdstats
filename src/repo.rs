use serde::Deserialize;

#[derive(Deserialize)]
pub struct Repository {
    pub full_name: String,
}

pub fn fetch_repositories(user: String) -> Vec<Repository> {
    let mut users_api_endpoint = String::from("https://api.github.com/users/");
    users_api_endpoint.push_str(&user);
    users_api_endpoint.push_str("/repos");

    return ureq::get(users_api_endpoint.as_str()).call()?.to_json()?;
}