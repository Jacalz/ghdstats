use std::env;

mod fetch;
use fetch::{Repo, fetch_repos, fetch_statistics};

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
