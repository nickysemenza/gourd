use serde::{Deserialize, Serialize};

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
    Other,
}
impl MeasureKind {
    pub fn from_str(s: &str) -> Self {
        match s {
            "weight" => Self::Weight,
            "volume" => Self::Volume,
            "money" => Self::Money,
            _ => Self::Other,
        }
    }
}
#[derive(Clone, PartialEq, PartialOrd, Debug)]
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
    Other(String),
}

impl Unit {
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
            Unit::Other(s) => return s,
        }
        .to_string()
    }
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
        let a = ingredient::parse_amount(s.as_str()).unwrap()[0].clone();
        Measure::parse(BareMeasurement::new(a.unit, a.value))
    }
    pub fn normalize(&self) -> Measure {
        let foo = match self.0 {
            Unit::Teaspoon | Unit::Milliliter | Unit::Gram | Unit::Cent | Unit::Other(_) => {
                return self.clone()
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
        // return Measure(foo.0, foo.1);
    }
    pub fn kind(&self) -> MeasureKind {
        return match self.0 {
            Unit::Gram => MeasureKind::Weight,
            Unit::Cent => MeasureKind::Money,
            Unit::Teaspoon | Unit::Milliliter => MeasureKind::Volume,

            Unit::Other(_) => MeasureKind::Other,
            _ => panic!("unit not normalized: {:?}", self),
        };
    }

    pub fn convert(
        &self,
        target: MeasureKind,
        mappings: Vec<(Measure, Measure)>,
    ) -> Option<Measure> {
        let curr_kind = self.kind();
        for (m_a, m_b) in mappings.into_iter() {
            let a = m_a.clone().normalize();
            let b = m_b.clone().normalize();
            let (kind_a, kind_b) = (a.kind(), b.kind());
            if kind_a == target && kind_b == curr_kind {
                return Some(Measure(a.0, a.1 / b.1 * self.1));
            } else if kind_a == curr_kind && kind_b == target {
                return Some(Measure(b.0, b.1 / a.1 * self.1));
            }
        }
        None
    }
    pub fn as_bare(self) -> BareMeasurement {
        let m = self.1;
        let (val, u, f) = match self.0 {
            Unit::Gram => {
                if m < 1000.0 {
                    (m, Unit::Gram, 1.0)
                } else {
                    (m, Unit::Kilogram, G_TO_K)
                }
            }
            Unit::Milliliter => {
                if m < 1000.0 {
                    (m, Unit::Milliliter, 1.0)
                } else {
                    (m, Unit::Liter, G_TO_K)
                }
            }
            Unit::Teaspoon => match m {
                m if { m < 3.0 } => (m, Unit::Teaspoon, 1.0),
                m if { m < 12.0 } => (m, Unit::Tablespoon, TSP_TO_TBSP),
                m if { m < CUP_TO_QUART * TSP_TO_CUP } => (m, Unit::Cup, TSP_TO_CUP),
                _ => (m, Unit::Teaspoon, 1.0),
            },

            Unit::Cent => match m {
                m if { m < CENTS_TO_DOLLAR } => (m, Unit::Cent, 1.0),
                _ => (m, Unit::Dollar, CENTS_TO_DOLLAR),
            },
            Unit::Other(o) => (m, Unit::Other(o), 1.0),
            _ => panic!("unit not normalized: {:?}", self),
        };
        return BareMeasurement::new(u.to_str(), val / f);
    }

    // Err("todo".to_string())
}
pub fn singular(s: &str) -> String {
    s.strip_suffix("s").unwrap_or(s).to_lowercase()
}

#[cfg(test)]
mod tests {

    use super::*;
    #[test]
    fn test_measure() {
        // let m1 = Measure::parse(Measurement("Tbsp".to_string(), 16.0));
        let m1 = Measure::from_string("16 tbsp".to_string());
        assert_eq!(m1, Measure(Unit::Teaspoon, 48.0));
        assert_eq!(m1.as_bare(), BareMeasurement::new("cup".to_string(), 1.0));
        assert_eq!(
            Measure::from_string("25.2 grams".to_string()).as_bare(),
            BareMeasurement::new("g".to_string(), 25.2)
        );
        assert_eq!(
            Measure::from_string("2500.2 grams".to_string()).as_bare(),
            BareMeasurement::new("kg".to_string(), 2.5002)
        );
        assert_eq!(
            Measure::from_string("12 foo".to_string()).as_bare(),
            BareMeasurement::new("foo".to_string(), 12.0)
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
            Measure::from_string("10.000001 cents".to_string()),
            Measure::from_string("1 egg".to_string())
                .convert(
                    MeasureKind::Money,
                    vec![(
                        Measure::from_string("12 eggs".to_string()),
                        Measure::from_string("1.20 dollar".to_string()),
                    )]
                )
                .unwrap()
        );
    }
}
