mod consts;
mod simple_tournaments;
mod utils;

use anyhow::Result;
use simple_tournaments::get_by_year_month;
use utils::{print_simple_tournaments, Capitalize};

#[tokio::main]
async fn main() -> Result<()> {
    let text_pl = "ŁOMŻA";
    println!("{}", text_pl);
    let title = text_pl.capitalize();
    println!("{}", title);

    let simple_tournaments = get_by_year_month(2024, 6).await?;

    print_simple_tournaments(&simple_tournaments);

    Ok(())
}
