use reqwest::header::{HeaderMap, HeaderValue};
use serde::Deserialize;

#[derive(Deserialize)]
pub struct Repo {
    pub full_name: String,
}

#[derive(Deserialize)]
struct Release {
    tag_name: String,
    assets: Vec<Asset>,
}

#[derive(Deserialize)]
struct Asset {
    name: String,
    download_count: u64,
}

pub fn fetch_repos(user: &str) -> Result<Vec<Repo>, reqwest::Error> {
    let client = reqwest::blocking::Client::new();
    let mut headers = HeaderMap::new();
    headers.insert("User-Agent", HeaderValue::from_static("ghdstats/v1.3.0"));

    client
        .get(format!("https://api.github.com/users/{user}/repos"))
        .headers(headers)
        .send()?
        .json()
}

pub fn fetch_statistics(repos: &Vec<Repo>) -> Result<(), reqwest::Error> {
    for repo in repos {
        print_repo_info(&repo.full_name)?;
    }
    Ok(())
}

pub fn print_repo_info(full_name: &str) -> Result<(), reqwest::Error> {
    let client = reqwest::blocking::Client::new();
    let mut headers = HeaderMap::new();
    headers.insert("User-Agent", HeaderValue::from_static("ghdstats/v1.3.0"));

    let info: Vec<Release> = client
        .get(format!("https://api.github.com/repos/{full_name}/releases",))
        .headers(headers)
        .send()?
        .json()?;

    println!("Releases for {full_name}:");
    if info.is_empty() {
        println!("- No releases!");
    }

    let mut total_downloads: u64 = 0;
    for release in info {
        if release.assets.is_empty() {
            continue;
        }

        let old_count = total_downloads;

        println!("{}:", release.tag_name);
        for asset in release.assets {
            if asset.download_count == 0 {
                continue;
            }

            total_downloads += asset.download_count;
            println!("- {}: {}", asset.name, asset.download_count);
        }

        if old_count == total_downloads {
            println!("- No downloads!");
        }
    }

    println!("Total downloads: {total_downloads}");
    Ok(())
}
