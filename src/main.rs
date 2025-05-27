mod api;
use std::env;
use std::error;

fn main() -> Result<(), Box<dyn error::Error>> {
    let args: Vec<String> = env::args().collect();

    let mut client = api::Client::new()?;
    match args.len() {
        0 | 1 | 4.. => {
            println!("Usage: gcdstats [user] [repository, optional]");
            return Ok(());
        }
        3 => client.add_repo(format!("{}/{}", &args[1], &args[2])),
        2 => {
            let parts: Vec<&str> = args[1].split('/').collect();
            match parts.len() {
                1 => client.lookup_repos(parts[0])?,
                2 => client.add_repo(args[1].clone()),
                _ => {
                    println!("Invalid input format");
                    return Ok(());
                }
            }
        }
    }

    client.print_downloads();
    Ok(())
}
