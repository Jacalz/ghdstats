use reqwest::header::{HeaderMap, HeaderValue};
use serde::Deserialize;
use std::error;
use std::io::{self, Write};

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

    client: reqwest::blocking::Client,
}

impl Client {
    pub fn new() -> Self {
        let mut headers = HeaderMap::new();
        headers.insert("User-Agent", HeaderValue::from_static("ghdstats/v1.3.0"));

        Self {
            repos: Vec::new(),
            client: reqwest::blocking::Client::builder()
                .default_headers(headers)
                .build()
                .unwrap(),
        }
    }

    pub fn add_repo(&mut self, full_name: String) {
        self.repos.push(Repo { full_name });
    }

    pub fn lookup_repos(&mut self, user: &str) -> Result<(), reqwest::Error> {
        match self
            .client
            .get(format!("https://api.github.com/users/{user}/repos"))
            .send()?
            .json()
        {
            Ok(repos) => {
                self.repos = repos;
                Ok(())
            }
            Err(err) => Err(err),
        }
    }

    pub fn print_downloads(&self) -> Result<(), Box<dyn error::Error>> {
        for repo in &self.repos {
            self.print_downloads_for_repo(&repo.full_name)?;
        }
        Ok(())
    }

    fn print_downloads_for_repo(&self, full_name: &str) -> Result<(), Box<dyn error::Error>> {
        let info: Vec<Release> = self
            .client
            .get(format!("https://api.github.com/repos/{full_name}/releases",))
            .send()?
            .json()?;

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
