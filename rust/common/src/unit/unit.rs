use std::fmt;

use anyhow::bail;
use petgraph::Graph;
use serde::{Deserialize, Serialize};
use tracing::debug;

#[derive(Clone, PartialEq, PartialOrd, Debug, Default, Serialize, Deserialize)]
pub struct BareMeasurement {
    pub unit: String,
    pub value: f32,
}

impl BareMeasurement {
    pub fn new(unit: String, value: f32) -> BareMeasurement {
        BareMeasurement { unit, value }
    }
}

#[derive(Clone, PartialEq, PartialOrd, Debug)]
pub struct Measure(Unit, f32);

#[derive(Clone, PartialEq, PartialOrd, Debug)]
pub enum MeasureKind {
    Weight,
    Volume,
    Money,
    Calories,
    Other,
}
impl MeasureKind {
    pub fn from_str(s: &str) -> Self {
        match s {
            "weight" => Self::Weight,
            "volume" => Self::Volume,
            "money" => Self::Money,
            "calories" => Self::Calories,
            _ => Self::Other,
        }
    }
}
#[derive(Clone, PartialEq, PartialOrd, Debug, Eq, Hash)]
pub enum Unit {
    Gram,
    Kilogram,
    Liter,
    Milliliter,
    Teaspoon,
    Tablespoon,
    Cup,
    Quart,
    FluidOunce,
    Ounce,
    Pound,
    Cent,
    Dollar,
    KCal,
    Other(String),
}

impl Unit {
    pub fn normalize(self) -> Unit {
        match self {
            Unit::Other(x) => return Unit::Other(singular(&x)),
            _ => return self,
        }
    }
    pub fn from_str(s: &str) -> Self {
        match s {
            "gram" | "g" => Self::Gram,
            "kilogram" | "kg" => Self::Kilogram,

            "oz" | "ounce" => Self::Ounce,
            "lb" | "pound" => Self::Pound,

            "ml" => Self::Milliliter,
            "l" => Self::Liter,

            "tsp" | "teaspoon" => Self::Teaspoon,
            "tbsp" | "tablespoon" => Self::Tablespoon,
            "c" | "cup" => Self::Cup,
            "q" | "quart" => Self::Quart,
            "fl oz" | "fluid oz" => Self::FluidOunce,

            "dollar" | "$" => Self::Dollar,
            "cent" => Self::Cent,

            "calorie" | "cal" | "kcal" => Self::KCal,
            _ => Self::Other(s.to_string()),
        }
    }
    pub fn to_str(self) -> String {
        match self {
            Unit::Gram => "g",
            Unit::Kilogram => "kg",
            Unit::Liter => "l",
            Unit::Milliliter => "ml",
            Unit::Teaspoon => "tsp",
            Unit::Tablespoon => "tbsp",
            Unit::Cup => "cup",
            Unit::Quart => "quart",
            Unit::FluidOunce => "fl oz",
            Unit::Ounce => "oz",
            Unit::Pound => "lb",
            Unit::Cent => "cent",
            Unit::Dollar => "$",
            Unit::KCal => "kcal",
            Unit::Other(s) => return singular(&s),
        }
        .to_string()
    }
}
pub fn unit_from_measurekind(m: MeasureKind) -> Unit {
    return match m {
        MeasureKind::Weight => Unit::Gram,
        MeasureKind::Volume => Unit::Milliliter,
        MeasureKind::Money => Unit::Cent,
        MeasureKind::Calories => Unit::KCal,
        MeasureKind::Other => Unit::Other("".to_string()),
    };
}
impl fmt::Display for Unit {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "{:?}", self)
    }
}

pub fn make_graph(mappings: Vec<(Measure, Measure)>) -> Graph<Unit, f32> {
    let mut g = Graph::<Unit, f32>::new();

    for (m_a, m_b) in mappings.into_iter() {
        let n_a = g
            .node_indices()
            .find(|i| g[*i] == m_a.0)
            .unwrap_or_else(|| g.add_node(m_a.0.clone().normalize()));
        let n_b = g
            .node_indices()
            .find(|i| g[*i] == m_b.0)
            .unwrap_or_else(|| g.add_node(m_b.0.clone().normalize()));
        let _c1 = g.add_edge(n_a, n_b, m_b.1 / m_a.1);
        let _c2 = g.add_edge(n_b, n_a, m_a.1 / m_b.1);
    }
    return g;
}
pub fn print_graph(g: Graph<Unit, f32>) -> String {
    return format!("{}", petgraph::dot::Dot::new(&g));
}

// multiplication factors
const TSP_TO_TBSP: f32 = 3.0;
const TSP_TO_FL_OZ: f32 = 2.0;
const G_TO_K: f32 = 1000.0;
const CUP_TO_QUART: f32 = 4.0;
const TSP_TO_CUP: f32 = 48.0;
const GRAM_TO_OZ: f32 = 28.3495;
const OZ_TO_LB: f32 = 16.0;
const CENTS_TO_DOLLAR: f32 = 100.0;

impl Measure {
    pub fn from_string(s: String) -> Measure {
        let a = ingredient::parse_amount(s.as_str())[0].clone();
        Measure::parse(BareMeasurement::new(singular(&a.unit), a.value))
    }
    pub fn normalize(&self) -> Measure {
        let foo = match &self.0 {
            Unit::Teaspoon | Unit::Milliliter | Unit::Gram | Unit::Cent | Unit::KCal => {
                return self.clone()
            }
            Unit::Other(x) => {
                let x2 = x.clone();
                let u2 = singular(&x2);
                return Measure(Unit::Other(u2), self.1);
            }

            Unit::Kilogram => (Unit::Gram, self.1 * G_TO_K),

            Unit::Ounce => (Unit::Gram, self.1 * GRAM_TO_OZ),
            Unit::Pound => (Unit::Gram, self.1 * GRAM_TO_OZ * OZ_TO_LB),

            Unit::Liter => (Unit::Milliliter, self.1 * G_TO_K),

            Unit::Tablespoon => (Unit::Teaspoon, self.1 * TSP_TO_TBSP),
            Unit::Cup => (Unit::Teaspoon, self.1 * TSP_TO_CUP),
            Unit::Quart => (Unit::Teaspoon, self.1 * CUP_TO_QUART * TSP_TO_CUP),
            Unit::FluidOunce => (Unit::Teaspoon, self.1 * TSP_TO_FL_OZ),

            Unit::Dollar => (Unit::Cent, self.1 * CENTS_TO_DOLLAR),
        };
        return Measure(foo.0, foo.1);
    }
    pub fn parse(m: BareMeasurement) -> Measure {
        Measure(Unit::from_str(singular(m.unit.as_ref()).as_ref()), m.value).normalize()
    }
    pub fn kind(&self) -> Result<MeasureKind, anyhow::Error> {
        match self.0 {
            Unit::Gram => Ok(MeasureKind::Weight),
            Unit::Cent => Ok(MeasureKind::Money),
            Unit::Teaspoon | Unit::Milliliter => Ok(MeasureKind::Volume),
            Unit::KCal => Ok(MeasureKind::Calories),
            Unit::Other(_) => Ok(MeasureKind::Other),
            Unit::Kilogram
            | Unit::Liter
            | Unit::Tablespoon
            | Unit::Cup
            | Unit::Quart
            | Unit::FluidOunce
            | Unit::Ounce
            | Unit::Pound
            | Unit::Dollar => bail!("unit not normalized: {:?}", self),
        }
    }

    #[tracing::instrument(name = "unit::convert")]
    pub fn convert(
        &self,
        target: MeasureKind,
        mappings: Vec<(Measure, Measure)>,
    ) -> Option<Measure> {
        let g = make_graph(mappings);

        let unit_a = self.0.clone();
        let unit_b = unit_from_measurekind(target);

        let n_a = g.node_indices().find(|i| g[*i] == unit_a)?;
        let n_b = g.node_indices().find(|i| g[*i] == unit_b)?;

        debug!("calculating {:?} to {:?}", n_a, n_b);
        if !petgraph::algo::has_path_connecting(&g, n_a, n_b, None) {
            debug!("convert failed for {:?}", self);
            return None;
        };

        let steps =
            petgraph::algo::astar(&g, n_a, |finish| finish == n_b, |e| *e.weight(), |_| 0.0)
                .unwrap()
                .1;
        let mut res: f32 = self.1;
        for x in 0..steps.len() - 1 {
            let edge = g
                .find_edge(*steps.get(x).unwrap(), *steps.get(x + 1).unwrap())
                .unwrap();
            res *= g.edge_weight(edge).unwrap();
        }
        let y = (res * 100.0).round() / 100.0;
        let result = Measure(unit_b, y);
        debug!("{:?} -> {:?} ({} hops)", self, result, steps.len());
        return Some(result);
    }

    pub fn as_bare(self) -> anyhow::Result<BareMeasurement> {
        let m = self.1;
        let (val, u, f) = match self.0 {
            Unit::Gram => (m, Unit::Gram, 1.0),
            Unit::Milliliter => (m, Unit::Milliliter, 1.0),
            Unit::Teaspoon => match m {
                // only for these measurements to we convert to the best fit, others stay bare due to the nature of the values
                m if { m < 3.0 } => (m, Unit::Teaspoon, 1.0),
                m if { m < 12.0 } => (m, Unit::Tablespoon, TSP_TO_TBSP),
                m if { m < CUP_TO_QUART * TSP_TO_CUP } => (m, Unit::Cup, TSP_TO_CUP),
                _ => (m, Unit::Teaspoon, 1.0),
            },
            Unit::Cent => (m, Unit::Dollar, CENTS_TO_DOLLAR),
            Unit::KCal => (m, Unit::KCal, 1.0),

            Unit::Other(o) => (m, Unit::Other(o), 1.0),
            Unit::Kilogram
            | Unit::Liter
            | Unit::Tablespoon
            | Unit::Cup
            | Unit::Quart
            | Unit::FluidOunce
            | Unit::Ounce
            | Unit::Pound
            | Unit::Dollar => bail!("unit not normalized: {:?}", self),
        };
        return Ok(BareMeasurement::new(u.to_str(), val / f));
    }
}
pub fn singular(s: &str) -> String {
    s.strip_suffix("s").unwrap_or(s).to_lowercase()
}

#[cfg(test)]
mod tests {

    use super::*;
    #[test]
    fn test_measure() {
        let m1 = Measure::from_string("16 tbsp".to_string());
        assert_eq!(m1, Measure(Unit::Teaspoon, 48.0));
        assert_eq!(
            m1.as_bare().unwrap(),
            BareMeasurement::new("cup".to_string(), 1.0)
        );
        assert_eq!(
            Measure::from_string("25.2 grams".to_string())
                .as_bare()
                .unwrap(),
            BareMeasurement::new("g".to_string(), 25.2)
        );
        assert_eq!(
            Measure::from_string("2500.2 grams".to_string())
                .as_bare()
                .unwrap(),
            BareMeasurement::new("g".to_string(), 2500.2)
        );
        assert_eq!(
            Measure::from_string("12 foo".to_string())
                .as_bare()
                .unwrap(),
            BareMeasurement::new("whole".to_string(), 12.0)
        );
    }

    #[test]
    fn test_convert() {
        let m = Measure::from_string("1 tbsp".to_string());
        let tbsp_dollars = (
            Measure::from_string("2 tbsp".to_string()),
            Measure::from_string("4 dollars".to_string()),
        );
        assert_eq!(
            Measure::from_string("2 dollars".to_string()),
            m.convert(MeasureKind::Money, vec![tbsp_dollars.clone()])
                .unwrap()
        );

        assert!(m
            .convert(MeasureKind::Volume, vec![tbsp_dollars.clone()])
            .is_none());
    }
    #[test]
    fn test_convert_lb() {
        let grams_dollars = (
            Measure::from_string("1 gram".to_string()),
            Measure::from_string("1 dollar".to_string()),
        );
        assert_eq!(
            Measure::from_string("2 dollars".to_string()),
            Measure::from_string("2 grams".to_string())
                .convert(MeasureKind::Money, vec![grams_dollars.clone()])
                .unwrap()
        );
        assert_eq!(
            Measure::from_string("56.699 dollars".to_string()),
            Measure::from_string("2 oz".to_string())
                .convert(MeasureKind::Money, vec![grams_dollars.clone()])
                .unwrap()
        );
        assert_eq!(
            Measure::from_string("226.796 dollars".to_string()),
            Measure::from_string(".5 lb".to_string())
                .convert(MeasureKind::Money, vec![grams_dollars.clone()])
                .unwrap()
        );
        assert_eq!(
            Measure::from_string("453.592 dollars".to_string()),
            Measure::from_string("1 lb".to_string())
                .convert(MeasureKind::Money, vec![grams_dollars.clone()])
                .unwrap()
        );
    }
    #[test]
    fn test_convert_other() {
        assert_eq!(
            Measure::from_string("10.0 cents".to_string()),
            Measure::from_string("1 whole".to_string())
                .convert(
                    MeasureKind::Money,
                    vec![(
                        Measure::from_string("12 whole".to_string()),
                        Measure::from_string("1.20 dollar".to_string()),
                    )]
                )
                .unwrap()
        );
    }
    #[test]
    fn test_convert_transitive() {
        assert_eq!(
            Measure::from_string("1 cent".to_string()),
            Measure::from_string("1 grams".to_string())
                .convert(
                    MeasureKind::Money,
                    vec![
                        (
                            Measure::from_string("1 cent".to_string()),
                            Measure::from_string("1 tsp".to_string()),
                        ),
                        (
                            Measure::from_string("1 grams".to_string()),
                            Measure::from_string("1 tsp".to_string()),
                        ),
                    ]
                )
                .unwrap()
        );
        assert_eq!(
            Measure::from_string("1 dollar".to_string()),
            Measure::from_string("1 grams".to_string())
                .convert(
                    MeasureKind::Money,
                    vec![
                        (
                            Measure::from_string("1 dollar".to_string()),
                            Measure::from_string("1 cup".to_string()),
                        ),
                        (
                            Measure::from_string("1 grams".to_string()),
                            Measure::from_string("1 cup".to_string()),
                        ),
                    ]
                )
                .unwrap()
        );
    }
    #[test]
    fn test_convert_kcal() {
        assert_eq!(
            Measure::from_string("200 kcal".to_string()),
            Measure::from_string("100 g".to_string())
                .convert(
                    MeasureKind::Calories,
                    vec![
                        (
                            Measure::from_string("20 cups".to_string()),
                            Measure::from_string("40 grams".to_string()),
                        ),
                        (
                            Measure::from_string("20 grams".to_string()),
                            Measure::from_string("40 kcal".to_string()),
                        )
                    ]
                )
                .unwrap()
        );
    }
}
