use rayon::ThreadPoolBuilder;
use rayon::prelude::*;
use reqwest::header::HeaderValue;
use serde::Deserialize;
use std::error::Error;
use std::io::{self, Write};

static RATE_LIMIT_ERROR_TEXT: &str = "exceeded GitHub API rate limit!";
static USER_AGENT: HeaderValue = HeaderValue::from_static("ghdstats/v1.3.0");

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

pub struct Client {
    pub repos: Vec<Repo>,
}

impl Client {
    pub fn new() -> Self {
        Self { repos: Vec::new() }
    }

    pub fn lookup_repos(&mut self, user: &str) -> Result<&Client, Box<dyn Error>> {
        let resp = reqwest::blocking::Client::new()
            .get(format!("https://api.github.com/users/{user}/repos"))
            .header("User-Agent", &USER_AGENT)
            .send()?;

        if !resp.status().is_success() {
            return Err(Box::new(io::Error::other(RATE_LIMIT_ERROR_TEXT)));
        }
        self.repos = resp.json()?;
        Ok(self)
    }

    pub fn print_all_downloads(&self) -> Result<(), Box<dyn Error>> {
        let pool = ThreadPoolBuilder::new().build()?;
        pool.install(|| {
            self.repos.par_iter().for_each(|repo| {
                Client::print_downloads_for_repo(&repo.full_name).unwrap();
            });
        });
        Ok(())
    }

    pub fn print_downloads_for_repo(full_name: &String) -> Result<(), Box<dyn Error>> {
        let resp = reqwest::blocking::Client::new()
            .get(format!("https://api.github.com/repos/{full_name}/releases"))
            .header("User-Agent", &USER_AGENT)
            .send()?;
        if !resp.status().is_success() {
            return Err(Box::new(io::Error::other(RATE_LIMIT_ERROR_TEXT)));
        }

        let info: Vec<Release> = resp.json()?;

        let mut buffer: Vec<u8> = Vec::with_capacity(1024);
        writeln!(&mut buffer, "Releases for {full_name}:")?;
        if info.is_empty() {
            writeln!(&mut buffer, "- No releases!")?;
        }

        let mut total_downloads: u64 = 0;
        for release in info {
            if release.assets.is_empty() {
                continue;
            }

            let old_count = total_downloads;

            writeln!(&mut buffer, "{}:", release.tag_name)?;
            for asset in release.assets {
                if asset.download_count == 0 {
                    continue;
                }

                total_downloads += asset.download_count;
                writeln!(&mut buffer, "- {}: {}", asset.name, asset.download_count)?;
            }

            if old_count == total_downloads {
                writeln!(&mut buffer, "- No downloads!")?;
            }
        }

        writeln!(&mut buffer, "Total downloads: {total_downloads}")?;
        io::stdout().write_all(&buffer)?;
        Ok(())
    }
}
