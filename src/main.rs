mod api;
use std::env;
use std::error;

fn main() -> Result<(), Box<dyn error::Error>> {
    let args: Vec<String> = env::args().collect();

    let mut client = api::Client::new();
    match args.len() {
        0 | 1 | 4.. => {
            println!("Usage: gcdstats [user] [repository, optional]");
            return Ok(());
        }
        3 => return client.print_downloads_for_repo(&format!("{}/{}", &args[1], &args[2])),
        2 => {
            let parts: Vec<&str> = args[1].split('/').collect();
            match parts.len() {
                1 => client.lookup_repos(parts[0])?,
                2 => return client.print_downloads_for_repo(&args[1]),
                _ => {
                    println!("Invalid input format");
                    return Ok(());
                }
            }
        }
    }

    client.print_all_downloads();
    Ok(())
}
