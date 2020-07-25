mod ingredient;
mod recipe;

use clap::{App, Arg, SubCommand};
use recipe::{APIClient, Recipe};
use std::io::prelude::*;
use std::{
    fs::{self, File},
    path::Path,
};

fn main() {
    let matches = App::new("food CLI")
        .version("1.0")
        .arg(
            Arg::with_name("url")
                .short("u")
                .long("url")
                .value_name("URL")
                .help("Sets a custom URL")
                .takes_value(true),
        )
        .arg(
            Arg::with_name("v")
                .short("v")
                .multiple(true)
                .help("Sets the level of verbosity"),
        )
        .subcommand(
            SubCommand::with_name("get")
                .about("downloads recipes by uuid")
                .arg(
                    Arg::with_name("uuid")
                        .short("id")
                        .long("uuid")
                        .help("get by uuid")
                        .takes_value(true),
                ),
        )
        .subcommand(
            SubCommand::with_name("upload")
                .about("upload recipe by file")
                .arg(
                    Arg::with_name("file")
                        .short("f")
                        .long("file")
                        .value_name("FILE")
                        .help("")
                        .takes_value(true),
                ),
        )
        .get_matches();

    // Gets a value for config if supplied by user, or defaults to "default.conf"
    let url = matches
        .value_of("url")
        .unwrap_or("http://localhost:4242/query");

    let c = APIClient {
        url: url.to_string(),
    };

    if let Some(matches) = matches.subcommand_matches("get") {
        // let uuid = .unwrap();

        let uuid = match matches.value_of("uuid") {
            None => panic!("uuid is required"),
            Some(u) => u,
        };

        let r = c.get_recipe(uuid.to_string());
        write_yaml_to_file(r);
    }
    if let Some(matches) = matches.subcommand_matches("upload") {
        let file = match matches.value_of("file") {
            None => panic!("file is required"),
            Some(u) => u,
        };
        let r = read_yaml_file(file.to_string());
        let res = c.set_recipe(r);
        println!("uploaded: {:#?}", res);
    }
}

pub fn write_yaml_to_file(r: Recipe) {
    let y = serde_yaml::to_string(&r).unwrap();
    println!("{}", y);
    let file_name = format!("{}.yaml", r.uuid);
    let path = Path::new(&file_name);
    let display = path.display();

    let mut file = match File::create(&path) {
        Err(why) => panic!("couldn't create {}: {}", display, why),
        Ok(file) => file,
    };

    match file.write_all(y.as_bytes()) {
        Err(why) => panic!("couldn't write to {}: {}", display, why),
        Ok(_) => println!("successfully wrote to {}", display),
    }
}

pub fn read_yaml_file(filename: String) -> Recipe {
    let contents = fs::read_to_string(filename).expect("Something went wrong reading the file");

    let recipe: Recipe = serde_yaml::from_str(&contents).unwrap();
    return recipe;
}
