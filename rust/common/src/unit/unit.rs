use anyhow::bail;
use ingredient::unit::singular;
use ingredient::unit::{kind::MeasureKind, Unit};
use ingredient::Amount;
use petgraph::Graph;
use serde::{Deserialize, Serialize};
use tracing::{debug, info};

use crate::new_ingredient_parser;
type MeasureGraph = Graph<Unit, f64>;

pub fn make_graph(mappings: Vec<(Measure, Measure)>) -> MeasureGraph {
    let mut g = Graph::<Unit, f64>::new();

    for (m_a, m_b) in mappings.into_iter() {
        let n_a = g
            .node_indices()
            .find(|i| g[*i] == m_a.unit)
            .unwrap_or_else(|| g.add_node(m_a.unit.clone().normalize()));
        let n_b = g
            .node_indices()
            .find(|i| g[*i] == m_b.unit)
            .unwrap_or_else(|| g.add_node(m_b.unit.clone().normalize()));
        let _c1 = g.add_edge(n_a, n_b, m_b.value / m_a.value);
        let _c2 = g.add_edge(n_b, n_a, m_a.value / m_b.value);
    }
    return g;
}
pub fn print_graph(g: MeasureGraph) -> String {
    return format!("{}", petgraph::dot::Dot::new(&g));
}

#[tracing::instrument(name = "unit::convert")]
pub fn convert(
    m: Measure,
    target: MeasureKind,
    mappings: Vec<(Measure, Measure)>,
) -> Option<Measure> {
    let g = make_graph(mappings);

    let unit_a = m.unit.clone();
    let unit_b = target.unit();

    let n_a = g.node_indices().find(|i| g[*i] == unit_a)?;
    let n_b = g.node_indices().find(|i| g[*i] == unit_b)?;

    debug!("calculating {:?} to {:?}", n_a, n_b);
    if !petgraph::algo::has_path_connecting(&g, n_a, n_b, None) {
        debug!("convert failed for {:?}", m);
        return None;
    };

    let steps = petgraph::algo::astar(&g, n_a, |finish| finish == n_b, |e| *e.weight(), |_| 0.0)
        .unwrap()
        .1;
    let mut factor: f64 = 1.0;
    for x in 0..steps.len() - 1 {
        let edge = g
            .find_edge(*steps.get(x).unwrap(), *steps.get(x + 1).unwrap())
            .unwrap();
        factor *= g.edge_weight(edge).unwrap();
    }

    let result = Measure::new(
        unit_b,
        (m.value * factor * 100.0).round() / 100.0,
        match m.upper_value {
            Some(x) => Some((x * factor * 100.0).round() / 100.0),
            None => None,
        },
    );
    debug!("{:?} -> {:?} ({} hops)", m, result, steps.len());
    return Some(result);
}

pub fn add_time_amounts(a: Vec<Amount>) -> Amount {
    let mut m = Measure::from_str("0 seconds");
    for x in a.into_iter() {
        m = m.add(Measure::parse(x)).unwrap();
    }
    m.as_bare().unwrap()
}

#[derive(Clone, PartialEq, PartialOrd, Debug, Serialize, Deserialize)]
pub struct Measure {
    unit: Unit,
    value: f64,
    upper_value: Option<f64>,
}

// multiplication factors
const TSP_TO_TBSP: f64 = 3.0;
const TSP_TO_FL_OZ: f64 = 2.0;
const G_TO_K: f64 = 1000.0;
const CUP_TO_QUART: f64 = 4.0;
const TSP_TO_CUP: f64 = 48.0;
const GRAM_TO_OZ: f64 = 28.3495;
const OZ_TO_LB: f64 = 16.0;
const CENTS_TO_DOLLAR: f64 = 100.0;
const SEC_TO_MIN: f64 = 60.0;
const SEC_TO_HOUR: f64 = 3600.0;
const SEC_TO_DAY: f64 = 86400.0;

impl Measure {
    pub fn new(unit: Unit, value: f64, upper_value: Option<f64>) -> Measure {
        Measure {
            unit,
            value,
            upper_value,
        }
    }
    pub fn from_string(s: String) -> Measure {
        let a = new_ingredient_parser(false).parse_amount(s.as_str())[0].clone();
        Measure::parse(Amount {
            unit: singular(&a.unit),
            value: a.value,
            upper_value: a.upper_value,
        })
    }
    pub fn from_str(s: &str) -> Measure {
        let a = new_ingredient_parser(false).parse_amount(s)[0].clone();
        Measure::parse(Amount {
            unit: singular(&a.unit),
            value: a.value,
            upper_value: a.upper_value,
        })
    }
    pub fn normalize(&self) -> Measure {
        let (unit, factor) = match &self.unit {
            Unit::Teaspoon
            | Unit::Milliliter
            | Unit::Gram
            | Unit::Cent
            | Unit::KCal
            | Unit::Farhenheit
            | Unit::Celcius // todo: convert to farhenheit?
            | Unit::Inch
            | Unit::Second => return self.clone(),
            Unit::Other(x) => {
                let x2 = x.clone();
                let u2 = singular(&x2);
                return Measure::new(Unit::Other(u2), self.value, self.upper_value);
            }

            Unit::Kilogram => (Unit::Gram, G_TO_K),

            Unit::Ounce => (Unit::Gram, GRAM_TO_OZ),
            Unit::Pound => (Unit::Gram, GRAM_TO_OZ * OZ_TO_LB),

            Unit::Liter => (Unit::Milliliter, G_TO_K),

            Unit::Tablespoon => (Unit::Teaspoon, TSP_TO_TBSP),
            Unit::Cup => (Unit::Teaspoon, TSP_TO_CUP),
            Unit::Quart => (Unit::Teaspoon, CUP_TO_QUART * TSP_TO_CUP),
            Unit::FluidOunce => (Unit::Teaspoon, TSP_TO_FL_OZ),

            Unit::Dollar => (Unit::Cent, CENTS_TO_DOLLAR),
            Unit::Day => (Unit::Second, SEC_TO_DAY),
            Unit::Hour => (Unit::Second, SEC_TO_HOUR),
            Unit::Minute => (Unit::Second, SEC_TO_MIN),
        };
        return Measure {
            unit,
            value: self.value * factor,
            upper_value: match self.upper_value {
                Some(x) => Some(x * factor),
                None => None,
            },
        };
    }
    pub fn add(&self, b: Measure) -> Result<Measure, anyhow::Error> {
        info!("adding {:?} to {:?}", self, b);

        if let MeasureKind::Other = b.kind().unwrap() {
            return Ok(self.clone());
        }

        if self.kind().unwrap() != b.kind().unwrap() {
            return Err(anyhow::anyhow!(
                "Cannot add measures of different kinds: {:#?} {:?}",
                self,
                b
            ));
        }
        Ok(Measure {
            unit: self.unit.clone(),
            value: self.value + b.value,
            upper_value: match (self.upper_value, b.upper_value) {
                (Some(a), Some(b)) => Some(a + b),
                (None, None) => None,
                (None, Some(b)) => Some(self.value + b),
                (Some(a), None) => Some(a + b.value),
            },
        })
    }
    pub fn parse(m: Amount) -> Measure {
        Measure::new(
            Unit::from_str(singular(m.unit.as_ref()).as_ref()),
            m.value,
            m.upper_value,
        )
        .normalize()
    }
    pub fn kind(&self) -> Result<MeasureKind, anyhow::Error> {
        match self.unit {
            Unit::Gram => Ok(MeasureKind::Weight),
            Unit::Cent => Ok(MeasureKind::Money),
            Unit::Teaspoon | Unit::Milliliter => Ok(MeasureKind::Volume),
            Unit::KCal => Ok(MeasureKind::Calories),
            Unit::Second => Ok(MeasureKind::Time),
            Unit::Farhenheit | Unit::Celcius => Ok(MeasureKind::Temperature), // todo: convert to farhenheit?
            Unit::Inch => Ok(MeasureKind::Length),
            Unit::Other(_) => Ok(MeasureKind::Other),
            Unit::Kilogram
            | Unit::Liter
            | Unit::Tablespoon
            | Unit::Cup
            | Unit::Quart
            | Unit::FluidOunce
            | Unit::Ounce
            | Unit::Pound
            | Unit::Dollar
            | Unit::Day
            | Unit::Minute
            | Unit::Hour => bail!("unit not normalized: {:?}", self),
        }
    }

    pub fn as_raw(self) -> Amount {
        Amount {
            unit: self.unit.to_str(),
            value: self.value,
            upper_value: self.upper_value,
        }
    }

    pub fn as_bare(self) -> anyhow::Result<Amount> {
        let (u, f) = match self.unit {
            Unit::Gram => (Unit::Gram, 1.0),
            Unit::Milliliter => (Unit::Milliliter, 1.0),
            Unit::Teaspoon => match self.value {
                // only for these measurements to we convert to the best fit, others stay bare due to the nature of the values
                m if { m < 3.0 } => (Unit::Teaspoon, 1.0),
                m if { m < 12.0 } => (Unit::Tablespoon, TSP_TO_TBSP),
                m if { m < CUP_TO_QUART * TSP_TO_CUP } => (Unit::Cup, TSP_TO_CUP),
                _ => (Unit::Quart, CUP_TO_QUART * TSP_TO_CUP),
            },
            Unit::Cent => (Unit::Dollar, CENTS_TO_DOLLAR),
            Unit::KCal => (Unit::KCal, 1.0),
            Unit::Second => match self.value {
                // only for these measurements to we convert to the best fit, others stay bare due to the nature of the values
                m if { m < SEC_TO_MIN } => (Unit::Second, 1.0),
                m if { m < SEC_TO_HOUR } => (Unit::Minute, SEC_TO_MIN),
                m if { m < SEC_TO_DAY } => (Unit::Hour, SEC_TO_HOUR),
                _ => (Unit::Day, SEC_TO_DAY),
            },
            Unit::Inch => (Unit::Inch, 1.0),
            Unit::Other(o) => (Unit::Other(o), 1.0),
            Unit::Kilogram
            | Unit::Liter
            | Unit::Tablespoon
            | Unit::Cup
            | Unit::Quart
            | Unit::FluidOunce
            | Unit::Ounce
            | Unit::Pound
            | Unit::Dollar
            | Unit::Farhenheit
            | Unit::Celcius // todo: convert to farhenheit?
            | Unit::Minute
            | Unit::Hour
            | Unit::Day => bail!("unit not normalized: {:?}", self),
        };
        return Ok(Amount {
            unit: u.to_str(),
            value: self.value / f,
            upper_value: match self.upper_value {
                Some(x) => Some(x / f),
                None => None,
            },
        });
    }
}

#[cfg(test)]
mod tests {

    use super::*;
    #[test]
    fn test_measure() {
        let m1 = Measure::from_str("16 tbsp");
        assert_eq!(m1, Measure::new(Unit::Teaspoon, 48.0, None));
        assert_eq!(m1.as_bare().unwrap(), Amount::new("cup", 1.0));
        assert_eq!(
            Measure::from_str("25.2 grams").as_bare().unwrap(),
            Amount::new("g", 25.2)
        );
        assert_eq!(
            Measure::from_str("2500.2 grams").as_bare().unwrap(),
            Amount::new("g", 2500.2)
        );
        assert_eq!(
            Measure::from_str("12 foo").as_bare().unwrap(),
            Amount::new("whole", 12.0)
        );
    }

    #[test]
    fn test_convert() {
        let m = Measure::from_str("1 tbsp");
        let tbsp_dollars = (Measure::from_str("2 tbsp"), Measure::from_str("4 dollars"));
        assert_eq!(
            Measure::from_str("2 dollars"),
            convert(m.clone(), MeasureKind::Money, vec![tbsp_dollars.clone()]).unwrap()
        );

        assert!(convert(m, MeasureKind::Volume, vec![tbsp_dollars.clone()]).is_none());
    }
    #[test]
    fn test_convert_lb() {
        let grams_dollars = (Measure::from_str("1 gram"), Measure::from_str("1 dollar"));
        assert_eq!(
            Measure::from_str("2 dollars"),
            convert(
                Measure::from_str("2 grams"),
                MeasureKind::Money,
                vec![grams_dollars.clone()]
            )
            .unwrap()
        );
        assert_eq!(
            Measure::from_str("56.699 dollars"),
            convert(
                Measure::from_str("2 oz"),
                MeasureKind::Money,
                vec![grams_dollars.clone()]
            )
            .unwrap()
        );
        assert_eq!(
            Measure::from_str("226.796 dollars"),
            convert(
                Measure::from_str(".5 lb"),
                MeasureKind::Money,
                vec![grams_dollars.clone()]
            )
            .unwrap()
        );
        assert_eq!(
            Measure::from_str("453.592 dollars"),
            convert(
                Measure::from_str("1 lb"),
                MeasureKind::Money,
                vec![grams_dollars.clone()]
            )
            .unwrap()
        );
    }
    #[test]
    fn test_convert_other() {
        assert_eq!(
            Measure::from_str("10.0 cents"),
            convert(
                Measure::from_str("1 whole"),
                MeasureKind::Money,
                vec![(
                    Measure::from_str("12 whole"),
                    Measure::from_str("1.20 dollar"),
                )]
            )
            .unwrap()
        );
    }
    #[test]
    fn test_convert_range() {
        assert_eq!(
            Measure::from_str("5-10 dollars"),
            convert(
                Measure::from_str("1-2 whole"),
                MeasureKind::Money,
                vec![(Measure::from_str("4 whole"), Measure::from_str("20 dollar"))]
            )
            .unwrap()
        );
    }
    #[test]
    fn test_convert_transitive() {
        assert_eq!(
            Measure::from_str("1 cent"),
            convert(
                Measure::from_str("1 grams"),
                MeasureKind::Money,
                vec![
                    (Measure::from_str("1 cent"), Measure::from_str("1 tsp"),),
                    (Measure::from_str("1 grams"), Measure::from_str("1 tsp"),),
                ]
            )
            .unwrap()
        );
        assert_eq!(
            Measure::from_str("1 dollar"),
            convert(
                Measure::from_str("1 grams"),
                MeasureKind::Money,
                vec![
                    (Measure::from_str("1 dollar"), Measure::from_str("1 cup"),),
                    (Measure::from_str("1 grams"), Measure::from_str("1 cup"),),
                ]
            )
            .unwrap()
        );
    }
    #[test]
    fn test_convert_kcal() {
        assert_eq!(
            Measure::from_str("200 kcal"),
            convert(
                Measure::from_str("100 g"),
                MeasureKind::Calories,
                vec![
                    (Measure::from_str("20 cups"), Measure::from_str("40 grams"),),
                    (Measure::from_str("20 grams"), Measure::from_str("40 kcal"),)
                ]
            )
            .unwrap()
        );
    }
    #[test]
    fn test_add() {
        assert_eq!(
            Measure::from_str("10 minutes")
                .add(Measure::from_str("2-3 minutes"))
                .unwrap(),
            Measure::from_str("12-13 minutes"),
        );
    }
}
