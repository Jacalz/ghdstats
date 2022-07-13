use std::env;

mod repo;
use repo::{Repository, fetch_repositories};

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() == 0 || args.len() > 2 {
        println!("Usage: gcdstats [user] [repository, optional]")
    }

    let repos: Vec<Repository>;
    if args.len() == 2 {
        repos = vec![Repository{full_name: args[0]+"/"+&args[1] }];
    } else if args[0].contains('/') {
        repos = vec![Repository{full_name: args[0] }];
    } else {
		repos = fetch_repositories(args[0]);
	}

    // fetch_statistics(repos);
}
