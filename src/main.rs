mod api;
use anyhow::Result;
use clap::{Parser, command};

#[derive(Parser)]
#[command(author, version, about, long_about = None)]
/// A CLI tool to fetch download statistics for GitHub repositories.
struct Args {
    /// The GitHub username to fetch download statistics for.
    user: String,

    /// Optionally, the specific repository to fetch download statistics for.
    repo: Option<String>,
}

fn main() -> Result<()> {
    let args = Args::parse();

    if let Some(repo) = args.repo {
        return api::Client::print_downloads_for_repo(&format!("{}/{}", &args.user, &repo));
    }

    return api::Client::new()
        .lookup_repos(&args.user)?
        .print_all_downloads();
}
