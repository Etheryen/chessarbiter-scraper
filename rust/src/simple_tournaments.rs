use crate::consts::BASE_URL;
use anyhow::{anyhow, Result};
use scraper::{selectable::Selectable, ElementRef, Html, Selector};
use std::fmt::Debug;
use thiserror::Error;
use url::Url;

#[derive(Debug)]
pub struct SimpleTournament {
    pub name: String,
    pub date: String,
    pub location: String,
    pub time_control: String,
    pub status: String,
    pub href: String,
}

pub async fn get_by_year_month(year: u32, month: u8) -> Result<Vec<SimpleTournament>> {
    let request_url = get_url(year, month)?;
    let html = get_html(request_url).await?;
    get_scraped_simple_tournaments(html)
}

fn get_scraped_simple_tournaments(html: Html) -> Result<Vec<SimpleTournament>> {
    let tbl_selector_text = r#"table[style="tbl"] > tbody > tr[class]"#;
    let tbl_selector = Selector::parse(tbl_selector_text)?;
    let rows = html.select(&tbl_selector);

    println!("{}", tbl_selector_text);

    let simple_tournaments = rows.filter_map(|row| parse_row(row).ok()).collect();

    Ok(simple_tournaments)
}

#[derive(Error, Debug)]
#[error("couldn't parse html")]
struct ScrapingError;

fn parse_row(row: ElementRef) -> Result<SimpleTournament> {
    let name_selector = Selector::parse("td > a")?;
    let name = row.select(&name_selector).next().unwrap().text().collect();

    let date_selector = Selector::parse("td[style]")?;
    let date = row
        .select(&date_selector)
        .next()
        .ok_or(ScrapingError)?
        .first_child()
        .ok_or(ScrapingError)?
        .value()
        .as_text()
        .ok_or(ScrapingError)?
        .to_string();

    let location_selector = Selector::parse("td + td > div.szary")?;
    let location = row
        .select(&location_selector)
        .next()
        .unwrap()
        .text()
        .collect::<String>()
        .chars()
        .take_while(|&ch| ch != '[')
        .collect::<String>()
        .as_str()
        .trim()
        .to_string();

    let time_control_selector = Selector::parse("td[width=\"12%\"] > div.szary")?;
    let time_control = row
        .select(&time_control_selector)
        .next()
        .unwrap()
        .text()
        .collect();

    let status_selector = Selector::parse("td[style] > div")?;
    let status = row
        .select(&status_selector)
        .next()
        .unwrap()
        .text()
        .collect();

    let href_selector = Selector::parse("td > a")?;
    let href = row
        .select(&href_selector)
        .next()
        .unwrap()
        .attr("href")
        .unwrap()
        .to_string();

    Ok(SimpleTournament {
        name,
        date,
        location,
        time_control,
        status,
        href,
    })
}

fn get_url(year: u32, month: u8) -> Result<Url> {
    let params = vec![("rok", year.to_string()), ("miesiac", month.to_string())];
    let mut request_url = Url::parse(BASE_URL)?;
    request_url.query_pairs_mut().extend_pairs(params);

    Ok(request_url)
}

async fn get_html(request_url: Url) -> Result<Html> {
    let resp = reqwest::get(request_url).await?;
    let text = resp.text().await?;
    Ok(Html::parse_document(&text))
}
