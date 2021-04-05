use serde::{Deserialize, Serialize};

#[derive(Clone, PartialEq, PartialOrd, Debug, Default, Serialize, Deserialize)]
pub struct BareMeasurement {
    unit: String,
    value: f32,
}

impl BareMeasurement {
    pub fn new(unit: String, value: f32) -> BareMeasurement {
        BareMeasurement { unit, value }
    }
}

#[derive(Clone, PartialEq, PartialOrd, Debug)]
pub struct Measure(Unit, f32);

// #[derive(Clone, PartialEq, PartialOrd, Debug)]
// pub enum Measure {
//     Other(BareMeasurement),
//     Grams(f32),
//     Ml(f32),
//     Teaspoon(f32),
//     Cent(f32),
// }

#[derive(Clone, PartialEq, PartialOrd, Debug)]
pub enum MeasureKind {
    Weight,
    Volume,
    Money,
    Other,
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

impl Measure {
    pub fn from_string(s: String) -> Measure {
        let a = ingredient::parse_amount(s.as_str()).unwrap()[0].clone();
        Measure::parse(BareMeasurement::new(a.unit, a.value))
    }
    pub fn normalize(self) -> Measure {
        let m = self.clone();
        let foo = match self.0 {
            Unit::Teaspoon | Unit::Milliliter | Unit::Gram | Unit::Cent | Unit::Other(_) => {
                return m
            }

            Unit::Kilogram => (Unit::Gram, m.1 * G_TO_K),

            Unit::Ounce => (Unit::Gram, m.1 * GRAM_TO_OZ),

            Unit::Liter => (Unit::Milliliter, m.1 * G_TO_K),

            Unit::Tablespoon => (Unit::Teaspoon, m.1 * TSP_TO_TBSP),
            Unit::Cup => (Unit::Teaspoon, m.1 * TSP_TO_CUP),
            Unit::Quart => (Unit::Teaspoon, m.1 * CUP_TO_QUART * TSP_TO_CUP),
            Unit::FluidOunce => (Unit::Teaspoon, m.1 * TSP_TO_FL_OZ),

            Unit::Dollar => (Unit::Cent, m.1 * 100.0),
        };
        return Measure(foo.0, foo.1);
    }
    pub fn parse(m: BareMeasurement) -> Measure {
        let foo = Measure(Unit::from_str(singular(m.unit.as_ref()).as_ref()), m.value).normalize();
        return Measure(foo.0, foo.1);
    }
    pub fn kind(self) -> MeasureKind {
        return match self.0 {
            Unit::Gram => MeasureKind::Weight,
            Unit::Cent => MeasureKind::Money,
            Unit::Teaspoon | Unit::Milliliter => MeasureKind::Volume,

            Unit::Other(_) => MeasureKind::Other,
            _ => panic!("unit not normalized: {:?}", self),
        };
    }

    pub fn convert(
        self,
        target: MeasureKind,
        mappings: Vec<(Measure, Measure)>,
    ) -> Option<Measure> {
        let curr_kind = self.clone().kind();
        for m in mappings.iter() {
            let (kind_a, kind_b) = (m.0.clone().kind(), m.1.clone().kind());
            if kind_a == target && kind_b == curr_kind {
                return Some(Measure(
                    m.0.clone().normalize().0.clone(),
                    m.0.clone().normalize().1 / m.1.clone().normalize().1 * self.clone().1,
                ));
            } else if kind_a == curr_kind && kind_b == target {
                dbg!(m);
                return Some(Measure(
                    m.1.clone().normalize().0.clone(),
                    m.1.clone().normalize().1 / dbg!(m.0.clone().normalize()).1 * self.clone().1,
                ));
            }
        }

        None
        // Measure(Unit::Other("foo".to_string()), 1.0)
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

            Unit::Cent => (m, Unit::Cent, 1.0),
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
            m.convert(MeasureKind::Money, vec![tbsp_dollars]).unwrap()
        );
    }
}
