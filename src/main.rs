mod api;
use anyhow::Result;
use std::env;
use std::process;

fn main() -> Result<()> {
    let args: Vec<String> = env::args().collect();
    match args.len() {
        0 | 1 | 4.. => {
            eprintln!("Usage: gcdstats [user] [repository, optional]");
            process::exit(1);
        }
        3 => api::Client::print_downloads_for_repo(&format!("{}/{}", &args[1], &args[2])),
        2 => {
            let parts: Vec<&str> = args[1].split('/').collect();
            match parts.len() {
                1 => api::Client::new()
                    .lookup_repos(parts[0])?
                    .print_all_downloads(),
                2 => api::Client::print_downloads_for_repo(&args[1]),
                _ => {
                    eprintln!("Invalid input format");
                    process::exit(1);
                }
            }
        }
    }
}
