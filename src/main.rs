use reqwest::header::{HeaderMap, HeaderValue};
use serde::Deserialize;
use std::env;

#[derive(Deserialize)]
struct Repo {
    full_name: String,
}

#[derive(Deserialize)]
struct Release {
    assets: Vec<Asset>,
}

#[derive(Deserialize)]
struct Asset {
    name: String,
    download_count: u64,
    updated_at: String,
}

fn main() {
    let args: Vec<String> = env::args().collect();

    let mut repos: Vec<Repo> = Vec::new();
    match args.len() {
        2 => repos = fetch_repos(&args[1]).unwrap(),
        3 => repos.push(Repo {
            full_name: format!("{}/{}", &args[1], &args[2]),
        }),
        _ => {
            println!("Usage: gcdstats [user] [repository, optional]");
            return;
        }
    }

    fetch_statistics(&repos).unwrap();
}

fn fetch_repos(user: &str) -> Result<Vec<Repo>, reqwest::Error> {
    let client = reqwest::blocking::Client::new();
    let mut headers = HeaderMap::new();
    headers.insert("User-Agent", HeaderValue::from_static("ghdstats/v1.3.0"));

    client
        .get(format!("https://api.github.com/users/{user}/repos"))
        .headers(headers)
        .send()?
        .json()
}

fn fetch_statistics(repos: &Vec<Repo>) -> Result<(), reqwest::Error> {
    for repo in repos {
        print_repo_info(repo)?;
    }
    Ok(())
}

fn print_repo_info(repo: &Repo) -> Result<(), reqwest::Error> {
    let client = reqwest::blocking::Client::new();
    let mut headers = HeaderMap::new();
    headers.insert("User-Agent", HeaderValue::from_static("ghdstats/v1.3.0"));

    let info: Vec<Release> = client
        .get(format!(
            "https://api.github.com/repos/{}/releases",
            repo.full_name
        ))
        .headers(headers)
        .send()?
        .json()?;

    let mut total_downloads: u64 = 0;

    for release in info {
        for asset in release.assets {
            total_downloads += asset.download_count;
        }
    }

    println!(
        "Total downloads for {}: {}",
        repo.full_name, total_downloads
    );
    Ok(())
}
