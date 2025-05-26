use reqwest::header::{HeaderMap, HeaderValue};
use serde::Deserialize;
use std::env;

#[derive(Deserialize)]
struct Repo {
    full_name: String,
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

fn main() -> Result<(), reqwest::Error> {
    let args: Vec<String> = env::args().collect();

    let mut repos: Vec<Repo> = Vec::new();
    match args.len() {
        0 | 1 | 4.. => {
            println!("Usage: gcdstats [user] [repository, optional]");
            return Ok(());
        }
        3 => repos.push(Repo {
            full_name: format!("{}/{}", &args[1], &args[2]),
        }),
        2 => {
            let parts: Vec<&str> = args[1].split('/').collect();
            match parts.len() {
                1 => repos = fetch_repos(parts[0])?,
                2 => repos.push(Repo {
                    full_name: String::from(&args[1]),
                }),
                _ => {
                    println!("Invalid input format");
                    return Ok(());
                }
            }
        }
    }

    fetch_statistics(&repos)?;
    Ok(())
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
        print_repo_info(&repo.full_name)?;
    }
    Ok(())
}

fn print_repo_info(full_name: &str) -> Result<(), reqwest::Error> {
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
