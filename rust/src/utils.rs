use crate::simple_tournaments::SimpleTournament;

pub trait Capitalize {
    fn capitalize(&self) -> String;
}

impl Capitalize for str {
    fn capitalize(&self) -> String {
        let mut c = self.chars();
        match c.next() {
            None => String::new(),
            Some(f) => f.to_uppercase().collect::<String>() + c.as_str().to_lowercase().as_str(),
        }
    }
}

pub fn print_simple_tournaments(tournaments: &[SimpleTournament]) {
    println!("[");
    tournaments.iter().enumerate().for_each(|(i, t)| {
        println!("  name: {}", t.name);
        println!("  date: {}", t.date);
        println!("  location: {}", t.location);
        println!("  time_control: {}", t.time_control);
        println!("  status: {}", t.status);
        println!("  href: {}", t.href);

        if i != tournaments.len() - 1 {
            println!();
        }
    });
    println!("]");
}
